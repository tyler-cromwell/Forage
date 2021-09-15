package models

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient struct {
	Collection *mongo.Collection
}

func (mc *MongoClient) GetOneDocument(ctx context.Context, filter bson.D) (*bson.M, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{"filter": filter})

	// Ask MongoDB to find the document
	var doc bson.M
	err := mc.Collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		log.WithError(err).Error("Decode failed")
		return nil, err
	} else {
		log.WithFields(logrus.Fields{"document": doc}).Info("Decode succeeded")
		return &doc, err
	}
}

func (mc *MongoClient) GetManyDocuments(ctx context.Context, filter bson.D) ([]bson.M, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "mongo.GetManyDocuments",
		"filter": filter,
	})

	// Ask MongoDB to find the documents
	docs := make([]bson.M, 0)
	cursor, err := mc.Collection.Find(ctx, filter)
	if err != nil {
		log.WithError(err).Error("Failure: Find")
		return nil, err
	}

	// Extract documents from MongoDB's response
	err = cursor.All(ctx, &docs)
	if err != nil {
		log.WithError(err).Error("Failure: All")
		return nil, err
	}

	// Cleanup
	err = cursor.Close(ctx)
	if err != nil {
		log.WithError(err).Error("Failure: Close")
		return nil, err // maybe return results anyway???
	} else {
		return docs, nil
	}
}

func (mc *MongoClient) PostOneDocument(ctx context.Context, doc interface{}) error {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{"document": doc})

	// Ask MongoDB to insert the document
	_, err := mc.Collection.InsertOne(ctx, doc, nil)
	if err != nil {
		log.WithError(err).Error("InsertOne failed")
		return err
	} else {
		log.Info("InsertOne succeeded")
		return nil
	}
}

/*
func (mc *Collection) postManyDocuments(ctx context.Context) error {
}

func (mc *Collection) putDocument(ctx) error {
}
*/

func (mc *MongoClient) DeleteOneDocument(ctx context.Context, filter bson.D) error {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{"filter": filter})

	// Ask MongoDB to insert the document
	_, err := mc.Collection.DeleteOne(ctx, filter)
	if err != nil {
		log.WithError(err).Error("DeleteOne failed")
		return err
	} else {
		log.Info("DeleteOne succeeded")
		return nil
	}
}
