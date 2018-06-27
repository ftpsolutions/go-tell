package worker

import (
	"github.com/kithix/go-tell/store"
)

type JobHandler func(job Job) error

type JobHandlers map[string]JobHandler

func NewMultiJobHandler(handlers ...JobHandler) JobHandler {
	return func(job store.Job) (err error) {
		for _, handler := range handlers {
			err = handler(job)
			if err != nil {
				return
			}
		}
		return
	}
}
