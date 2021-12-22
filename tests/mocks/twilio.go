package mocks

import (
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/tyler-cromwell/forage/clients"
)

type MockTwilioStruct struct{}

func NewTwilioClientWrapper(accountSid, authToken, phoneFrom, phoneTo string) *clients.Twilio {
	c := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: accountSid,
		Password: authToken,
	})
	ms := MockTwilioStruct{}
	client := clients.Twilio{
		From:      phoneFrom,
		To:        phoneTo,
		Client:    c,
		Interface: &ms,
	}
	return &client
}

func (ms *MockTwilioStruct) CreateMessage(params *openapi.CreateMessageParams) (*openapi.ApiV2010Message, error) {
	var sid string
	var err error

	if *params.Body == "Error" {
		err = fmt.Errorf("Error")
	} else {
		sid = "Hello"
	}

	return &openapi.ApiV2010Message{
		Sid: &sid,
	}, err
}
