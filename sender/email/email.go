package email

import (
	gotell "github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/sender"
)

func validateEmail(job *gotell.Job) error {
	return nil
}

func GenericEmailHandler(emailSender sender.ByEmail, job *gotell.Job) error {
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

func MakeWithAttachmentsHandler(emailSender sender.ByEmailWithAttachments) gotell.JobHandler {
	return func(job gotell.Job) error {
		emailSender.WithAttachments(job.Data.Attachments)
		return GenericEmailHandler(
			emailSender,
			&job,
		)
	}
}

func MakeHandler(emailSender sender.ByEmail) gotell.JobHandler {
	return func(job gotell.Job) error {
		return GenericEmailHandler(emailSender, &job)
	}
}
