package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/adlio/trello"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/tyler-cromwell/forage/clients"
	"github.com/tyler-cromwell/forage/tests/mocks"
)

func TestTrelloClient(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/members/{mid}", func(response http.ResponseWriter, request *http.Request) {
		member := trello.Member{
			ID: "mid",
		}
		b, err := json.Marshal(member)
		if err != nil {
			panic(err)
		}
		response.WriteHeader(http.StatusOK)
		response.Write(b)
	})

	router.HandleFunc("/members/{mid}/boards", func(response http.ResponseWriter, request *http.Request) {
		boards := make([]trello.Board, 1)
		boards[0] = trello.Board{
			ID:   "board",
			Name: "Board",
		}
		b, err := json.Marshal(boards)
		if err != nil {
			panic(err)
		}
		response.WriteHeader(http.StatusOK)
		response.Write(b)
	})

	router.HandleFunc("/boards/{bid}/lists", func(response http.ResponseWriter, request *http.Request) {
		lists := make([]trello.List, 1)
		lists[0] = trello.List{
			ID:   "list",
			Name: "List",
		}
		l, err := json.Marshal(lists)
		if err != nil {
			panic(err)
		}
		response.WriteHeader(http.StatusOK)
		response.Write(l)
	})

	router.HandleFunc("/lists/{lid}/cards", func(response http.ResponseWriter, request *http.Request) {
		cards := make([]trello.Card, 1)
		cards[0] = trello.Card{
			ID:   "shopping_list",
			Name: "Shopping List",
		}
		c, err := json.Marshal(cards)
		if err != nil {
			panic(err)
		}
		response.WriteHeader(http.StatusOK)
		response.Write(c)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	/*
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request, send mock response, etc.
			logrus.WithFields(logrus.Fields{"request": r}).Info()
		})
	*/

	t.Run("NewTrelloClientWrapper", func(t *testing.T) {
		client := clients.NewTrelloClientWrapper("", "", "", "", "", "")
		require.NotNil(t, client)
	})

	t.Run("GetShoppingList", func(t *testing.T) {
		client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
		require.NotNil(t, client)

		card, err := client.GetShoppingList()
		require.NoError(t, err)
		require.NotNil(t, card)
	})

	t.Run("CreateShoppingList", func(t *testing.T) {

	})

	t.Run("AddToShoppingList", func(t *testing.T) {

	})
}
