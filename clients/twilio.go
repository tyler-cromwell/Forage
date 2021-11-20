package clients

import (
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Twilio struct {
	From   string
	To     string
	Client *twilio.RestClient
}

func NewTwilioClientWrapper(accountSid, authToken, phoneFrom, phoneTo string) *Twilio {
	client := Twilio{
		From: phoneFrom,
		To:   phoneTo,
		Client: twilio.NewRestClientWithParams(twilio.RestClientParams{
			Username: accountSid,
			Password: authToken,
		}),
	}
	return &client
}

func (tc *Twilio) SendMessage(phoneFrom, phoneTo, message string) (string, error) {
	// Prepare message
	params := &openapi.CreateMessageParams{}
	params.SetFrom(phoneFrom)
	params.SetTo(phoneTo)
	params.SetBody(message)

	// Send it
	resp, err := tc.Client.ApiV2010.CreateMessage(params)
	if err != nil {
		return "", err
	} else {
		return *resp.Sid, err
	}
}
