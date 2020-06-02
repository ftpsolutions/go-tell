package store

import (
	"sync"

	"github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/store/mem"
)

type BasicStore struct {
	sync.RWMutex

	storage   gotell.Storage
	receivers *receiverList
}

func Open(storage gotell.Storage) *BasicStore {
	if storage == nil {
		storage = gotell.Storage(mem.Open())
	}
	return &BasicStore{
		storage:   storage,
		receivers: newReceiverList(),
	}
}

// Wrappers for internal storage
func (s *BasicStore) GetJob() (*gotell.Job, error) {
	s.RLock()
	defer s.RUnlock()
	return s.storage.GetJob()
}

func (s *BasicStore) UpdateJob(job *gotell.Job) error {
	s.Lock()
	defer s.Unlock()
	return s.storage.UpdateJob(job)
}

func (s *BasicStore) DeleteJob(job *gotell.Job) error {
	s.Lock()
	defer s.Unlock()
	return s.storage.DeleteJob(job)
}
func (s *BasicStore) AddJob(job *gotell.Job) error {
	s.Lock()
	defer s.Unlock()
	var err error
	job.Status = gotell.StatusJobCreated
	receiver := s.receivers.Get()
	if receiver != nil {
		// Update our job as it's about to be worked on.
		job.Status = gotell.StatusJobPending
		// Setup a defer to ensure the receiver is given the job after the state is handled.
		defer func() {
			if err != nil {
				// Reset our jobs state.
				job.Status = ""
				// Ensure the receiver is returned.
				s.receivers.Add(receiver)
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

func (s *BasicStore) WaitToDoJob() (chan *gotell.Job, error) {
	s.RLock()
	defer s.RUnlock()
	receiver := make(chan *gotell.Job, 1)
	job, err := s.storage.GetJob()
	if err != nil {
		// If there are no jobs, add us to a pending list
		if err == gotell.ErrorNoJobFound {
			s.receivers.Add(receiver)
			return receiver, nil
		}
		// Something else happened, bail.
		return nil, err
	}

	receiver <- job
	close(receiver)
	return receiver, nil
}

// Helper methods
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

type receiverList struct {
	pending []chan *gotell.Job
	sync.Mutex
}

func newReceiverList() *receiverList {
	return &receiverList{
		pending: make([]chan *gotell.Job, 0),
	}
}

func (r *receiverList) Get() chan *gotell.Job {
	r.Lock()
	defer r.Unlock()
	if len(r.pending) > 0 {
		receiver := r.pending[0]
		r.pending = r.pending[1:]
		return receiver
	}
	return nil
}

func (r *receiverList) Add(receiver chan *gotell.Job) {
	r.Lock()
	r.pending = append(r.pending, receiver)
	r.Unlock()
}
