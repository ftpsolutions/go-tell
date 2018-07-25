package worker

import (
	"errors"
	"log"
	"os"

	gotell "github.com/ftpsolutions/go-tell"
)

var stdoutLogger = log.New(os.Stdout, "", 1)

type Worker struct {
	jobHandler    gotell.JobHandler
	store         gotell.Store
	stopChan      chan struct{}
	retryStrategy RetryStrategy

	Logger *log.Logger
}

func (w *Worker) Close() error {
	w.stopChan <- struct{}{}
	return nil
}

// Should this have error handling to report to the main worker loop?
func (w *Worker) handleJob(job *gotell.Job) {
	w.Logger.Printf("Handling job %v: %v", job.ID, job.Data)

	// Send our job as values to the handler.
	err := w.jobHandler(gotell.Job{
		ID:     job.ID,
		Status: job.Status,
		Type:   job.Type,
		Data:   job.Data,
	})

	if err != nil {
		w.Logger.Printf("Error handling job ID %v: %v, retrying", job.ID, err)
		err = w.retryStrategy(w.store, job)
		if err != nil {
			w.Logger.Println("Error retrying job", job.ID, err)
		}
		return
	}
	w.Logger.Println("Job completed", job.ID)

	err = w.store.CompleteJob(job)
	if err != nil {
		w.Logger.Println("Error completing job", job.ID, err)
		// TODO work out what to do when the job fails to complete.
	}
}

func run(w *Worker) {
	w.Logger.Println("Worker starting", w)
	for {
		waitingForAJob := w.store.WaitToDoJob()

		w.Logger.Println("Worker waiting", w)
		select {
		case job := <-waitingForAJob:
			if job == nil {
				continue
			}
			w.handleJob(job)

		case <-w.stopChan:
			w.Logger.Println("Worker stopping", w)
			return
		}
	}
}

func Open(
	store gotell.Store,
	jobHandler gotell.JobHandler,
	retryStrategy RetryStrategy,
	logger *log.Logger,
) (*Worker, error) {
	// Default to stdout logging
	if logger == nil {
		logger = stdoutLogger
	}
	// Default to one retry attempt
	if retryStrategy == nil {
		retryStrategy = OneAttempt
	}
	w := Worker{
		jobHandler:    jobHandler,
		store:         store,
		retryStrategy: retryStrategy,
		stopChan:      make(chan struct{}),

		Logger: logger,
	}
	go run(&w)
	return &w, nil
}

type RetryStrategy func(gotell.Store, *gotell.Job) error

func OneAttempt(store gotell.Store, job *gotell.Job) error {
	return store.FailedJob(job)
}

func AlwaysRetry(store gotell.Store, job *gotell.Job) error {
	return store.ReturnJob(job)
}

func RetryUntil(failureLimit int) RetryStrategy {
	return func(store gotell.Store, job *gotell.Job) error {
		job.RetryCount++
		if job.RetryCount >= failureLimit {
			err := store.FailedJob(job)
			if err != nil {
				return err
			}
			return errors.New("Retry limit reached")
		}
		return store.ReturnJob(job)
	}
}
