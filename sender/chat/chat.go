package chat

import (
	"github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/sender"
)

func validateChat(job *gotell.Job) error {
	return nil
}

func MakeHandler(chatSender sender.ByChat) gotell.JobHandler {
	return func(job gotell.Job) error {
		err := validateChat(&job)
		if err != nil {
			return err
		}

		chatSender.WithBody(job.Data.Body)
		chatSender.To(job.Data.To)

		return chatSender.Send()

	}
}
