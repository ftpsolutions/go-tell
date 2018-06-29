package worker

import (
	"testing"

	"github.com/ftpsolutions/go-tell/store/storetest"

	"github.com/ftpsolutions/go-tell/store"
	"github.com/ftpsolutions/go-tell/store/mem"
)

func TestRetryUntil(t *testing.T) {
	s := store.Basic(memstorage.Open())
	retrier := RetryUntil(2)
	job := &store.Job{
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
	if job.Status != store.StatusJobCreated {
		t.Error("Expected job to have created status")
	}
	err = retrier(s, job)
	if err == nil {
		t.Error("Expected error retrying")
	}
	if job.Status != store.StatusJobError {
		t.Error("Expected job to have error status")
	}
}
