package worker

import (
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
	err = retrier(s, job)
	if err != nil {
		t.Error(err)
	}
	if job.Status != gotell.StatusJobCreated {
		t.Error("Expected job to have created status")
	}
	err = retrier(s, job)
	if err == nil {
		t.Error("Expected error retrying")
	}
	if job.Status != gotell.StatusJobError {
		t.Error("Expected job to have error status")
	}
}
