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
	log := logrus.WithFields(logrus.Fields{
		"at":     "mongo.GetOneDocument",
		"filter": filter,
	})

	// Ask MongoDB to find the document
	var doc bson.M
	err := mc.Collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		log.WithError(err).Error("Failed to decode document")
		return nil, err
	} else {
		return &doc, nil
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
		log.WithError(err).Error("Failed to find documents")
		return nil, err
	}

	// Extract documents from MongoDB's response
	err = cursor.All(ctx, &docs)
	if err != nil {
		log.WithError(err).Error("Failed to decode documents")
		return nil, err
	}

	// Cleanup
	err = cursor.Close(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to close cursor")
		return nil, err // maybe return results anyway???
	} else {
		return docs, nil
	}
}

func (mc *MongoClient) PostOneDocument(ctx context.Context, doc interface{}) error {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":       "mongo.PostOneDocument",
		"document": doc,
	})

	// Ask MongoDB to insert the document
	_, err := mc.Collection.InsertOne(ctx, doc, nil)
	if err != nil {
		log.WithError(err).Error("Failed to insert document")
		return err
	} else {
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
	log := logrus.WithFields(logrus.Fields{
		"at":     "mongo.DeleteOneDocument",
		"filter": filter,
	})

	// Ask MongoDB to delete the document
	_, err := mc.Collection.DeleteOne(ctx, filter)
	if err != nil {
		log.WithError(err).Error("Failed to delete document")
		return err
	} else {
		return nil
	}
}
