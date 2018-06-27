package chat

import (
	"github.com/ftpsolutions/go-tell/sender"
	"github.com/ftpsolutions/go-tell/store"
	"github.com/ftpsolutions/go-tell/worker"
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
