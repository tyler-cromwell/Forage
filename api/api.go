package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/tyler-cromwell/forage/clients"
	"github.com/tyler-cromwell/forage/models"
)

var mongoClient *models.MongoClient

func getExpiring(response http.ResponseWriter, request *http.Request) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getExpiring",
		"method": "GET",
	})

	// Filter by food expiring within 2 days
	now := time.Now()
	lookahead := time.Now().Add(time.Hour * 24 * 2)
	filter := bson.M{"$and": []bson.M{
		{
			"$or": []bson.M{
				{
					"expirationDate": bson.M{
						"$gte": primitive.NewDateTimeFromTime(now),
						"$lte": primitive.NewDateTimeFromTime(lookahead),
					},
				},
				{
					"sellBy": bson.M{
						"$gte": primitive.NewDateTimeFromTime(now),
						"$lte": primitive.NewDateTimeFromTime(lookahead),
					},
				},
			},
		},
		{
			"haveStocked": bson.M{
				"$eq": true,
			},
		},
	}}

	// Define sorting criteria
	opts := options.Find()
	opts.SetSort(bson.D{{"expirationDate", -1}})

	// Grab the documents
	documents, err := mongoClient.GetManyDocuments(context.Background(), filter, opts)
	if err != nil {
		log.WithError(err).Error("Failed to identify expiring items")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
	} else {
		// Prepare to respond with documents
		marshalled, err := json.Marshal(documents)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode documents")
			RespondWithError(response, log, http.StatusInternalServerError, err.Error())
		} else {
			log.WithFields(logrus.Fields{"quantity": len(documents), "size": len(marshalled), "status": http.StatusOK}).Debug("Success")
			response.WriteHeader(http.StatusOK)
			response.Write(marshalled)
		}
	}
}

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
	if err != nil && err.Error() == "the provided hex string is not a valid ObjectID" {
		// Invalid document id provided
		log.WithFields(logrus.Fields{"id": id, "status": http.StatusBadRequest}).WithError(err).Warn("Failed to parse document id")
		RespondWithError(response, log, http.StatusBadRequest, err.Error())
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"id": id, "status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse document id")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
		return
	}

	// Create filter
	filter := bson.D{{"_id", oid}}
	log = log.WithFields(logrus.Fields{"filter": filter})

	// Attempt to get the document
	document, err := mongoClient.GetOneDocument(request.Context(), filter)
	if err != nil && err.Error() == "mongo: no documents in result" {
		// Search completed but no document was found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).WithError(err).Warn("Failed to get document")
		RespondWithError(response, log, http.StatusNotFound, err.Error())
	} else if err != nil {
		// Search failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to get document")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
	} else {
		// Prepare to respond with document
		marshalled, err := json.Marshal(document)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode document")
			RespondWithError(response, log, http.StatusInternalServerError, err.Error())
		} else {
			log.WithFields(logrus.Fields{"size": len(marshalled), "status": http.StatusOK}).Debug("Success")
			response.WriteHeader(http.StatusOK)
			response.Write(marshalled)
		}
	}
}

func getManyDocuments(response http.ResponseWriter, request *http.Request) {
	// Extract query parameters
	queryParams := request.URL.Query()
	qpFrom := queryParams.Get("from")
	qpName := queryParams.Get("name")
	qpType := queryParams.Get("type")
	qpTo := queryParams.Get("to")

	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getManyDocuments",
		"method": "GET",
	})

	// Check if query parameters are present
	var timeFrom time.Time
	var timeTo time.Time
	filterExpires := bson.M{}
	filterName := bson.M{}
	filterType := bson.M{}
	filterHaveStocked := bson.M{
		"haveStocked": bson.M{
			"$eq": true,
		},
	}

	if qpFrom != "" {
		from, err := strconv.ParseInt(qpFrom, 10, 64)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse from date")
			RespondWithError(response, log, http.StatusInternalServerError, err.Error())
			return
		}
		timeFrom = time.Unix(0, from*int64(time.Millisecond))
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$gte": primitive.NewDateTimeFromTime(timeFrom),
			},
		}
	}
	if qpName != "" {
		filterName = bson.M{"name": qpName}
	}
	if qpType != "" {
		filterType = bson.M{"type": qpType}
	}
	if qpTo != "" {
		to, err := strconv.ParseInt(qpTo, 10, 64)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse to date")
			RespondWithError(response, log, http.StatusInternalServerError, err.Error())
			return
		}
		timeTo = time.Unix(0, to*int64(time.Millisecond))
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$lte": primitive.NewDateTimeFromTime(timeTo),
			},
		}
	}

	if qpFrom != "" && qpTo != "" {
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$gte": primitive.NewDateTimeFromTime(timeFrom),
				"$lte": primitive.NewDateTimeFromTime(timeTo),
			},
		}
	}

	// Create filter
	filter := bson.M{"$and": []bson.M{
		filterHaveStocked,
		filterName,
		filterType,
		filterExpires,
	}}
	log = log.WithFields(logrus.Fields{"filter": filter})

	// Attempt to get the documents
	documents, err := mongoClient.GetManyDocuments(request.Context(), filter, nil)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to find documents")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
	} else {
		// Prepare to respond with documents
		marshalled, err := json.Marshal(documents)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode documents")
			RespondWithError(response, log, http.StatusInternalServerError, err.Error())
		} else {
			log.WithFields(logrus.Fields{"quantity": len(documents), "size": len(marshalled), "status": http.StatusOK}).Debug("Success")
			response.WriteHeader(http.StatusOK)
			response.Write(marshalled)
		}
	}
}

/*
func postManyDocuments(response http.ResponseWriter, request *http.Request) {

}
*/

