package storetest

import (
	"github.com/ftpsolutions/go-tell"
)

// MockStorage is a simple struct for mocking internal storage behaviour
type MockStorage struct {
	MockAddJob    func(job *gotell.Job) error
	MockGetJob    func() (*gotell.Job, error)
	MockUpdateJob func(job *gotell.Job) error
	MockDeleteJob func(job *gotell.Job) error
}

//AddJob '
func (s *MockStorage) AddJob(job *gotell.Job) error {
	return s.MockAddJob(job)
}

//GetJob '
func (s *MockStorage) GetJob() (*gotell.Job, error) {
	return s.MockGetJob()
}

//UpdateJob '
func (s *MockStorage) UpdateJob(job *gotell.Job) error {
	return s.MockUpdateJob(job)
}

//DeleteJob '
func (s *MockStorage) DeleteJob(job *gotell.Job) error {
	return s.MockDeleteJob(job)
}
