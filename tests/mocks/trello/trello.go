package trello

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/adlio/trello"
	"github.com/gorilla/mux"
	"github.com/tyler-cromwell/forage/clients"
)

type MockTrello struct {
	Key                        string
	Token                      string
	MemberID                   string
	BoardName                  string
	ListName                   string
	LabelsStr                  string
	OverrideGetShoppingList    func() (*trello.Card, error)
	OverrideCreateShoppingList func(*time.Time, []string, []string) (string, error)
	OverrideAddToShoppingList  func([]string) (string, error)
}

func (mmc *MockTrello) GetShoppingList() (*trello.Card, error) {
	if mmc.OverrideGetShoppingList != nil {
		return mmc.OverrideGetShoppingList()
	} else {
		var card trello.Card
		return &card, nil
	}
}

func (mmc *MockTrello) CreateShoppingList(dueDate *time.Time, applyLabels []string, listItems []string) (string, error) {
	if mmc.OverrideCreateShoppingList != nil {
		return mmc.OverrideCreateShoppingList(dueDate, applyLabels, listItems)
	} else {
		return "", nil
	}
}

func (mmc *MockTrello) AddToShoppingList(itemNames []string) (string, error) {
	if mmc.OverrideAddToShoppingList != nil {
		return mmc.OverrideAddToShoppingList(itemNames)
	} else {
		return "", nil
	}
}

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

func MockGetMember(router *mux.Router) *mux.Router {
	router.HandleFunc("/members/{mid}", func(response http.ResponseWriter, request *http.Request) {
		member := trello.Member{
			ID: "mid",
		}
		b, _ := json.Marshal(member)
		response.WriteHeader(http.StatusOK)
		response.Write(b)
	})
	return router
}

func MockGetMemberError(router *mux.Router) *mux.Router {
	router.HandleFunc("/members/{mid}", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusInternalServerError)
	})
	return router
}

func MockMemberGetBoards(router *mux.Router) *mux.Router {
	router.HandleFunc("/members/{mid}/boards", func(response http.ResponseWriter, request *http.Request) {
		boards := make([]trello.Board, 1)
		boards[0] = trello.Board{
			ID:   "board",
			Name: "Board",
		}
		b, _ := json.Marshal(boards)
		response.WriteHeader(http.StatusOK)
		response.Write(b)
	})
	return router
}

func MockMemberGetBoardsError(router *mux.Router) *mux.Router {
	router.HandleFunc("/members/{mid}/boards", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusInternalServerError)
	})
	return router
}

func MockBoardGetLabels(router *mux.Router) *mux.Router {
	router.HandleFunc("/boards/{bid}/labels", func(response http.ResponseWriter, request *http.Request) {
		labels := make([]trello.Label, 1)
		labels[0] = trello.Label{
			ID:   "label",
			Name: "Label",
		}
		l, _ := json.Marshal(labels)
		response.WriteHeader(http.StatusOK)
		response.Write(l)
	})
	return router
}

func MockBoardGetLabelsError(router *mux.Router) *mux.Router {
	router.HandleFunc("/boards/{bid}/labels", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusInternalServerError)
	})
	return router
}

func MockBoardGetLists(router *mux.Router) *mux.Router {
	router.HandleFunc("/boards/{bid}/lists", func(response http.ResponseWriter, request *http.Request) {
		lists := make([]trello.List, 1)
		lists[0] = trello.List{
			ID:   "list",
			Name: "List",
		}
		l, _ := json.Marshal(lists)
		response.WriteHeader(http.StatusOK)
		response.Write(l)
	})
	return router
}

func MockBoardGetListsError(router *mux.Router) *mux.Router {
	router.HandleFunc("/boards/{bid}/lists", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusInternalServerError)
	})
	return router
}

func MockListGetCards(router *mux.Router) *mux.Router {
	router.HandleFunc("/lists/{lid}/cards", func(response http.ResponseWriter, request *http.Request) {
		cards := make([]trello.Card, 1)
		cards[0] = trello.Card{
			ID:   "shopping_list",
			Name: "Shopping List",
		}
		c, _ := json.Marshal(cards)
		response.WriteHeader(http.StatusOK)
		response.Write(c)
	})
	return router
}

