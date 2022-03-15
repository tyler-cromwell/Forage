package tests

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

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
			router = trelloMocks.MockGetMember(router)
			router = trelloMocks.MockMemberGetBoards(router)
			router = trelloMocks.MockBoardGetLists(router)
			router = trelloMocks.MockListGetCards(router)
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
			router = trelloMocks.MockGetMemberError(router)
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
			router = trelloMocks.MockGetMember(router)
			router = trelloMocks.MockMemberGetBoardsError(router)
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
			router = trelloMocks.MockGetMember(router)
			router = trelloMocks.MockMemberGetBoards(router)
			router = trelloMocks.MockBoardGetListsError(router)
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
			router = trelloMocks.MockGetMember(router)
			router = trelloMocks.MockMemberGetBoards(router)
			router = trelloMocks.MockBoardGetLists(router)
			router = trelloMocks.MockListParamCardsError(router)
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
		t.Run("Success", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = trelloMocks.MockGetMember(router)
			router = trelloMocks.MockMemberGetBoards(router)
			router = trelloMocks.MockBoardGetLists(router)
			router = trelloMocks.MockBoardGetLabels(router)
			router = trelloMocks.MockListAddCards(router)
			router = trelloMocks.MockCardSetPos(router)
			router = trelloMocks.MockCreateChecklist(router)
			//Mock CreateCheckItem endpoint? (shouldn't be needed)
			//Mock AddIDLabel endpoint? (shouldn't be needed)
			server := httptest.NewServer(router)
			defer server.Close()
			client := trelloMocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)
			dueDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			card, err := client.CreateShoppingList(&dueDate, []string{}, []string{})
			require.NoError(t, err)
			require.NotNil(t, card)
		})
	})

	t.Run("AddToShoppingList", func(t *testing.T) {

	})
}
