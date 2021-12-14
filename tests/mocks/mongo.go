package mocks

import (
	"context"

	"github.com/tyler-cromwell/forage/clients"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func NewMongoClientWrapper(mt *mtest.T, ctx context.Context, mongoUri string) (*clients.Mongo, error) {
	return &clients.Mongo{
		Uri:        mongoUri,
		Client:     mt.Client,
		Database:   mt.DB,
		Collection: mt.Coll,
	}, nil
}
