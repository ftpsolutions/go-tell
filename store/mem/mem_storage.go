package memstorage

import (
	"github.com/kithix/go-tell/store"
)

type MemStore struct {
	jobs []*store.Job
}

func (m *MemStore) AddJob(job *store.Job) error {
	job.Status = store.StatusJobCreated
	m.jobs = append(m.jobs, job)
	return nil
}

func (m *MemStore) GetJob() (*store.Job, error) {
	for _, job := range m.jobs {
		if job.Status == store.StatusJobCreated {
			job.Status = store.StatusJobPending
			return job, nil
		}
	}
	return nil, store.ErrorNoJobFound
}

func (m *MemStore) UpdateJob(job *store.Job) error {
	for _, j := range m.jobs {
		if j.ID == job.ID {
			j = job
			return nil
		}
	}
	return store.ErrorNoJobFound
}

func (m *MemStore) DeleteJob(job *store.Job) error {
	for i, j := range m.jobs {
		if j.ID == job.ID {
			m.jobs[i] = m.jobs[len(m.jobs)-1]
			m.jobs[len(m.jobs)-1] = nil
			m.jobs = m.jobs[:len(m.jobs)-1]

			return nil
		}
	}
	return store.ErrorNoJobFound
}

// Creates a new running memstore
func Open() *MemStore {
	return &MemStore{
		make([]*store.Job, 0),
	}
}
