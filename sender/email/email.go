package email

import (
	"github.com/kithix/go-tell/sender"
	"github.com/kithix/go-tell/store"
	"github.com/kithix/go-tell/worker"
)

func validateEmail(job *store.Job) error {
	return nil
}

func genericEmailHandler(emailSender sender.ByEmail, job *store.Job) error {
	err := validateEmail(job)
	if err != nil {
		return err
	}

	if job.Data.From != "" {
		emailSender.From(job.Data.From)
	}

	emailSender.WithSubject(job.Data.Subject)
	emailSender.WithBody(job.Data.Body)
	emailSender.WithTag(job.Data.Tag)
	emailSender.To(job.Data.To, job.Data.CC...)

	if !job.Data.Tracking {
		emailSender.WithNoTracking()
	}

	if job.Data.HTML {
		return emailSender.SendHtml()
	}
	return emailSender.Send()
}

func MakeWithAttachmentsHandler(emailSender sender.ByEmailWithAttachments) worker.JobHandler {
	return func(job store.Job) error {
		emailSender.WithAttachments(job.Data.Attachments)
		return genericEmailHandler(
			emailSender,
			&job,
		)
	}
}

func MakeHandler(emailSender sender.ByEmail) worker.JobHandler {
	return func(job store.Job) error {
		return genericEmailHandler(emailSender, &job)
	}
}
