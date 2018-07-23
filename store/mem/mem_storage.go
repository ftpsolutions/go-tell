package mem

import gotell "github.com/ftpsolutions/go-tell"

type MemStorage struct {
	jobs []*gotell.Job
}

func (m *MemStorage) AddJob(job *gotell.Job) error {
	job.Status = gotell.StatusJobCreated
	m.jobs = append(m.jobs, job)
	return nil
}

func (m *MemStorage) GetJob() (*gotell.Job, error) {
	for _, job := range m.jobs {
		if job.Status == gotell.StatusJobCreated {
			job.Status = gotell.StatusJobPending
			return job, nil
		}
	}
	return nil, gotell.ErrorNoJobFound
}

func (m *MemStorage) UpdateJob(job *gotell.Job) error {
	for _, j := range m.jobs {
		if j.ID == job.ID {
			j = job
			return nil
		}
	}
	return gotell.ErrorNoJobFound
}

func (m *MemStorage) DeleteJob(job *gotell.Job) error {
	for i, j := range m.jobs {
		if j.ID == job.ID {
			m.jobs[i] = m.jobs[len(m.jobs)-1]
			m.jobs[len(m.jobs)-1] = nil
			m.jobs = m.jobs[:len(m.jobs)-1]

			return nil
		}
	}
	return gotell.ErrorNoJobFound
}

// Creates a new running MemStorage
func Open() *MemStorage {
	return &MemStorage{
		make([]*gotell.Job, 0),
	}
}
