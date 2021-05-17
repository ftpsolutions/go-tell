package sms

import (
	"testing"

	"github.com/stretchr/testify/assert"

	gotell "github.com/ftpsolutions/go-tell"
)

const (
	australianMobileNumber                         = "0412345678"
	australianMobileNumberWithCountryCode          = "+61412345678"
	australianMobileNumberWithCountryCodeAlternate = "0061412345678"
)

func TestValidateSms(t *testing.T) {
	job := gotell.Job{
		Type: gotell.JobTypeSMS,
		Data: gotell.JobData{
			To: australianMobileNumber,
			CC: []string{
				australianMobileNumber,
				australianMobileNumberWithCountryCode,
				australianMobileNumberWithCountryCodeAlternate,
			},
		},
	}

	err := validateSMS(&job)

	assert.Nil(t, err)
}
