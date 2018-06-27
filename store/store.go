package store

import (
	"errors"
	"sync"
)

var (
	ErrorNoJobFound = errors.New("No job found")
)

type Storage interface {
	AddJob(job *Job) error
	GetJob() (*Job, error)
	UpdateJob(job *Job) error
	DeleteJob(job *Job) error
}

type Store interface {
	Storage
	WaitToDoJob() chan *Job
	CompleteJob(job *Job) error
	ReturnJob(job *Job) error
	FailedJob(job *Job) error
}

type BasicStore struct {
	storage Storage

	pendingLock      sync.Mutex
	pendingReceivers []chan *Job
}

func Basic(storage Storage) *BasicStore {
	return &BasicStore{
		storage:          storage,
		pendingReceivers: make([]chan *Job, 0),
	}
}

// Wrappers for internal storage
func (s *BasicStore) GetJob() (*Job, error)    { return s.storage.GetJob() }
func (s *BasicStore) UpdateJob(job *Job) error { return s.storage.UpdateJob(job) }
func (s *BasicStore) DeleteJob(job *Job) error { return s.storage.DeleteJob(job) }
func (s *BasicStore) AddJob(job *Job) error {
	var err error
	job.Status = StatusJobCreated
	receiver := s.getReceiver()
	if receiver != nil {
		// Update our job as it's about to be worked on.
		job.Status = StatusJobPending
		// After adding our job we check for failures
		defer func() {
			if err != nil {
				// Reset our state
				job.Status = ""
				// Ensure the receiver is returned
				s.addReceiver(receiver)
				return
			}
			receiver <- job
			close(receiver)
		}()
	}

	err = s.storage.AddJob(job)
	if err != nil {
		return err
	}
	return nil
}

func (s *BasicStore) WaitToDoJob() chan *Job {
	receiver := make(chan *Job, 1)

	job, err := s.storage.GetJob()
	if err != nil {
		s.addReceiver(receiver)
		return receiver
	}

	receiver <- job
	close(receiver)
	return receiver
}

func (s *BasicStore) CompleteJob(job *Job) error {
	job.Status = StatusJobComplete
	return s.storage.UpdateJob(job)
}

func (s *BasicStore) FailedJob(job *Job) error {
	job.Status = StatusJobError
	return s.storage.UpdateJob(job)
}

func (s *BasicStore) ReturnJob(job *Job) error {
	job.Status = StatusJobCreated
	return s.storage.UpdateJob(job)
}

func (s *BasicStore) addReceiver(receiver chan *Job) {
	s.pendingLock.Lock()
	s.pendingReceivers = append(s.pendingReceivers, receiver)
	s.pendingLock.Unlock()
}

func (s *BasicStore) getReceiver() chan *Job {
	s.pendingLock.Lock()
	defer s.pendingLock.Unlock()
	if len(s.pendingReceivers) > 0 {
		receiver := s.pendingReceivers[0]
		s.pendingReceivers = s.pendingReceivers[1:]
		return receiver
	}
	return nil
}
