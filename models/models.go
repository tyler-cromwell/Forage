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
	logger := logrus.WithFields(logrus.Fields{
		"method": "GET",
		"filter": filter,
	})

	var doc bson.M
	err := mc.Collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		logger.WithError(err).Error("FindOne failed")
		return nil, err
	} else {
		logrus.WithFields(logrus.Fields{"document": doc}).Info("FindOne succeeded")
		return &doc, err
	}
}

func (mc *MongoClient) GetManyDocuments(ctx context.Context, logger *logrus.Entry, filter bson.D) ([]bson.M, error) {
	// Specify common fields
	logger = logrus.WithFields(logrus.Fields{"filter": filter})

	// Ask MongoDB to find the documents
	docs := make([]bson.M, 0)
	cursor, err := mc.Collection.Find(ctx, filter)
	if err != nil {
		logger.WithError(err).Error("GetManyDocuments failed")
		return nil, err
	}

	// Extract documents from MongoDB's response
	err = cursor.All(ctx, &docs)
	if err != nil {
		logger.WithError(err).Error("GetManyDocuments failed")
		return nil, err
	}

	// Cleanup
	err = cursor.Close(ctx)
	if err != nil {
		logger.WithError(err).Error("GetManyDocuments failed")
		return nil, err // return results anyway???
	}

	logger.WithFields(logrus.Fields{"quantity": len(docs)}).Info("GetManyDocuments succeeded")
	return docs, err
}

func (mc *MongoClient) PostOneDocument(ctx context.Context, doc interface{}) error {
	logger := logrus.WithFields(logrus.Fields{
		"method":   "POST",
		"document": doc,
	})

	_, err := mc.Collection.InsertOne(ctx, doc, nil)
	if err != nil {
		logger.WithError(err).Error("PostOneDocument failed")
	} else {
		logger.Info("PostOneDocument succeeded")
	}
	return err
}

/*
func (mc *Collection) postManyDocuments(ctx context.Context) error {
}

func (mc *Collection) putDocument(ctx) error {
}
*/

func (mc *MongoClient) DeleteOneDocument(ctx context.Context, filter bson.D) error {
	logger := logrus.WithFields(logrus.Fields{
		"method": "DELETE",
		"filter": filter,
	})

	_, err := mc.Collection.DeleteOne(ctx, filter)
	if err != nil {
		logger.WithError(err).Error("DeleteOneDocument failed")
	} else {
		logger.Info("DeleteOneDocument succeeded")
	}
	return err
}
