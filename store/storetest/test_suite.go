package storetest

import (
	"testing"

	"github.com/ftpsolutions/go-tell/store"
)

func errCheck(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func StorageSuite(storeBuilder func() (store.Storage, func()), t *testing.T) {
	// TODO cleanUpFunc may not be run if the test fails at get job. need to confirm
	t.Log("-AddJob-")
	s, cleanUpFunc := storeBuilder()
	AddJob(s, t)
	cleanUpFunc()

	t.Log("-GetJob-")
	s, cleanUpFunc = storeBuilder()
	GetJob(s, t)
	cleanUpFunc()

	t.Log("-UpdateJob-")
	s, cleanUpFunc = storeBuilder()
	UpdateJob(s, t)
	cleanUpFunc()

	t.Log("-DeleteJob-")
	s, cleanUpFunc = storeBuilder()
	DeleteJob(s, t)
	cleanUpFunc()
}

func AddJob(s store.Storage, t *testing.T) {
	errCheck(s.AddJob(&store.Job{
		ID: BuildJobID(1),
	}), t)
	// TODO Should the job with the same ID be able to be 'added' twice?
}

func DeleteJob(s store.Storage, t *testing.T) {
	errCheck(s.AddJob(&store.Job{
		ID: BuildJobID(1),
	}), t)
	errCheck(s.AddJob(&store.Job{
		ID: BuildJobID(2),
	}), t)
	errCheck(s.DeleteJob(&store.Job{
		ID: BuildJobID(1),
	}), t)
	errCheck(s.DeleteJob(&store.Job{
		ID: BuildJobID(2),
	}), t)
	_, err := s.GetJob()
	if err != store.ErrorNoJobFound {
		t.Log(err)
		t.Error("Incorrect error returned")
	}
}

func GetJob(s store.Storage, t *testing.T) {
	errCheck(s.AddJob(&store.Job{
		ID: BuildJobID(1),
	}), t)

	job, err := s.GetJob()
	errCheck(err, t)
	if job.ID != BuildJobID(1) {
		t.Error("Incorrect job retrieved")
	}
	if job.Status != store.StatusJobPending {
		t.Error("Job status not set to pending")
	}
	// Should not be able to get anymore jobs as 1/1 is pending.
	job, err = s.GetJob()
	if err != store.ErrorNoJobFound {
		t.Log(err)
		t.Error("Incorrect error returned")
	}
}

func UpdateJob(s store.Storage, t *testing.T) {
	errCheck(s.AddJob(&store.Job{
		ID: BuildJobID(1),
	}), t)
	errCheck(s.AddJob(&store.Job{
		ID: BuildJobID(2),
	}), t)
}
