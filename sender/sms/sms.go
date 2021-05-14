package sms

import (
	"errors"
	"fmt"
	gotell "github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/sender"
	"regexp"
	"strings"
)

func validateSMS(job *gotell.Job) error {
	re, _ := regexp.Compile(`^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)
	numbers := append([]string{job.Data.To}, job.Data.CC...)
	failedNumbers := []string{}
	for _, s := range numbers {
		if !re.MatchString(s) {
			failedNumbers = append(failedNumbers, s)
		}
	}

	if len(failedNumbers) > 0 {
		return errors.New(fmt.Sprintf("Mobile phone numbers do not match regex. %s", strings.Join(failedNumbers, ", ")))
	}
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
