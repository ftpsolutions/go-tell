package store

import (
	"sync"

	"github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/store/mem"
)

type BasicStore struct {
	sync.Mutex

	storage          gotell.Storage
	pendingReceivers []chan *gotell.Job
}

func Open(storage gotell.Storage) *BasicStore {
	if storage == nil {
		storage = gotell.Storage(mem.Open())
	}
	return &BasicStore{
		storage:          storage,
		pendingReceivers: make([]chan *gotell.Job, 0),
	}
}

// Wrappers for internal storage
func (s *BasicStore) GetJob() (*gotell.Job, error)    { return s.storage.GetJob() }
func (s *BasicStore) UpdateJob(job *gotell.Job) error { return s.storage.UpdateJob(job) }
func (s *BasicStore) DeleteJob(job *gotell.Job) error { return s.storage.DeleteJob(job) }
func (s *BasicStore) AddJob(job *gotell.Job) error {
	var err error
	job.Status = gotell.StatusJobCreated
	receiver := s.getReceiver()
	if receiver != nil {
		// Update our job as it's about to be worked on.
		job.Status = gotell.StatusJobPending
		// Setup a defer to ensure we handle error state for processing receiver.
		defer func() {
			if err != nil {
				// Reset our jobs state.
				job.Status = ""
				// Ensure the receiver is returned.
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

func (s *BasicStore) WaitToDoJob() chan *gotell.Job {
	receiver := make(chan *gotell.Job, 1)

	job, err := s.storage.GetJob()
	if err != nil {
		s.addReceiver(receiver)
		return receiver
	}

	receiver <- job
	close(receiver)
	return receiver
}

func (s *BasicStore) CompleteJob(job *gotell.Job) error {
	job.Status = gotell.StatusJobComplete
	return s.storage.UpdateJob(job)
}

func (s *BasicStore) FailedJob(job *gotell.Job) error {
	job.Status = gotell.StatusJobError
	return s.storage.UpdateJob(job)
}

func (s *BasicStore) ReturnJob(job *gotell.Job) error {
	job.Status = gotell.StatusJobCreated
	return s.storage.UpdateJob(job)
}

func (s *BasicStore) addReceiver(receiver chan *gotell.Job) {
	s.Lock()
	s.pendingReceivers = append(s.pendingReceivers, receiver)
	s.Unlock()
}

func (s *BasicStore) getReceiver() chan *gotell.Job {
	s.Lock()
	defer s.Unlock()
	if len(s.pendingReceivers) > 0 {
		receiver := s.pendingReceivers[0]
		s.pendingReceivers = s.pendingReceivers[1:]
		return receiver
	}
	return nil
}