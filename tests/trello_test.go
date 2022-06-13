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
	"github.com/tyler-cromwell/forage/tests/mocks"
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
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockListGetCards(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			card, err := client.GetShoppingList()
			require.NoError(t, err)
			require.NotNil(t, card)
		})

		t.Run("ErrorGetMembers", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMemberError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			card, err := client.GetShoppingList()
			require.Error(t, err)
			require.Nil(t, card)
		})

		t.Run("ErrorGetBoards", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoardsError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			card, err := client.GetShoppingList()
			require.Error(t, err)
			require.Nil(t, card)
		})

		t.Run("ErrorGetLists", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetListsError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			card, err := client.GetShoppingList()
			require.Error(t, err)
			require.Nil(t, card)
		})

		t.Run("ErrorGetCards", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockListParamCardsError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			card, err := client.GetShoppingList()
			require.Error(t, err)
			require.Nil(t, card)
		})
	})

	t.Run("CreateShoppingList", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockBoardGetLabels(router)
			router = mocks.MockListAddCards(router)
			router = mocks.MockCardSetPos(router)
			router = mocks.MockCreateChecklist(router)
			router = mocks.MockCreateCheckItem(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			dueDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			card, err := client.CreateShoppingList(&dueDate, []string{"label"}, []string{})
			require.NoError(t, err)
			require.NotNil(t, card)
		})

		t.Run("ErrorGetMember", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMemberError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			dueDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			card, err := client.CreateShoppingList(&dueDate, []string{"label"}, []string{})
			require.Error(t, err)
			require.Empty(t, card)
		})

		t.Run("ErrorMemberGetBoards", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoardsError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			dueDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			card, err := client.CreateShoppingList(&dueDate, []string{"label"}, []string{})
			require.Error(t, err)
			require.Empty(t, card)
		})

		t.Run("ErrorBoardGetLists", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetListsError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			dueDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			card, err := client.CreateShoppingList(&dueDate, []string{"label"}, []string{})
			require.Error(t, err)
			require.Empty(t, card)
		})

		t.Run("ErrorBoardGetLabels", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockBoardGetLabelsError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			dueDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			card, err := client.CreateShoppingList(&dueDate, []string{"label"}, []string{})
			require.Error(t, err)
			require.Empty(t, card)
		})

		t.Run("ErrorListParamCards", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockBoardGetLabels(router)
			router = mocks.MockListParamCardsError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			dueDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			card, err := client.CreateShoppingList(&dueDate, []string{"Label"}, []string{})
			require.Error(t, err)
			require.Empty(t, card)
		})

		t.Run("ErrorCardsParam", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockBoardGetLabels(router)
			router = mocks.MockListAddCards(router)
			router = mocks.MockCardsParamError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			dueDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			card, err := client.CreateShoppingList(&dueDate, []string{"Label"}, []string{})
			require.Error(t, err)
			require.Empty(t, card)
		})

		t.Run("ErrorCreateChecklist", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockBoardGetLabels(router)
			router = mocks.MockListAddCards(router)
			router = mocks.MockCardSetPos(router)
			router = mocks.MockCreateChecklistError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			dueDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			card, err := client.CreateShoppingList(&dueDate, []string{"Label"}, []string{})
			require.Error(t, err)
			require.Empty(t, card)
		})

		t.Run("ErrorCreateCheckItem", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockBoardGetLabels(router)
			router = mocks.MockListAddCards(router)
			router = mocks.MockCardSetPos(router)
			router = mocks.MockCreateChecklist(router)
			router = mocks.MockCreateCheckItemError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			dueDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			card, err := client.CreateShoppingList(&dueDate, []string{"Label"}, []string{"ListItem"})
			require.Error(t, err)
			require.Empty(t, card)
		})

		t.Run("ErrorAddIDLabel", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockBoardGetLabels(router)
			router = mocks.MockListAddCards(router)
			router = mocks.MockCardSetPos(router)
			router = mocks.MockCreateChecklist(router)
			router = mocks.MockCreateCheckItem(router)
			router = mocks.MockAddIDLabelError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			dueDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			card, err := client.CreateShoppingList(&dueDate, []string{"Label"}, []string{})
			require.Error(t, err)
			require.Empty(t, card)
		})
	})

	t.Run("AddToShoppingList", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockListGetCardsWithCheckLists(router)
			router = mocks.MockGetChecklist(router)
			router = mocks.MockCreateCheckItem(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			url, err := client.AddToShoppingList([]string{""})
			require.NoError(t, err)
			require.NotEmpty(t, url)
		})

		t.Run("Success2", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockListGetCardsWithCheckLists(router)
			router = mocks.MockGetChecklistEmpty(router)
			router = mocks.MockCreateCheckItem(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			url, err := client.AddToShoppingList([]string{"Gyoza"})
			require.NoError(t, err)
			require.NotEmpty(t, url)
		})

		t.Run("ErrorGetShoppingList", func(t *testing.T) {
			// Same as GetShoppingList/ErrorGetMembers
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMemberError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			url, err := client.AddToShoppingList([]string{"Gyoza"})
			require.Error(t, err)
			require.Empty(t, url)
		})

		t.Run("ErrorIDCheckLists", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockListGetCards(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			url, err := client.AddToShoppingList([]string{"Gyoza"})
			require.Error(t, err)
			require.Empty(t, url)
		})

		t.Run("ErrorGetChecklist", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockListGetCardsWithCheckLists(router)
			router = mocks.MockGetChecklistError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			url, err := client.AddToShoppingList([]string{"Gyoza"})
			require.Error(t, err)
			require.Empty(t, url)
		})

		t.Run("ErrorCreateCheckItem", func(t *testing.T) {
			router := mux.NewRouter().StrictSlash(true)
			router = mocks.MockGetMember(router)
			router = mocks.MockMemberGetBoards(router)
			router = mocks.MockBoardGetLists(router)
			router = mocks.MockListGetCardsWithCheckLists(router)
			router = mocks.MockGetChecklistEmpty(router)
			router = mocks.MockCreateCheckItemError(router)
			server := httptest.NewServer(router)
			defer server.Close()
			client := mocks.NewTrelloClientWrapper(server, "apikey", "apitoken", "mid", "Board", "List", "Label")
			require.NotNil(t, client)

			url, err := client.AddToShoppingList([]string{"Gyoza"})
			require.Error(t, err)
			require.Empty(t, url)
		})
	})
}
