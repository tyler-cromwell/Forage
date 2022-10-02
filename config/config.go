package config

import (
	"context"
	"time"

	"github.com/adlio/trello"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MongoCollectionIngredients = "ingredients"
const MongoCollectionRecipes = "recipes"

type MongoHandle interface {
	Collections(context.Context) ([]string, error)
	FindOneDocument(context.Context, string, bson.D) (*bson.M, error)
	FindDocuments(context.Context, string, bson.M, *options.FindOptions) ([]bson.M, error)
	//InsertOneDocument(context.Context, string, interface{}) error
	InsertManyDocuments(context.Context, string, []interface{}) error
	UpdateOneDocument(context.Context, string, bson.D, interface{}) (int64, int64, error)
	DeleteOneDocument(context.Context, string, bson.D) error
	DeleteManyDocuments(context.Context, string, bson.M) (int64, error)
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
	Lookahead      time.Duration
	Silence        bool
	Time           string
	Timezone       string
	LogrusLevel    logrus.Level
	ListenSocket   string
	Scheduler      *gocron.Scheduler
	Mongo          MongoHandle  //*clients.Mongo
	Trello         TrelloHandle //*clients.Trello
	Twilio         TwilioHandle //*clients.Twilio
}
