package gotell

type Storage interface {
	AddJob(job *Job) error
	GetJob() (*Job, error)
	UpdateJob(job *Job) error
	DeleteJob(job *Job) error
	DeleteOldestJob(jobType string) error
}

type Store interface {
	Storage
	WaitToDoJob() (chan *Job, error)
	StopWaiting(chan *Job)
	CompleteJob(job *Job) error
	ReturnJob(job *Job) error
	FailedJob(job *Job) error
}
