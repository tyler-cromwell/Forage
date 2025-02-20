package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/config"
	"go.mongodb.org/mongo-driver/bson"
)

func getCookable(response http.ResponseWriter, request *http.Request) {
	// Setup
	ctx := request.Context()
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getCookable",
		"method": "GET",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Create filter
	filter := bson.M{"isCookable": true}
	log.WithFields(logrus.Fields{"value": filter}).Debug("Filter data")

	// Grab the documents
	documents, err := configuration.Mongo.FindDocuments(ctx, config.MongoCollectionRecipes, filter, nil)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to identify cookable recipes")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else {
		log.WithFields(logrus.Fields{"quantity": len(documents), "value": documents}).Debug("Documents found")

		// Prepare to respond with documents
		marshalled, err := json.Marshal(documents)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode documents")
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(err.Error()))
		} else {
			log.WithFields(logrus.Fields{"quantity": len(documents), "size": len(marshalled), "status": http.StatusOK}).Info("Succeeded")
			response.WriteHeader(http.StatusOK)
			response.Write(marshalled)
		}
	}
}