func putOneDocument(response http.ResponseWriter, request *http.Request) {
	// Extract route parameter
	vars := mux.Vars(request)
	id := vars["id"]

	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.putOneDocument",
		"method": "PUT",
	})

	// Parse document id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil && err.Error() == "the provided hex string is not a valid ObjectID" {
		// Invalid document id provided
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Warn("Failed to parse document id")
		RespondWithError(response, log, http.StatusBadRequest, err.Error())
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse document id")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
		return
	}

	// Construct filter
	filter := bson.D{{"_id", oid}}
	log = log.WithFields(logrus.Fields{"filter": filter})

	// Get document fields from body
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse request body")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
		return
	}

	// Parse update fields
	var fields map[string]interface{}
	err = json.Unmarshal(bytes, &fields)
	if err != nil && strings.HasPrefix(err.Error(), "invalid character") {
		// Invalid request body
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Warn("Failed to decode update fields")
		RespondWithError(response, log, http.StatusBadRequest, err.Error())
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to decode update fields")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
		return
	}

	// Construct update instructions
	interim := bson.M{}
	for k, v := range fields {
		interim[k] = v
	}
	update := bson.M{"$set": interim}

	// Attempt to put the document
	matched, _, err := mongoClient.PutOneDocument(request.Context(), filter, update)
	if matched == 0 {
		// Update completed but no document was found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).WithError(err).Warn("Failed to update document")
		RespondWithError(response, log, http.StatusNotFound, err.Error())
	} else if err != nil {
		// Update failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to update document")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
	} else {
		// Respond
		log.WithFields(logrus.Fields{"status": http.StatusOK}).Debug("Success")
		response.WriteHeader(http.StatusOK)
	}
}

func deleteOneDocument(response http.ResponseWriter, request *http.Request) {
	// Extract route parameter
	vars := mux.Vars(request)
	id := vars["id"]

	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.deleteOneDocument",
		"method": "GET",
	})

	// Parse document id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse document id")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
		return
	}

	// Create filter
	filter := bson.D{{"_id", oid}}
	log = log.WithFields(logrus.Fields{"filter": filter})

	// Attempt to delete the document
	err = mongoClient.DeleteOneDocument(request.Context(), filter)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to delete document")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
	} else {
		log.WithFields(logrus.Fields{"status": http.StatusOK}).Debug("Success")
		response.WriteHeader(http.StatusOK)
	}
}

/*
func deleteManyDocuments(response http.ResponseWriter, request *http.Request) {

}
*/

func ListenAndServe(tcpSocket string) {
	uri := "mongodb://127.0.0.1:27017"

	// Initialize context/timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	// Launch job to periodically check for expiring food
	interval := 24 * time.Hour
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		lookahead := time.Hour * 24 * 2
		logrus.WithFields(logrus.Fields{"interval": interval, "lookahead": lookahead}).Info("Expiration watch job started")
		for {
			select {
			case <-ticker.C:
				// Specify common fields
				log := logrus.WithFields(logrus.Fields{"at": "api.expirationJob"})

				// Filter by food expiring within 2 days
				now := time.Now()
				later := time.Now().Add(time.Hour * 24 * 2)
				filter := bson.M{"$or": []bson.M{
					{
						"expirationDate": bson.M{
							"$gte": primitive.NewDateTimeFromTime(now),
							"$lte": primitive.NewDateTimeFromTime(later),
						},
					},
					{
						"sellBy": bson.M{
							"$gte": primitive.NewDateTimeFromTime(now),
							"$lte": primitive.NewDateTimeFromTime(later),
						},
					},
				}}

				// Grab the documents
				documents, err := mongoClient.GetManyDocuments(context.Background(), filter, nil)
				if err != nil {
					log.WithError(err).Error("Failed to identify expiring items")
				} else {
					quantity := len(documents)
					log.WithFields(logrus.Fields{"quantity": quantity}).Info("Items expiring")
					// if > 0, push an event (SMS via Twilio? Email? Schedule shopping in Google Calendar?, Prepare a Peapod order?)

					// Skip if nothing is expiring
					if quantity == 0 {
						continue
					}

					// Get message participants
					phoneFrom := os.Getenv("TWILIO_PHONE_FROM")
					phoneTo := os.Getenv("TWILIO_PHONE_TO")

					// Compose message
					var message string
					if quantity == 1 {
						message = fmt.Sprint(quantity, " item expiring soon")
					} else {
						message = fmt.Sprint(quantity, " items expiring soon")
					}

					// Send the message
					_, err := clients.SendMessage(phoneFrom, phoneTo, message)
					if err != nil {
						log.WithFields(logrus.Fields{"from": phoneFrom, "to": phoneTo, "message": message}).WithError(err).Error("Failed to send Twilio message")
					} else {
						log.WithFields(logrus.Fields{"from": phoneFrom, "to": phoneTo, "message": message}).Debug("Sent Twilio message")
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// Define route actions/methods
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/documents/{id}", getOneDocument).Methods("GET")
	router.HandleFunc("/documents/{id}", putOneDocument).Methods("PUT")
	router.HandleFunc("/documents/{id}", deleteOneDocument).Methods("DELETE")
	router.HandleFunc("/documents", getManyDocuments).Methods("GET")
	//router.HandleFunc("/documents", postManyDocuments).Methods("POST")
	//router.HandleFunc("/documents", deleteManyDocuments).Methods("DELETE")
	router.HandleFunc("/expiring", getExpiring).Methods("GET")

	logrus.WithFields(logrus.Fields{"socket": tcpSocket}).Info("Listening")
	err = http.ListenAndServe(tcpSocket, router)
	if err != nil {
		logrus.WithFields(logrus.Fields{"socket": tcpSocket}).WithError(err).Fatal("Failed to listen for and serve requests")
	}
}
