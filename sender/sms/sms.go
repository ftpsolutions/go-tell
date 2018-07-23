package sms

import (
	gotell "github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/sender"
)

func validateSMS(job *gotell.Job) error {
	return nil
}

func MakeSMSHandler(smsSender sender.BySMS) gotell.JobHandler {
	return func(job gotell.Job) error {
		err := validateSMS(&job)
		if err != nil {
			return err
		}

		smsSender.From(job.Data.From)
		smsSender.WithBody(job.Data.Body)
		smsSender.To(job.Data.To, job.Data.CC...)

		return smsSender.Send()
	}
}
