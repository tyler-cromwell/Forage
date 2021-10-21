package models

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	result := mc.Collection.FindOne(ctx, filter)
	err := result.Err()
	if err != nil && err.Error() == "mongo: no documents in result" {
		// Search completed but no document was found
		log.WithError(err).Warn("Failed to find document")
		return nil, err
	} else if err != nil {
		// Search failed
		log.WithError(err).Error("Failed to find document")
		return nil, err
	}

	// MongoDB found the document
	err = result.Decode(&doc)
	if err != nil {
		log.WithError(err).Error("Failed to decode document")
		return nil, err
	} else {
		return &doc, nil
	}
}

func (mc *MongoClient) GetManyDocuments(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "mongo.GetManyDocuments",
		"filter": filter,
	})

	// Ask MongoDB to find the documents
	docs := make([]bson.M, 0)
	cursor, err := mc.Collection.Find(ctx, filter, opts)
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
func (mc *MongoClient) PostManyDocuments(ctx context.Context, documents []bson.M) error {
}
*/

func (mc *MongoClient) PutOneDocument(ctx context.Context, filter bson.D, update interface{}) (int64, int64, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "mongo.PutOneDocument",
		"filter": filter,
	})

	// Ask MongoDB to update the document
	result, err := mc.Collection.UpdateOne(ctx, filter, update, nil)
	if result.MatchedCount == 0 {
		// Update completed but no document was found
		customError := fmt.Errorf("no document matching filter")
		log.WithError(customError).Warn("Failed to update document")
		return result.MatchedCount, result.ModifiedCount, customError
	} else if err != nil {
		// Update failed
		log.WithError(err).Error("Failed to update document")
		return result.MatchedCount, result.ModifiedCount, err
	} else {
		return result.MatchedCount, result.ModifiedCount, nil
	}
}

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

/*
func (mc *MongoClient) DeleteManyDocuments(ctx context.Context, filter bson.D) error {

}
*/
