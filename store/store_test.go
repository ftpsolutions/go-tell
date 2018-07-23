package store

import (
	"testing"

	"github.com/ftpsolutions/go-tell/store/storetest"

	"github.com/ftpsolutions/go-tell"
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
	jobChan := s.WaitToDoJob()
	job := <-jobChan
	if job.ID != storetest.BuildJobID(1) {
		t.Errorf("Expected Job.ID:%v, got %v", storetest.BuildJobID(1), job.ID)
	}

	// Only available after being added
	jobChan = s.WaitToDoJob()
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
