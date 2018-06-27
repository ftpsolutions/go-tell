package storetest

import (
	"github.com/kithix/go-tell/store"
)

type MockStorage struct {
	jobWasAdded bool
	fakeChan    chan *store.Job
	AddedJobs   []*store.Job
}

func (s *MockStorage) AddJob(job *store.Job) error {
	s.jobWasAdded = true
	s.AddedJobs = append(s.AddedJobs, job)
	return nil
}

func (s *MockStorage) JobWasAdded() bool {
	return s.jobWasAdded
}

func (s *MockStorage) WaitToDoJob() chan *store.Job {
	return s.fakeChan
}

func (s *MockStorage) JobReady(id uint8) {
	go func() {
		s.fakeChan <- &store.Job{ID: BuildJobID(id)}
	}()
}

func (s *MockStorage) CompleteJob(job *store.Job) error {
	return nil
}

func (s *MockStorage) FailedJob(job *store.Job) error {
	return nil
}

func (s *MockStorage) ReturnJob(job *store.Job) error {
	return nil
}

func (s *MockStorage) DeleteJob(job *store.Job) error {
	return nil
}

func (s MockStorage) Open() *MockStorage {
	return &MockStorage{
		fakeChan: make(chan *store.Job),
	}
}
