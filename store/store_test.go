package store

import (
	"runtime"
	"testing"
	"time"

	"github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/store/storetest"
)

func TestBasicStore_GetJob(t *testing.T) {
	s := Open(nil)
	_, err := s.GetJob()
	if err != gotell.ErrorNoJobFound {
		t.Errorf("Expected error retrieving job: %v, got Error: %v", gotell.ErrorNoJobFound, err)
	}
	err = s.AddJob(&gotell.Job{
		ID: storetest.BuildJobID(1),
	})
	if err != nil {
		t.Errorf("Unexpected error when adding job: %v", err)
	}
	job, err := s.GetJob()
	if err != nil {
		t.Errorf("Unexpected error when adding job: %v", err)
	}
	if job.ID != storetest.BuildJobID(1) {
		t.Errorf("Expected Job.ID:%v, got %v", storetest.BuildJobID(1), job.ID)
	}
	job, err = s.GetJob()
	if err != gotell.ErrorNoJobFound {
		t.Errorf("Expected error retrieving job: %v, got Error: %v", gotell.ErrorNoJobFound, err)
	}
	if job != nil {
		t.Errorf("Found a job: %v, expected none", *job)
	}
}

func TestBasicStore_WaitToDoJob(t *testing.T) {
	s := Open(nil)
	// Immediately available.
	err := s.AddJob(&gotell.Job{
		ID: storetest.BuildJobID(1),
	})
	if err != nil {
		t.Errorf("Unexpected error when adding job: %v", err)
	}
	jobChan, err := s.WaitToDoJob()
	if err != nil {
		t.Errorf("Error getting job channel: %v", err)
	}
	job := <-jobChan
	if job.ID != storetest.BuildJobID(1) {
		t.Errorf("Expected Job.ID:%v, got %v", storetest.BuildJobID(1), job.ID)
	}

	// Only available after being added
	jobChan, err = s.WaitToDoJob()
	if err != nil {
		t.Errorf("Error getting job channel: %v", err)
	}
	err = s.AddJob(&gotell.Job{
		ID: storetest.BuildJobID(2),
	})
	if err != nil {
		t.Errorf("Unexpected error when adding job: %v", err)
	}
	job = <-jobChan
	if job.ID != storetest.BuildJobID(2) {
		t.Errorf("Expected Job.ID:%v, got %v", storetest.BuildJobID(2), job.ID)
	}
}

func TestBasicStore_ReceiverRaceCondition(t *testing.T) {
	mockStorage := &storetest.MockStorage{}
	signal := make(chan struct{})

	// Will catch when you are waiting to 'get a job' as a job is being added.
	s := Open(mockStorage)

	mockStorage.MockAddJob = func(job *gotell.Job) error {
		return nil
	}
	mockStorage.MockGetJob = func() (*gotell.Job, error) {
		// Trigger add job
		signal <- struct{}{}
		runtime.Gosched()
		// Wait for add job to finish
		select {
		// If add job was able to run and finish before the timeout. We lost.
		case <-signal:
		// Time after to escape on locks.
		case <-time.After(1 * time.Nanosecond):
		}
		// In both outcomes no job is found.
		return nil, gotell.ErrorNoJobFound
	}

	// Wait for get job to fire before we start adding our job.
	// Once the job has been added, signal back
	go func() {
		<-signal
		// Add job one,
		s.AddJob(&gotell.Job{
			ID: storetest.BuildJobID(1),
		})
		close(signal)
	}()

	// GetJob will block until the first job has been added.
	jobChan, err := s.WaitToDoJob()
	if err != nil {
		t.Errorf("Error getting job channel: %v", err)
	}

	// Once we have a chan, add job2
	s.AddJob(&gotell.Job{
		ID: storetest.BuildJobID(2),
	})

	// We should receive 1
	job := <-jobChan
	if job.ID != storetest.BuildJobID(1) {
		t.Errorf("Expected Job.ID:%v, got %v", storetest.BuildJobID(1), job.ID)
	}
}
