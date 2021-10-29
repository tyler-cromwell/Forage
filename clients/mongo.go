package clients

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	Collection *mongo.Collection
}

func (mc *Mongo) FindOneDocument(ctx context.Context, filter bson.D) (*bson.M, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "mongo.FindOneDocument",
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
		log.Debug("Success")
		return &doc, nil
	}
}

func (mc *Mongo) FindDocuments(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "mongo.FindDocuments",
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
		log.Debug("Success")
		return docs, nil
	}
}

func (mc *Mongo) InsertOneDocument(ctx context.Context, doc interface{}) error {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":       "mongo.InsertOneDocument",
		"document": doc,
	})

	// Ask MongoDB to insert the document
	_, err := mc.Collection.InsertOne(ctx, doc, nil)
	if err != nil {
		log.WithError(err).Error("Failed to insert document")
		return err
	} else {
		log.Debug("Success")
		return nil
	}
}

func (mc *Mongo) InsertManyDocuments(ctx context.Context, docs []interface{}) error {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":       "mongo.InsertManyDocuments",
		"quantity": len(docs),
	})

	// Ask MongoDB to insert the documents
	_, err := mc.Collection.InsertMany(ctx, docs, nil)
	if err != nil {
		log.WithError(err).Error("Failed to insert documents")
		return err
	} else {
		log.Debug("Success")
		return nil
	}
}

func (mc *Mongo) UpdateOneDocument(ctx context.Context, filter bson.D, update interface{}) (int64, int64, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "mongo.UpdateOneDocument",
		"filter": filter,
	})

	// Ask MongoDB to update the document
	result, err := mc.Collection.UpdateOne(ctx, filter, update, nil)
	if result != nil && result.MatchedCount == 0 {
		// Update completed but no document was found
		customError := fmt.Errorf("no document matching filter")
		log.WithError(customError).Warn("Failed to update document")
		return result.MatchedCount, result.ModifiedCount, customError
	} else if err != nil {
		// Update failed
		log.WithError(err).Error("Failed to update document")
		return result.MatchedCount, result.ModifiedCount, err
	} else {
		log.Debug("Success")
		return result.MatchedCount, result.ModifiedCount, nil
	}
}

func (mc *Mongo) DeleteOneDocument(ctx context.Context, filter bson.D) error {
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
		log.Debug("Success")
		return nil
	}
}

func (mc *Mongo) DeleteManyDocuments(ctx context.Context, filter bson.M) (int64, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "mongo.DeleteManyDocuments",
		"filter": filter,
	})

	// Ask MongoDB to delete the documents
	result, err := mc.Collection.DeleteMany(ctx, filter)
	if err != nil {
		log.WithError(err).Error("Failed to delete documents")
		return result.DeletedCount, err
	} else {
		log.WithFields(logrus.Fields{"quantity": result.DeletedCount}).Debug("Success")
		return result.DeletedCount, nil
	}
}
