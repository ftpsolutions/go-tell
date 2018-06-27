package sms

import (
	"github.com/kithix/go-tell/sender"
	"github.com/kithix/go-tell/store"
	"github.com/kithix/go-tell/worker"
)

func validateSMS(job *store.Job) error {
	return nil
}

func MakeSMSHandler(smsSender sender.BySMS) worker.JobHandler {
	return func(job store.Job) error {
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
