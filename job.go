package gotell

import (
	"errors"

	"github.com/gofrs/uuid"
)

const StatusJobCreated = ""
const StatusJobPending = "pending"
const StatusJobComplete = "complete"
const StatusJobError = "error"

var (
	StaleJobError = errors.New("Job is stale")
)

const (
	JobTypeSMS   = "SMS"
	JobTypeEmail = "Email"
	JobTypeChat  = "Chat"
)

type JobHandler func(job Job) error

// Job stores the state and raw information to turn it into a task to be executed.
type Job struct {
	ID     uuid.UUID
	Type   string // "email", "sms", "chat" , ?
	Status string
	Data   JobData
	// Failure Information
	RetryCount int
}

type JobData struct {
	From        string
	Subject     string
	Body        string
	Tag         string
	Tracking    bool
	HTML        bool
	To          string
	CC          []string
	Attachments map[string][]byte
}

// Human representation of a job.
func (d JobData) String() string {
	emails := d.To
	for _, email := range d.CC {
		emails += ", " + email
	}
	return d.Subject + " - " + emails
}

func generateJobID() uuid.UUID {
	// Panics if unable to generate an ID.
	return uuid.Must(uuid.NewV4())
}

func validateChatJob(body string, to ...string) error {
	return nil
}

func BuildChatJob(body string, to ...string) (*Job, error) {
	err := validateChatJob(body, to...)
	if err != nil {
		return &Job{}, err
	}
	return &Job{
		ID:     generateJobID(),
		Type:   JobTypeChat,
		Status: StatusJobCreated,
		Data: JobData{
			Body: body,
			To:   to[0],
			CC:   to[1:],
		},
	}, nil
}

func validateEmailJob(from, subject, body, tag string, to ...string) error {
	return nil
}

//BuildEmailJob
// from is email
// subject is string max len??
// to is email
func BuildEmailJob(from, subject, body, tag string, tracking, isHTML bool, attachments map[string][]byte, to ...string) (*Job, error) {
	err := validateEmailJob(from, subject, body, tag, to...)
	if err != nil {
		return &Job{}, err
	}

	return &Job{
		ID:     generateJobID(),
		Type:   JobTypeEmail,
		Status: StatusJobCreated,
		Data: JobData{
			From:        from,
			Subject:     subject,
			Body:        body,
			Tag:         tag,
			Tracking:    tracking,
			HTML:        isHTML,
			Attachments: attachments,
			To:          to[0],
			CC:          to[1:],
		},
	}, nil
}

//TODO Test this when we want to implement
func validateSMSJob(from, body string, to ...string) error {
	return nil
}

// BuildSMSJob
// from is a phone number
// body is the content
// to is a list of target phone numbers
func BuildSMSJob(from, body string, to ...string) (Job, error) {
	err := validateSMSJob(from, body, to...)
	if err != nil {
		return Job{}, err
	}
	return Job{
		ID:     generateJobID(),
		Type:   JobTypeSMS,
		Status: StatusJobCreated,
		Data: JobData{
			From: from,
			Body: body,
			To:   to[0],
			CC:   to[1:],
		},
	}, nil
}
