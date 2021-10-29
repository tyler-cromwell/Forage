package clients

import (
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Twilio struct {
	Client *twilio.RestClient
}

func (twilio *Twilio) SendMessage(phoneFrom, phoneTo, message string) (string, error) {
	// Prepare message
	params := &openapi.CreateMessageParams{}
	params.SetFrom(phoneFrom)
	params.SetTo(phoneTo)
	params.SetBody(message)

	// Send it
	resp, err := twilio.Client.ApiV2010.CreateMessage(params)
	if err != nil {
		return "", err
	} else {
		return *resp.Sid, err
	}
}
