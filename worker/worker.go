package worker

import (
	"io/ioutil"
	"log"

	"github.com/kithix/go-tell/store"
)

var nullLogger = log.New(ioutil.Discard, "", 0)

type Job = store.Job
type Store = store.Store

type Worker struct {
	jobHandler JobHandler
	store      Store
	stopChan   chan struct{}

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
		// TODO insert retry strategy here
		return
	}
	w.Logger.Println("Job completed", job)
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

// Takes a store, jobHandler, logger
func Open(store Store, jobHandler JobHandler, logger *log.Logger) (*Worker, error) {
	if logger == nil {
		logger = nullLogger
	}
	stopChan := make(chan struct{})
	w := Worker{
		jobHandler: jobHandler,
		store:      store,
		stopChan:   stopChan,

		Logger: logger,
	}
	go run(&w)
	return &w, nil
}
