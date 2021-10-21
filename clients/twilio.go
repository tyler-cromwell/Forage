package clients

import (
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var twilioClient *twilio.RestClient

func init() {
	// Get Twilio credentials
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	// Initialize Twilio client
	twilioClient = twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: accountSid,
		Password: authToken,
	})
}

func SendMessage(phoneFrom, phoneTo, message string) (string, error) {
	// Prepare message
	params := &openapi.CreateMessageParams{}
	params.SetFrom(phoneFrom)
	params.SetTo(phoneTo)
	params.SetBody(message)

	// Send it
	resp, err := twilioClient.ApiV2010.CreateMessage(params)
	if err != nil {
		return "", err
	} else {
		return *resp.Sid, err
	}
}
