package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongo struct {
	OverrideFindOneDocument     func(context.Context, bson.D) (*bson.M, error)
	OverrideFindManyDocuments   func(context.Context, bson.M, *options.FindOptions) ([]bson.M, error)
	OverrideInsertOneDocument   func(context.Context, interface{}) error
	OverrideInsertManyDocuments func(context.Context, []interface{}) error
	OverrideUpdateOneDocument   func(context.Context, bson.D, interface{}) (int64, int64, error)
	OverrideDeleteOneDocument   func(context.Context, bson.D) error
	OverrideDeleteManyDocuments func(context.Context, bson.M) (int64, error)
}

func (mmc *MockMongo) FindOneDocument(ctx context.Context, filter bson.D) (*bson.M, error) {
	if mmc.OverrideFindOneDocument != nil {
		return mmc.OverrideFindOneDocument(ctx, filter)
	} else {
		var doc bson.M
		return &doc, nil
	}
}

func (mmc *MockMongo) FindDocuments(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	if mmc.OverrideFindManyDocuments != nil {
		return mmc.OverrideFindManyDocuments(ctx, filter, opts)
	} else {
		docs := make([]bson.M, 0)
		return docs, nil
	}
}

func (mmc *MockMongo) InsertOneDocument(ctx context.Context, doc interface{}) error {
	return nil
}

func (mmc *MockMongo) InsertManyDocuments(ctx context.Context, docs []interface{}) error {
	if mmc.OverrideInsertManyDocuments != nil {
		return mmc.OverrideInsertManyDocuments(ctx, docs)
	} else {
		return nil
	}
}

func (mmc *MockMongo) UpdateOneDocument(ctx context.Context, filter bson.D, update interface{}) (int64, int64, error) {
	if mmc.OverrideUpdateOneDocument != nil {
		return mmc.OverrideUpdateOneDocument(ctx, filter, update)
	} else {
		return 1, 1, nil
	}
}

func (mmc *MockMongo) DeleteOneDocument(ctx context.Context, filter bson.D) error {
	if mmc.OverrideDeleteOneDocument != nil {
		return mmc.OverrideDeleteOneDocument(ctx, filter)
	} else {
		return nil
	}
}

func (mmc *MockMongo) DeleteManyDocuments(ctx context.Context, filter bson.M) (int64, error) {
	return 1, nil
}
