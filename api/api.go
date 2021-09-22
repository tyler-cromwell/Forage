package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/tyler-cromwell/forage/models"
)

var mongoClient *models.MongoClient

func getOneDocument(response http.ResponseWriter, request *http.Request) {
	// Extract route parameter
	vars := mux.Vars(request)
	id := vars["id"]

	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getOneDocument",
		"method": "GET",
	})

	// Parse document id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse document id")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create filter
	filter := bson.D{{"_id", oid}}
	log = log.WithFields(logrus.Fields{
		"filter": filter,
	})

	// Attempt to get the document
	document, err := mongoClient.GetOneDocument(request.Context(), filter)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to find document")
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		// Prepare to respond with document
		marshalled, err := json.Marshal(document)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode document")
			response.WriteHeader(http.StatusInternalServerError)
		} else {
			log.WithFields(logrus.Fields{"size": len(marshalled), "status": http.StatusOK}).Info("Success")
			response.WriteHeader(http.StatusOK)
			response.Write(marshalled)
		}
	}
}

func getManyDocuments(response http.ResponseWriter, request *http.Request) {
	// Extract query parameters
	queryParams := request.URL.Query()
	qpName := queryParams.Get("name")
	qpType := queryParams.Get("type")

	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getManyDocuments",
		"method": "GET",
	})

	filterName := bson.D{{}}
	filterType := bson.D{{}}
	if qpName != "" {
		filterName = bson.D{{"name", qpName}}
	}
	if qpType != "" {
		filterType = bson.D{{"type", qpType}}
	}

	filter := bson.D{{"$and", []bson.D{
		filterName,
		filterType,
	}}}
	log = log.WithFields(logrus.Fields{
		"filter": filter,
	})

	// Attempt to get the documents
	documents, err := mongoClient.GetManyDocuments(request.Context(), filter)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to find documents")
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		// Prepare to respond with documents
		marshalled, err := json.Marshal(documents)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode documents")
			response.WriteHeader(http.StatusInternalServerError)
		} else {
			log.WithFields(logrus.Fields{"quantity": len(documents), "size": len(marshalled), "status": http.StatusOK}).Info("Success")
			response.WriteHeader(http.StatusOK)
			response.Write(marshalled)
		}
	}
}

func ListenAndServe(tcpSocket string) {
	uri := "mongodb://127.0.0.1:27017"

	// Initialize context/timeout
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Initialize MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		logrus.WithFields(logrus.Fields{"uri": uri}).WithError(err).Fatal("Failure initialize MongoDB client")
	}

	// Connect to database instance
	err = client.Connect(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"uri": uri}).WithError(err).Fatal("Failed to connect to MongoDB instance")
	}
	defer client.Disconnect(ctx)

	// Specify database & collection
	database := client.Database("forage")
	collection := database.Collection("data")
	mongoClient = &models.MongoClient{Collection: collection}

	// Define route actions/methods
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/documents/{id}", getOneDocument).Methods("GET")
	router.HandleFunc("/documents", getManyDocuments).Methods("GET")

	logrus.WithFields(logrus.Fields{"socket": tcpSocket}).Info("Listening")
	err = http.ListenAndServe(tcpSocket, router)
	if err != nil {
		logrus.WithFields(logrus.Fields{"socket": tcpSocket}).WithError(err).Fatal("Failed to listen for and serve requests")
	}
}
