package tests

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/tyler-cromwell/forage/clients"
	"github.com/tyler-cromwell/forage/tests/mocks"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoClient(t *testing.T) {
	// Discard logging output
	logrus.SetOutput(ioutil.Discard)

	// Setup context
	//	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Mock the Mongo database
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("NewMongoClientWrapper", func(mt *mtest.T) {
		// Case: "Success"
		client, err := clients.NewMongoClientWrapper(ctx, "mongodb://127.0.0.1:27017") // Expects an actual instance running
		require.NoError(mt, err)
		require.NotNil(mt, client)

		// Case: "Failed to initialize MongoDB client"
		client, err = clients.NewMongoClientWrapper(ctx, "mongodb://0.0.0.0:0")
		require.Error(mt, err)
		require.Nil(mt, client)

		// Case: "Failed to connect to MongoDB instance"
		client, err = clients.NewMongoClientWrapper(ctx, "mongodb://0.0.0.0:1")
		require.Error(mt, err)
		require.Nil(mt, client)
	})

	mt.Run("FindOneDocument", func(mt *mtest.T) {
		client, err := mocks.NewMongoClientWrapper(mt, ctx, "")
		require.NoError(mt, err)
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		document1 := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", id1},
			{"name", "john"},
			{"email", "john.doe@test.com"},
		})

		// Case: "Success"
		mt.ClearMockResponses()
		mt.AddMockResponses(document1)
		filter := bson.D{{"_id", id1}}
		_, err = client.FindOneDocument(ctx, filter)
		require.NoError(mt, err)

		// Case: "Failed to find document" (Search completed but document not found)
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: "mongo: no documents in result"}))
		filter = bson.D{{"_id", id2}}
		_, err = client.FindOneDocument(ctx, filter)
		require.Error(mt, err)

		// Case: "Failed to find document" (Actual failure)
		mt.ClearMockResponses()
		filter = bson.D{{"_id", id1}}
		_, err = client.FindOneDocument(ctx, filter)
		require.Error(mt, err)

		mt.ClearMockResponses()
	})

	mt.Run("FindDocuments", func(mt *mtest.T) {
		client, err := mocks.NewMongoClientWrapper(mt, ctx, "")
		require.NoError(mt, err)
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		document1 := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", id1},
			{"name", "john"},
			{"email", "john.doe@test.com"},
		})
		document2 := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", id2},
			{"name", "john"},
			{"email", "foo.bar@test.com"},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)

		// Case: "Success"
		mt.ClearMockResponses()
		mt.AddMockResponses(document1, document2, killCursors)
		filter := bson.M{"_id": id1}
		_, err = client.FindDocuments(ctx, filter, nil)
		require.NoError(mt, err)

		// Case: "Failed to find documents"
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: "mongo: no documents in result"}))
		filter = bson.M{"_id": id1}
		_, err = client.FindDocuments(ctx, filter, nil)
		require.Error(mt, err)

		// Case: "Failed to decode documents"
		mt.ClearMockResponses()
		mt.AddMockResponses(document1)
		filter = bson.M{"_id": id1}
		_, err = client.FindDocuments(ctx, filter, nil)
		require.Error(mt, err)

		mt.ClearMockResponses()
	})

	mt.Run("InsertOneDocument", func(mt *mtest.T) {
		client, err := mocks.NewMongoClientWrapper(mt, ctx, "")
		require.NoError(mt, err)

		doc1 := bson.D{
			{"name", "john"},
			{"email", "john.doe@test.com"},
		}

		// Case: "Success"
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err = client.InsertOneDocument(ctx, doc1)
		require.NoError(mt, err)

		// Case: "Failed to insert document"
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: "command failure"}))
		err = client.InsertOneDocument(ctx, doc1)
		require.Error(mt, err)

		mt.ClearMockResponses()
	})

	mt.Run("InsertManyDocuments", func(mt *mtest.T) {
		client, err := mocks.NewMongoClientWrapper(mt, ctx, "")
		require.NoError(mt, err)

		doc1 := bson.D{
			{"name", "john"},
			{"email", "john.doe@test.com"},
		}

		// Case: "Success"
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err = client.InsertManyDocuments(ctx, []interface{}{doc1})
		require.NoError(mt, err)

		// Case: "Failed to insert document"
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: "command failure"}))
		err = client.InsertManyDocuments(ctx, []interface{}{doc1})
		require.Error(mt, err)

		mt.ClearMockResponses()
	})

	mt.Run("UpdateOneDocument", func(mt *mtest.T) {
		client, err := mocks.NewMongoClientWrapper(mt, ctx, "")
		require.NoError(mt, err)
		id1 := primitive.NewObjectID()

		doc1 := bson.D{
			{"_id", id1},
			{"name", "john"},
			{"email", "john.doe@test.com"},
		}

		// Case: "Success"
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateSuccessResponse(), bson.D{
			{"ok", 1},
			{"value", doc1},
		})
		filter := bson.D{{"_id", id1}}
		_, _, err = client.UpdateOneDocument(ctx, filter, bson.M{"$set": doc1})
		require.NoError(mt, err)

		// Case: "Failed to update document"
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: "command failure"}))
		filter = bson.D{{"_id", id1}}
		_, _, err = client.UpdateOneDocument(ctx, filter, doc1)
		require.Error(mt, err)

		mt.ClearMockResponses()
	})

	mt.Run("DeleteOneDocument", func(mt *mtest.T) {
		client, err := mocks.NewMongoClientWrapper(mt, ctx, "")
		require.NoError(mt, err)
		id1 := primitive.NewObjectID()

		// Case: "Success"
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		filter := bson.D{{"_id", id1}}
		err = client.DeleteOneDocument(ctx, filter)
		require.NoError(mt, err)

		// Case: "Failed to delete document"
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: "command failure"}))
		filter = bson.D{{"_id", id1}}
		err = client.DeleteOneDocument(ctx, filter)
		require.Error(mt, err)

		mt.ClearMockResponses()
	})

	mt.Run("DeleteManyDocuments", func(mt *mtest.T) {
		client, err := mocks.NewMongoClientWrapper(mt, ctx, "")
		require.NoError(mt, err)
		id1 := primitive.NewObjectID()

		// Case: "Success"
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		filter := bson.M{"_id": bson.M{"$in": []primitive.ObjectID{id1}}}
		_, err = client.DeleteManyDocuments(ctx, filter)
		require.NoError(mt, err)

		// Case: "Failed to delete documents"
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: "command failure"}))
		filter = bson.M{"_id": bson.M{"$in": []primitive.ObjectID{id1}}}
		_, err = client.DeleteManyDocuments(ctx, filter)
		require.Error(mt, err)

		mt.ClearMockResponses()
	})
}
