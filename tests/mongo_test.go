package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tyler-cromwell/forage/tests/mocks"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoClient(t *testing.T) {
	// Discard logging output
	//logrus.SetOutput(ioutil.Discard)

	// Mock the Mongo database
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("FindOneDocument", func(mt *mtest.T) {
		//cases := []struct{}{}
		client, err := mocks.NewMongoClientWrapper(mt, context.Background(), "")
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		document1 := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", id1},
			{"name", "john"},
			{"email", "john.doe@test.com"},
		})

		// Case: "Success" (happy case)
		mt.ClearMockResponses()
		mt.AddMockResponses(document1)
		filter := bson.D{{"_id", id1}}
		_, err = client.FindOneDocument(context.Background(), filter)
		require.NoError(mt, err)

		// Case: "Failed to find document" (Search completed but document not found)
		mt.ClearMockResponses()
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Message: "mongo: no documents in result",
		}))
		filter = bson.D{{"_id", id2}}
		_, err = client.FindOneDocument(context.Background(), filter)
		require.Error(mt, err)

		// Case: "Failed to find document" (Actual failure)
		mt.ClearMockResponses()
		filter = bson.D{{"_id", id1}}
		_, err = client.FindOneDocument(context.Background(), filter)
		require.Error(mt, err)

		mt.ClearMockResponses()
	})

	mt.Run("FindDocuments", func(mt *mtest.T) {
		//cases := []struct{}{}
		client, err := mocks.NewMongoClientWrapper(mt, context.Background(), "")
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

		// Documents found & decoded successfully (happy case)
		mt.ClearMockResponses()
		mt.AddMockResponses(document1, document2, killCursors)
		filter := bson.M{"_id": id1}
		_, err = client.FindDocuments(context.Background(), filter, nil)
		t.Log(err)

		// Case: "Failed to find documents"
	})

	mt.Run("InsertOneDocument", func(mt *mtest.T) {
	})

	mt.Run("InsertManyDocuments", func(mt *mtest.T) {
	})

	mt.Run("UpdateOneDocument", func(mt *mtest.T) {
	})

	mt.Run("DeleteOneDocument", func(mt *mtest.T) {
	})

	mt.Run("DeleteManyDocuments", func(mt *mtest.T) {
	})
}
