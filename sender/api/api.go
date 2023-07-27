package api

import (
	gotell "github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/sender"
)

func validateApi(job *gotell.Job) error {
	return nil
}

func genericApiHandler(apiSender sender.ByApi, job *gotell.Job) error {
	err := validateApi(job)
	if err != nil {
		return err
	}
	apiSender.WithBody(job.Data.Body)
	return apiSender.Send()
}

func MakeHandler(apiSender sender.ByApi) gotell.JobHandler {
	return func(job gotell.Job) error {
		return genericApiHandler(apiSender, &job)
	}
}
