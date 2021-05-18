package sms

import (
	"errors"
	"fmt"
	gotell "github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/sender"
	"regexp"
	"strings"
)

// https://stackoverflow.com/questions/123559/how-to-validate-phone-numbers-using-regex?page=1&tab=votes#tab-top
var reStrip = regexp.MustCompile(`[^\+\d]+`)
var re = regexp.MustCompile(`^(\+?\d{9,15})$`)

func stripPhoneNumberCharacters(number string) string {
	number = reStrip.ReplaceAllString(number, "")
	return number
}

func validateSMS(job *gotell.Job) error {
	numbers := append([]string{job.Data.To}, job.Data.CC...)
	var failedNumbers []string
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

func transformToForJob(data gotell.JobData) (string, []string) {
	var results []string
	for _, number := range data.CC {
		results = append(results, stripPhoneNumberCharacters(number))
	}
	return stripPhoneNumberCharacters(data.To), results
}

func MakeSMSHandler(smsSender sender.BySMS) gotell.JobHandler {
	return func(job gotell.Job) error {
		to, cc := transformToForJob(job.Data)
		job.Data.To = to
		job.Data.CC = cc
		fmt.Println(job.Data.To)
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
