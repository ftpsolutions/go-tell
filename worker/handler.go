package worker

import "github.com/ftpsolutions/go-tell"

type JobHandlers map[string]gotell.JobHandler

func NewMultiJobHandler(handlers ...gotell.JobHandler) gotell.JobHandler {
	return func(job gotell.Job) (err error) {
		for _, handler := range handlers {
			err = handler(job)
			if err != nil {
				return
			}
		}
		return
	}
}
