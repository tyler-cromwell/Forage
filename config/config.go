package config

import (
	"context"
	"time"

	"github.com/adlio/trello"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoHandle interface {
	FindOneDocument(context.Context, bson.D) (*bson.M, error)
	FindDocuments(context.Context, bson.M, *options.FindOptions) ([]bson.M, error)
	InsertOneDocument(context.Context, interface{}) error
	InsertManyDocuments(context.Context, []interface{}) error
	UpdateOneDocument(context.Context, bson.D, interface{}) (int64, int64, error)
	DeleteOneDocument(context.Context, bson.D) error
	DeleteManyDocuments(context.Context, bson.M) (int64, error)
}

type TrelloHandle interface {
	GetShoppingList() (*trello.Card, error)
	CreateShoppingList(*time.Time, []string, []string) (string, error)
	AddToShoppingList([]string) (string, error)
}

type TwilioHandle interface {
	ComposeMessage(int, int, string) string
	SendMessage(string, string, string) (string, error)
}

type Configuration struct {
	ContextTimeout time.Duration
	Interval       time.Duration
	Lookahead      time.Duration
	LogrusLevel    logrus.Level
	ListenSocket   string
	Mongo          MongoHandle  //*clients.Mongo
	Trello         TrelloHandle //*clients.Trello
	Twilio         TwilioHandle //*clients.Twilio
}
