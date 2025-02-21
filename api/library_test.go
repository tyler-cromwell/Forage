package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/tests/mocks"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestLibrary(t *testing.T) {
	// Discard logging output
	logrus.SetOutput(io.Discard)

	t.Run("isCookable", func(t *testing.T) {
		ctx := context.Background()

		mcErr := mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}
		mc := mocks.MockMongo{
			OverrideFindManyDocuments: func(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return []primitive.M{}, nil
			},
		}

		cases := []struct {
			mock   mocks.MockMongo
			recipe primitive.M
			want   bool
			err    error
		}{
			{mc, primitive.M{"_id": "hello"}, false, nil},
			{mcErr, primitive.M{"_id": "hello", "ingredients": []interface{}{}}, false, fmt.Errorf(errorBasic)},
			{mc, primitive.M{"_id": "hello", "ingredients": []interface{}{}}, true, nil},
		}
		for _, c := range cases {
			configuration.Mongo = &c.mock
			got, err := isCookable(ctx, &c.recipe)
			if got != c.want {
				t.Errorf("isCookable(\"%+v\"), got (\"%t\", \"%s\"), want (\"%t\", \"%s\")", c.recipe, got, err, c.want, c.err)
			}
		}
	})

	// Reverse logrus output change
	log.SetOutput(os.Stdout)
}
