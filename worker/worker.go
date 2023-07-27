package worker

import (
	"errors"
	"log"
	"os"
	"time"

	gotell "github.com/ftpsolutions/go-tell"
)

var stdoutLogger = log.New(os.Stdout, "", 1)

type Worker struct {
	jobHandler    gotell.JobHandler
	store         gotell.Store
	runningChan   chan struct{}
	stoppingChan  chan struct{}
	stoppedChan   chan struct{}
	retryStrategy RetryStrategy

	Logger *log.Logger
}

func (w *Worker) Close() error {
	//tell the worker to stop
	w.stoppingChan <- struct{}{}
	//wait for it to actually stop
	<-w.stoppedChan
	return nil
}

func (w *Worker) WaitTillRunning() {
	//wait for worker to actually start its run loop (for testing purposes)
	<-w.runningChan
}

// Should this have error handling to report to the main worker loop?
func (w *Worker) handleJob(job *gotell.Job) {
	w.Logger.Printf("Handling job %v: %v - Created at %v", job.ID, job.Data, job.Created)

	// Send our job as values to the handler.
	// TODO make this a job.Copy()
	err := w.jobHandler(gotell.Job{
		ID:     job.ID,
		Status: job.Status,
		Type:   job.Type,
		Data:   job.Data,

		RetryCount: job.RetryCount,
		Created:    job.Created,
	})

	if err != nil {
		w.Logger.Printf("Error handling job ID %v: %v, retrying", job.ID, err)
		err = w.retryStrategy(w.store, job, err)
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
	w.runningChan <- struct{}{}
	for {
		waitingForAJob, err := w.store.WaitToDoJob()
		if err != nil {
			w.Logger.Println("Error retrieving job: ", err, " trying again")
			time.Sleep(5 * time.Second)
			continue
		}

		w.Logger.Println("Worker waiting", w)
		select {
		case job := <-waitingForAJob:
			if job == nil {
				continue
			}
			w.handleJob(job)

		case <-w.stoppingChan:
			w.Logger.Println("Worker stopping", w)
			w.store.StopWaiting(waitingForAJob)
			w.stoppedChan <- struct{}{}
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
		runningChan:   make(chan struct{}, 1),
		stoppingChan:  make(chan struct{}),
		stoppedChan:   make(chan struct{}),

		Logger: logger,
	}
	go run(&w)
	return &w, nil
}

type RetryStrategy func(gotell.Store, *gotell.Job, error) error

func OneAttempt(store gotell.Store, job *gotell.Job, _ error) error {
	return store.FailedJob(job)
}

func AlwaysRetry(store gotell.Store, job *gotell.Job, _ error) error {
	return store.ReturnJob(job)
}

func RetryUntil(failureLimit int) RetryStrategy {
	return func(store gotell.Store, job *gotell.Job, _ error) error {
		job.RetryCount++
		if job.RetryCount >= failureLimit {
			err := store.FailedJob(job)
			if err != nil {
				return err
			}
			return errors.New("retry limit reached")
		}
		return store.ReturnJob(job)
	}
}
