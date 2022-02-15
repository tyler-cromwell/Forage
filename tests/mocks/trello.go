package mocks

import (
	"net/http/httptest"

	"github.com/adlio/trello"
	"github.com/tyler-cromwell/forage/clients"
)

func NewTrelloClientWrapper(mockServer *httptest.Server, apiKey, apiToken, memberID, boardName, listName, labels string) *clients.Trello {
	client := clients.Trello{
		Key:       apiKey,
		Token:     apiToken,
		MemberID:  memberID,
		BoardName: boardName,
		ListName:  listName,
		LabelsStr: labels,
		Client:    trello.NewClient(apiKey, apiToken),
	}

	client.Client.BaseURL = mockServer.URL // Something like "http://127.0.0.1:53791"
	return &client
}
