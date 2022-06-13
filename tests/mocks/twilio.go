package mocks

import (
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/tyler-cromwell/forage/clients"
)

type MockTwilio struct {
	From                   string
	To                     string
	OverrideComposeMessage func(int, int, string) string
	OverrideSendMessage    func(string, string, string) (string, error)
}

type MockTwilioStruct struct{}

func (mtc *MockTwilio) ComposeMessage(quantity, quantityExpired int, url string) string {
	if mtc.OverrideComposeMessage != nil {
		return mtc.OverrideComposeMessage(quantity, quantityExpired, url)
	} else {
		return ""
	}
}

func (mtc *MockTwilio) SendMessage(phoneFrom, phoneTo, message string) (string, error) {
	if mtc.OverrideSendMessage != nil {
		return mtc.OverrideSendMessage(phoneFrom, phoneTo, message)
	} else {
		return "", nil
	}
}

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
