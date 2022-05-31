package clients

import (
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioInterface interface {
	CreateMessage(*openapi.CreateMessageParams) (*openapi.ApiV2010Message, error)
}

type Twilio struct {
	From      string
	To        string
	Client    *twilio.RestClient
	Interface TwilioInterface
}

func NewTwilioClientWrapper(accountSid, authToken, phoneFrom, phoneTo string) *Twilio {
	c := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: accountSid,
		Password: authToken,
	})
	client := Twilio{
		From:      phoneFrom,
		To:        phoneTo,
		Client:    c,
		Interface: c.ApiV2010, // Need this in order to mock SMS messaging
	}
	return &client
}

func (tc *Twilio) InnerStruct() *Twilio {
	return tc
}

func (tc *Twilio) ComposeMessage(quantity, quantityExpired int, url string) string {
	var message string
	if quantity == 1 {
		message = fmt.Sprintf("%d item expiring soon and %d already expired! View shopping list: %s", quantity, quantityExpired, url)
	} else if quantity > 1 {
		message = fmt.Sprintf("%d items expiring soon and %d already expired! View shopping list: %s", quantity, quantityExpired, url)
	} else if quantity <= 0 && quantityExpired == 1 {
		message = fmt.Sprintf("%d item expired! View shopping list: %s", quantityExpired, url)
	} else if quantity <= 0 && quantityExpired > 1 {
		message = fmt.Sprintf("%d items expired! View shopping list: %s", quantityExpired, url)
	}
	return message
}

func (tc *Twilio) SendMessage(phoneFrom, phoneTo, message string) (string, error) {
	// Prepare message
	params := &openapi.CreateMessageParams{}
	params.SetFrom(phoneFrom)
	params.SetTo(phoneTo)
	params.SetBody(message)

	// Send it
	//resp, err := tc.Client.ApiV2010.CreateMessage(params)
	resp, err := tc.Interface.CreateMessage(params)
	if err != nil {
		return "", err
	} else {
		return *resp.Sid, err
	}
}
