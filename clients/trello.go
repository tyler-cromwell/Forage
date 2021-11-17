package clients

import (
	"fmt"
	"strings"
	"time"

	"github.com/adlio/trello"
)

type Trello struct {
	Key       string
	Token     string
	MemberID  string
	BoardName string
	ListName  string
	Client    *trello.Client
}

func (tc *Trello) GetShoppingList() (*trello.Card, error) {
	var board *trello.Board
	var list *trello.List
	var card *trello.Card

	// Get Trello member
	member, err := tc.Client.GetMember(tc.MemberID, trello.Defaults())
	if err != nil {
		return nil, err
	}

	// Get all boards
	boards, err := member.GetBoards(trello.Defaults())
	if err != nil {
		return nil, err
	}
	for _, b := range boards {
		if b.Name == tc.BoardName {
			board = b
			break
		}
	}

	// Get List with given name
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		return nil, err
	}
	for _, l := range lists {
		if l.Name == tc.ListName {
			list = l
			break
		}
	}

	// Get Card with name "Shopping List"
	cards, err := list.GetCards()
	if err != nil {
		return nil, err
	}
	for _, c := range cards {
		if c.Name == "Shopping List" {
			card = c
			break
		}
	}

	return card, nil
}

func (tc *Trello) CreateShoppingList(dueDate *time.Time, applyLabels []string, listItems []string) (string, error) {
	var board *trello.Board
	var list *trello.List
	var labelIDs []string

	// Get Trello member
	member, err := tc.Client.GetMember(tc.MemberID, trello.Defaults())
	if err != nil {
		return "", err
	}

	// Get all boards
	boards, err := member.GetBoards(trello.Defaults())
	if err != nil {
		return "", err
	}
	for _, b := range boards {
		if b.Name == tc.BoardName {
			board = b
			break
		}
	}

	// Get List with given name
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		return "", err
	}
	for _, l := range lists {
		if l.Name == tc.ListName {
			list = l
			break
		}
	}

	// Get labels
	labels, err := board.GetLabels(trello.Defaults())
	if err != nil {
		return "", err
	}
	for _, l := range labels {
		for _, al := range applyLabels {
			if l.Name == al {
				labelIDs = append(labelIDs, l.ID)
			}
		}
	}

	// Construct card
	cardName := "Shopping List"
	card := &trello.Card{
		Name: cardName,
		Desc: "A list of items that must be bought in the near future.",
		Due:  dueDate,
	}

	// Add shopping list card
	err = list.AddCard(card, trello.Defaults())
	if err != nil {
		return "", err
	}

	// Set the card's position in the list
	err = card.SetPos(1.0)
	if err != nil {
		return "", err
	}

	// Add checklist to card
	checklist, err := tc.Client.CreateChecklist(card, "Groceries", trello.Defaults())
	if err != nil {
		return "", err
	} else {
		// Add items to the checklist
		for _, item := range listItems {
			_, err := checklist.CreateCheckItem(item)
			if err != nil {
				return "", err
			}
		}
	}

	// Add labels to the card
	for _, labelID := range labelIDs {
		err = card.AddIDLabel(labelID)
		if err != nil {
			return "", err
		}
	}

	return card.URL, nil
}

func (tc *Trello) AddToShoppingList(itemNames []string) (string, error) {
	card, err := tc.GetShoppingList()
	if err != nil {
		return "", err
	}

	checklistIDs := card.IDCheckLists
	if len(checklistIDs) == 0 {
		return "", fmt.Errorf("no checklists attached to card")
	}

	checklist, err := tc.Client.GetChecklist(checklistIDs[0], trello.Defaults())
	if err != nil {
		return "", err
	} else {
		existingItems := checklist.CheckItems

		// Add items to the checklist
		for _, name := range itemNames {
			exists := false
			for _, existingItem := range existingItems {
				if strings.HasPrefix(existingItem.Name, name) {
					exists = true
					break
				}
			}

			if !exists {
				_, err := checklist.CreateCheckItem(name)
				if err != nil {
					return "", err
				}
			}
		}
	}

	return card.URL, nil
}
