package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/tyler-cromwell/forage/models"
)

func getManyDocuments(response http.ResponseWriter, request *http.Request, mc *models.MongoClient) {
	// Specify common fields
	logger := logrus.WithFields(logrus.Fields{"method": "GET"})

	// Attempt to get the documents
	filter := bson.D{{}} // Effectively gets ALL documents
	documents, err := mc.GetManyDocuments(request.Context(), logger, filter)
	if err != nil {
		logger.WithError(err).Error("GetManyDocuments failed")
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		// Prepare to respond with documents
		marshalled, err := json.Marshal(documents)
		if err != nil {
			logger.WithError(err).Error("json.Marshal failed")
			response.WriteHeader(http.StatusInternalServerError)
		} else {
			logger.Debug("json.Marshal succeeded")
			response.WriteHeader(http.StatusOK)
			response.Write(marshalled)
		}
	}
}

func ListenAndServe(tcpSocket string) {
	// Specify common fields
	logger := logrus.WithFields(logrus.Fields{})

	// Initialize context/timeout
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Initialize MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		logrus.Fatal(err)
	}

	// Connect to database instance
	err = client.Connect(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Specify database & collection
	database := client.Database("forage")
	collection := database.Collection("data")
	mc := models.MongoClient{Collection: collection}

	// Define route actions/methods
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		getManyDocuments(response, request, &mc)
	})).Methods("GET")

	logger.WithFields(logrus.Fields{"socket": tcpSocket}).Info("Listening")
	http.ListenAndServe(tcpSocket, router)
}
