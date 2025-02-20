package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func isCookable(ctx context.Context, recipe *primitive.M) (bool, error) {
	// Setup
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.isCookable",
		"recipe": (*recipe)["_id"],
	})

	// Determine if recipe is cookable
	log.Trace("Begin cookable determination")
	defer log.Trace("End cookable determination")
	ingredientNames := (*recipe)["ingredients"]

	if ingredientNames == nil {
		return false, fmt.Errorf("no ingredients specified")
	}

	var length int
	if v, ok := ingredientNames.([]interface{}); ok {
		length = len(v)
	} else if v, ok := ingredientNames.(primitive.A); ok {
		length = len(v)
	}
	//iids := ingreidentNames.([]interface{}) // "ingredients" is a requried attribute
	filterMany := bson.M{"$and": []bson.M{
		{
			"expirationDate": bson.M{
				"$gt": int64(time.Now().UTC().UnixNano()) / int64(time.Millisecond),
			},
		},
		{
			"haveStocked": bson.M{
				"$eq": true,
			},
		},
		{
			"name": bson.M{
				"$in": ingredientNames,
			},
		},
	}}
	log.WithFields(logrus.Fields{"value": filterMany}).Debug("Filter data")

	ingredients, err := configuration.Mongo.FindDocuments(ctx, config.MongoCollectionIngredients, filterMany, nil)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to get documents")
		return false, err
	} else {
		result := len(ingredients) >= length
		log.WithFields(logrus.Fields{"expect": length, "have": len(ingredients), "value": result}).Debug("Determined")
		return result, nil
	}
}
