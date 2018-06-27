package chat

import (
	"github.com/kithix/go-tell/sender"
	"github.com/kithix/go-tell/store"
	"github.com/kithix/go-tell/worker"
)

func validateChat(job *store.Job) error {
	return nil
}

func MakeHandler(chatSender sender.ByChat) worker.JobHandler {
	return func(job store.Job) error {
		err := validateChat(&job)
		if err != nil {
			return err
		}

		chatSender.WithBody(job.Data.Body)
		chatSender.To(job.Data.To)

		return chatSender.Send()

	}
}
