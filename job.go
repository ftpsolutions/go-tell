package gotell

import (
	"time"

	"github.com/gofrs/uuid"
)

const StatusJobCreated = ""
const StatusJobPending = "pending"
const StatusJobComplete = "complete"
const StatusJobError = "error"

const (
	JobTypeSMS   = "SMS"
	JobTypeEmail = "Email"
	JobTypeChat  = "Chat"
	JobTypeApi   = "Api"
)

type JobHandler func(job Job) error

// Job stores the state and raw information to turn it into a task to be executed.
type Job struct {
	ID        uuid.UUID
	Type      string // "email", "sms", "chat" , ?
	Status    string
	RelatedId string
	Data      JobData

	// Failure Information
	RetryCount int
	Created    time.Time
}

func (job Job) Clone() Job {
	return Job{
		ID:        job.ID,
		Type:      job.Type,
		Status:    job.Status,
		RelatedId: job.RelatedId,
		Data:      job.Data,

		RetryCount: job.RetryCount,
		Created:    job.Created,
	}
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

func BuildChatJob(body string, relatedId string, to ...string) (*Job, error) {
	err := validateChatJob(body, to...)
	if err != nil {
		return &Job{}, err
	}
	return &Job{
		ID:        generateJobID(),
		Type:      JobTypeChat,
		Status:    StatusJobCreated,
		RelatedId: relatedId,
		Data: JobData{
			Body: body,
			To:   to[0],
			CC:   to[1:],
		},
		Created: time.Now(),
	}, nil
}

func validateEmailJob(from, subject, body, tag string, to ...string) error {
	return nil
}

// BuildEmailJob
// from is email
// subject is string max len??
// to is email
func BuildEmailJob(from, subject, body, tag string, tracking, isHTML bool, attachments map[string][]byte, relatedId string, to ...string) (*Job, error) {
	err := validateEmailJob(from, subject, body, tag, to...)
	if err != nil {
		return &Job{}, err
	}

	return &Job{
		ID:        generateJobID(),
		Type:      JobTypeEmail,
		Status:    StatusJobCreated,
		RelatedId: relatedId,
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
		Created: time.Now(),
	}, nil
}

// TODO Test this when we want to implement
func validateSMSJob(from, body string, to ...string) error {
	return nil
}

// BuildSMSJob
// from is a phone number
// body is the content
// to is a list of target phone numbers
func BuildSMSJob(from, body string, relatedId string, to ...string) (Job, error) {
	err := validateSMSJob(from, body, to...)
	if err != nil {
		return Job{}, err
	}
	return Job{
		ID:        generateJobID(),
		Type:      JobTypeSMS,
		Status:    StatusJobCreated,
		RelatedId: relatedId,
		Data: JobData{
			From: from,
			Body: body,
			To:   to[0],
			CC:   to[1:],
		},
		Created: time.Now(),
	}, nil
}

func BuildAPIJob(payload string, relatedId string) (Job, error) {
	return Job{
		ID:        generateJobID(),
		Type:      JobTypeApi,
		Status:    StatusJobCreated,
		RelatedId: relatedId,
		Data: JobData{
			Body: payload,
		},
		Created: time.Now(),
	}, nil
}
