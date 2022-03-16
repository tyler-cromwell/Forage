package trello

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/adlio/trello"
	"github.com/gorilla/mux"
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