func MockListGetCardsWithCheckLists(router *mux.Router) *mux.Router {
	router.HandleFunc("/lists/{lid}/cards", func(response http.ResponseWriter, request *http.Request) {
		cards := make([]trello.Card, 1)
		cards[0] = trello.Card{
			ID:           "shopping_list",
			Name:         "Shopping List",
			URL:          "www.mock.url.com",
			IDCheckLists: []string{"groceries"},
			Checklists: []*trello.Checklist{
				{
					ID:   "groceries",
					Name: "Groceries",
				},
			},
		}
		c, _ := json.Marshal(cards)
		response.WriteHeader(http.StatusOK)
		response.Write(c)
	})
	return router
}

func MockListAddCards(router *mux.Router) *mux.Router {
	router.HandleFunc("/lists/{lid}/cards", func(response http.ResponseWriter, request *http.Request) {
		card := trello.Card{
			ID:   "shopping_list",
			Name: "Shopping List",
		}
		c, _ := json.Marshal(card)
		response.WriteHeader(http.StatusOK)
		response.Write(c)
	})
	return router
}

func MockListParamCardsError(router *mux.Router) *mux.Router {
	router.HandleFunc("/lists/{lid}/cards", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusInternalServerError)
	})
	return router
}

func MockCardSetPos(router *mux.Router) *mux.Router {
	router.HandleFunc("/cards/{cid}", func(response http.ResponseWriter, request *http.Request) {
		card := trello.Card{
			ID:   "shopping_list",
			Name: "Shopping List",
		}
		c, _ := json.Marshal(card)
		response.WriteHeader(http.StatusOK)
		response.Write(c)
	})
	return router
}

func MockCardsParamError(router *mux.Router) *mux.Router {
	router.HandleFunc("/cards/{cid}", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusInternalServerError)
	})
	return router
}

func MockCreateChecklist(router *mux.Router) *mux.Router {
	router.HandleFunc("/cards/{cid}/checklists", func(response http.ResponseWriter, request *http.Request) {
		card := trello.Card{
			ID:   "shopping_list",
			Name: "Shopping List",
			Checklists: []*trello.Checklist{
				{
					ID:   "groceries",
					Name: "Groceries",
				},
			},
		}
		c, _ := json.Marshal(card)
		response.WriteHeader(http.StatusOK)
		response.Write(c)
	})
	return router
}

func MockCreateChecklistError(router *mux.Router) *mux.Router {
	router.HandleFunc("/cards/{cid}/checklists", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusInternalServerError)
	})
	return router
}

func MockGetChecklist(router *mux.Router) *mux.Router {
	router.HandleFunc("/checklists/{cid}", func(response http.ResponseWriter, request *http.Request) {
		checklist := trello.Checklist{
			ID:   "groceries",
			Name: "Groceries",
			CheckItems: []trello.CheckItem{
				{
					ID:   "gyoza",
					Name: "Gyoza",
				},
			},
		}
		c, _ := json.Marshal(checklist)
		response.WriteHeader(http.StatusOK)
		response.Write(c)
	})
	return router
}

func MockGetChecklistEmpty(router *mux.Router) *mux.Router {
	router.HandleFunc("/checklists/{cid}", func(response http.ResponseWriter, request *http.Request) {
		checklist := trello.Checklist{
			ID:         "groceries",
			Name:       "Groceries",
			CheckItems: []trello.CheckItem{},
		}
		c, _ := json.Marshal(checklist)
		response.WriteHeader(http.StatusOK)
		response.Write(c)
	})
	return router
}

func MockGetChecklistError(router *mux.Router) *mux.Router {
	router.HandleFunc("/checklists/{cid}", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusInternalServerError)
	})
	return router
}

func MockCreateCheckItem(router *mux.Router) *mux.Router {
	router.HandleFunc("/checklists/{cid}/checkItems", func(response http.ResponseWriter, request *http.Request) {
		checkItem := trello.CheckItem{}
		c, _ := json.Marshal(checkItem)
		response.WriteHeader(http.StatusOK)
		response.Write(c)
	})
	return router
}

func MockCreateCheckItemError(router *mux.Router) *mux.Router {
	router.HandleFunc("/checklists/{cid}/checkItems", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusInternalServerError)
	})
	return router
}

func MockAddIDLabelError(router *mux.Router) *mux.Router {
	router.HandleFunc("/cards/{cid}/idLabels", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusInternalServerError)
	})
	return router
}
