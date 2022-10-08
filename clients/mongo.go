package clients

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	Uri      string
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoClientWrapper(ctx context.Context, mongoUri string) (*Mongo, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":  "mongo.NewMongoClientWrapper",
		"uri": mongoUri,
	})

	// Log diagnostic information
	log.Trace("Begin function")
	defer log.Trace("End function")

	// Initialize MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.WithError(err).Error("Failed to initialize MongoDB client")
		return nil, err
	}

	// Connect to database instance
	_ = client.Connect(ctx)
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.WithError(err).Error("Failed to connect to MongoDB instance")
		return nil, err
	} else {
		log.Info("Connected to MongoDB")
	}

	// Specify database & collection
	database := client.Database("forage")
	return &Mongo{
		Uri:      mongoUri,
		Client:   client,
		Database: database,
	}, nil
}

func (mc *Mongo) Collections(ctx context.Context) ([]string, error) {
	return mc.Database.ListCollectionNames(ctx, bson.D{})
}

func (mc *Mongo) FindOneDocument(ctx context.Context, collection string, filter bson.D) (*bson.M, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":         "mongo.FindOneDocument",
		"collection": collection,
		"database":   mc.Uri,
		"filter":     filter,
	})

	// Log diagnostic information
	log.Trace("Begin function")
	defer log.Trace("End function")

	// Ask MongoDB to find the document
	var doc bson.M
	result := mc.Database.Collection(collection).FindOne(ctx, filter)
	log.WithFields(logrus.Fields{"value": result}).Debug("Result")
	err := result.Decode(&doc)
	if err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
		// Search completed but no document was found
		log.WithError(err).Warn("Failed to find document")
		return nil, err
	} else if err != nil {
		// Actual failure
		log.WithError(err).Error("Failed to find document")
		return nil, err
	}

	log.Info("Succeeded")
	return &doc, nil
}

func (mc *Mongo) FindDocuments(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":         "mongo.FindDocuments",
		"collection": collection,
		"database":   mc.Uri,
		"filter":     filter,
		"opts":       opts,
	})

	// Log diagnostic information
	log.Trace("Begin function")
	defer log.Trace("End function")

	// Ask MongoDB to find the documents
	docs := make([]bson.M, 0)
	cursor, err := mc.Database.Collection(collection).Find(ctx, filter, opts)
	if err != nil {
		log.WithError(err).Error("Failed to find documents")
		return nil, err
	}

	// Extract documents from MongoDB's response
	err = cursor.All(ctx, &docs)
	if err != nil {
		log.WithError(err).Error("Failed to decode documents")
		return nil, err
	} else {
		log.WithFields(logrus.Fields{"quantity": len(docs)}).Info("Succeeded")
		return docs, nil
	}
}

/*
func (mc *Mongo) InsertOneDocument(ctx context.Context, collection string, doc interface{}) error {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":       "mongo.InsertOneDocument",
		"collection": collection,
		"database": mc.Uri,
		"document": doc,
	})

	// Log diagnostic information
	log.Trace("Begin function")
	defer log.Trace("End function")

	// Ask MongoDB to insert the document
	_, err := mc.Database.Collection(collection).InsertOne(ctx, doc, nil)
	if err != nil {
		log.WithError(err).Error("Failed to insert document")
		return err
	} else {
		log.Info("Succeeded")
		return nil
	}
}
*/

func (mc *Mongo) InsertManyDocuments(ctx context.Context, collection string, docs []interface{}) error {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":         "mongo.InsertManyDocuments",
		"collection": collection,
		"database":   mc.Uri,
		"quantity":   len(docs),
	})

	// Log diagnostic information
	log.Trace("Begin function")
	defer log.Trace("End function")

	// Ask MongoDB to insert the documents
	_, err := mc.Database.Collection(collection).InsertMany(ctx, docs, nil)
	if err != nil {
		log.WithError(err).Error("Failed to insert documents")
		return err
	} else {
		log.Info("Succeeded")
		return nil
	}
}

func (mc *Mongo) UpdateOneDocument(ctx context.Context, collection string, filter bson.D, update interface{}) (int64, int64, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":         "mongo.UpdateOneDocument",
		"collection": collection,
		"database":   mc.Uri,
		"filter":     filter,
	})

	// Log diagnostic information
	log.Trace("Begin function")
	defer log.Trace("End function")

	// Ask MongoDB to update the document
	result, err := mc.Database.Collection(collection).UpdateOne(ctx, filter, update, nil)
	if err != nil {
		// Update failed
		log.WithError(err).Error("Failed to update document")
		return 0, 0, err
	} else {
		log.WithFields(logrus.Fields{"matched": result.MatchedCount, "modified": result.ModifiedCount}).Info("Succeeded")
		return result.MatchedCount, result.ModifiedCount, nil
	}
}

func (mc *Mongo) DeleteOneDocument(ctx context.Context, collection string, filter bson.D) error {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":         "mongo.DeleteOneDocument",
		"collection": collection,
		"database":   mc.Uri,
		"filter":     filter,
	})

	// Log diagnostic information
	log.Trace("Begin function")
	defer log.Trace("End function")

	// Ask MongoDB to delete the document
	_, err := mc.Database.Collection(collection).DeleteOne(ctx, filter)
	if err != nil {
		log.WithError(err).Error("Failed to delete document")
		return err
	} else {
		log.Info("Succeeded")
		return nil
	}
}

func (mc *Mongo) DeleteManyDocuments(ctx context.Context, collection string, filter bson.M) (int64, error) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":         "mongo.DeleteManyDocuments",
		"collection": collection,
		"database":   mc.Uri,
		"filter":     filter,
	})

	// Log diagnostic information
	log.Trace("Begin function")
	defer log.Trace("End function")

	// Ask MongoDB to delete the documents
	result, err := mc.Database.Collection(collection).DeleteMany(ctx, filter)
	if err != nil {
		log.WithError(err).Error("Failed to delete documents")
		return 0, err
	} else {
		log.WithFields(logrus.Fields{"quantity": result.DeletedCount}).Info("Succeeded")
		return result.DeletedCount, nil
	}
}
