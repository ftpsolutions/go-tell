package worker

import (
	"errors"
	"log"
	"os"

	"github.com/ftpsolutions/go-tell/store"
)

var stdoutLogger = log.New(os.Stdout, "", 1)

type Job = store.Job
type Store = store.Store

type Worker struct {
	jobHandler    JobHandler
	store         Store
	stopChan      chan struct{}
	retryStrategy RetryStrategy

	Logger *log.Logger
}

func (w *Worker) Close() error {
	w.stopChan <- struct{}{}
	return nil
}

// Should this have error handling to report to the main worker loop?
func (w *Worker) handleJob(job *Job) {
	w.Logger.Println("Handling job", job)

	// Send our job as values to the handler.
	err := w.jobHandler(Job{
		ID:     job.ID,
		Status: job.Status,
		Type:   job.Type,
		Data:   job.Data,
	})

	if err != nil {
		w.Logger.Println("Error handling job", job, err)
		err = w.retryStrategy(w.store, job)
		if err != nil {
			w.Logger.Println("Error retrying job", job, err)
		}
		return
	}
	w.Logger.Println("Job completed", job)

	err = w.store.CompleteJob(job)
	if err != nil {
		w.Logger.Println("Error completing job", job, err)
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
	store Store,
	jobHandler JobHandler,
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

type RetryStrategy func(Store, *Job) error

// Retry Strategies
func OneAttempt(store Store, job *Job) error {
	return store.FailedJob(job)
}

func AlwaysRetry(store Store, job *Job) error {
	return store.ReturnJob(job)
}

func RetryUntil(failureLimit int) RetryStrategy {
	return func(store Store, job *Job) error {
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
