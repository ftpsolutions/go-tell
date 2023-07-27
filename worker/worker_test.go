package worker

import (
	"log"
	"testing"

	gotell "github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/store"
	"github.com/ftpsolutions/go-tell/store/mem"
	"github.com/ftpsolutions/go-tell/store/storetest"
)

func TestRetryUntil(t *testing.T) {
	s := store.Open(mem.Open())
	retrier := RetryUntil(2)
	job := &gotell.Job{
		ID: storetest.BuildJobID(1),
	}
	err := s.AddJob(job)
	if err != nil {
		t.Error(err)
	}
	err = retrier(s, job, err)
	if err != nil {
		t.Error(err)
	}
	if job.Status != gotell.StatusJobCreated {
		t.Error("Expected job to have created status")
	}
	err = retrier(s, job, err)
	if err == nil {
		t.Error("Expected error retrying")
	}
	if job.Status != gotell.StatusJobError {
		t.Error("Expected job to have error status")
	}
}

func TestOneWorkerCanReceiveAJob(t *testing.T) {
	s := store.Open(mem.Open())

	var processedJob *gotell.Job
	handler := func(job gotell.Job) error {
		processedJob = &job
		return nil
	}

	log.Println("Starting worker")
	worker1, _ := Open(s, handler, nil, nil)

	worker1.WaitTillRunning()

	job := &gotell.Job{
		ID: storetest.BuildJobID(1),
	}
	s.AddJob(job)

	log.Println("Closing worker1")
	worker1.Close()

	//worker should have done the thing now
	if processedJob == nil {
		t.Fatal("Job should have been processed")
	}
}

func TestTwoWorkersCanReceiveJobs(t *testing.T) {
	s := store.Open(mem.Open())

	var processedJob *gotell.Job
	handler := func(job gotell.Job) error {
		processedJob = &job
		return nil
	}

	log.Println("Starting worker1")
	worker1, _ := Open(s, handler, nil, nil)
	worker1.WaitTillRunning()

	log.Println("Starting worker2")
	worker2, _ := Open(s, handler, nil, nil)
	worker2.WaitTillRunning()

	job := &gotell.Job{
		ID: storetest.BuildJobID(1),
	}
	s.AddJob(job)

	log.Println("Closing workers")
	worker1.Close()
	worker2.Close()

	//a worker should have done the thing now
	if processedJob == nil {
		t.Fatal("Job should have been processed")
	}
}

func TestClosingWorkerDoesntBreakThings(t *testing.T) {
	s := store.Open(mem.Open())

	var processedJob *gotell.Job
	handler := func(job gotell.Job) error {
		processedJob = &job
		return nil
	}

	log.Println("Starting worker1")
	worker1, _ := Open(s, handler, nil, nil)

	worker1.WaitTillRunning()

	log.Println("Closing worker1")
	worker1.Close()

	log.Println("Starting worker2")
	worker2, _ := Open(s, handler, nil, nil)

	worker2.WaitTillRunning()

	job := &gotell.Job{
		ID: storetest.BuildJobID(1),
	}
	s.AddJob(job)

	log.Println("Closing worker2")
	worker2.Close()

	//worker should have done the thing now
	if processedJob == nil {
		t.Fatal("Job should have been processed")
	}
}
