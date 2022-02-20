package tests

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/tyler-cromwell/forage/clients"
	trelloMocks "github.com/tyler-cromwell/forage/tests/mocks/trello"
)

func TestTrelloClient(t *testing.T) {
	logrus.SetOutput(os.Stdout)

	t.Run("NewTrelloClientWrapper", func(t *testing.T) {
		client := clients.NewTrelloClientWrapper("", "", "", "", "", "")
		require.NotNil(t, client)
	})

	t.Run("GetShoppingList", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = trelloMocks.MockMembersParam(router)
			router = trelloMocks.MockMembersParamBoards(router)
			router = trelloMocks.MockBoardsParamLists(router)
			router = trelloMocks.MockListsParamCards(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := trelloMocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			card, err := client.GetShoppingList()
			require.NoError(t, err)
			require.NotNil(t, card)
		})

		t.Run("ErrorGetMembers", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = trelloMocks.MockMembersParamError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := trelloMocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			card, err := client.GetShoppingList()
			require.Error(t, err)
			require.Nil(t, card)
		})

		t.Run("ErrorGetBoards", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = trelloMocks.MockMembersParam(router)
			router = trelloMocks.MockMembersParamBoardsError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := trelloMocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			card, err := client.GetShoppingList()
			require.Error(t, err)
			require.Nil(t, card)
		})

		t.Run("ErrorGetLists", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = trelloMocks.MockMembersParam(router)
			router = trelloMocks.MockMembersParamBoards(router)
			router = trelloMocks.MockBoardsParamListsError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := trelloMocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			card, err := client.GetShoppingList()
			require.Error(t, err)
			require.Nil(t, card)
		})

		t.Run("ErrorGetCards", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = trelloMocks.MockMembersParam(router)
			router = trelloMocks.MockMembersParamBoards(router)
			router = trelloMocks.MockBoardsParamLists(router)
			router = trelloMocks.MockListsParamCardsError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := trelloMocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			card, err := client.GetShoppingList()
			require.Error(t, err)
			require.Nil(t, card)
		})
	})

	t.Run("CreateShoppingList", func(t *testing.T) {

	})

	t.Run("AddToShoppingList", func(t *testing.T) {

	})
}
