package sms

import (
	"testing"

	"github.com/stretchr/testify/assert"

	gotell "github.com/ftpsolutions/go-tell"
)

const (
	australianMobileNumber                         = "04 12345678"
	australianMobileNumberWithCountryCode          = "+61412345678"
	australianMobileNumberWithCountryCodeAlternate = "0061412345678"
	americanMobileNumber                           = "+1 (123) 456–7890"
	britishMobileNumber                            = "07890 123456"
	britishMobileNumberWithCountryCode             = "+447890123456"
	thaiMobileNumber                               = "+66 2-123-4567"
	czechMobileNumber                              = "+420 2 / 12 34 56 78"
)

func TestValidateSms(t *testing.T) {
	job := gotell.Job{
		Type: gotell.JobTypeSMS,
		Data: gotell.JobData{
			To: australianMobileNumber,
			CC: []string{
				australianMobileNumberWithCountryCode,
				australianMobileNumberWithCountryCodeAlternate,
				americanMobileNumber,
				britishMobileNumber,
				britishMobileNumberWithCountryCode,
				thaiMobileNumber,
				czechMobileNumber,
			},
		},
	}

	to, cc := transformToForJob(job.Data)
	job.Data.To = to
	job.Data.CC = cc
	err := validateSMS(&job)

	assert.Nil(t, err)
}

func TestJobTransform(t *testing.T) {
	data := gotell.JobData{
		To: australianMobileNumber,
		CC: []string{
			australianMobileNumberWithCountryCode,
			australianMobileNumberWithCountryCodeAlternate,
			americanMobileNumber,
			britishMobileNumber,
			britishMobileNumberWithCountryCode,
			thaiMobileNumber,
			czechMobileNumber,
		},
	}

	to, result := transformToForJob(data)

	assert.Equal(t, to, "0412345678")
	assert.Equal(t, result, []string{"+61412345678", "0061412345678", "+11234567890", "07890123456", "+447890123456", "+6621234567", "+420212345678"})
}
