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

	"github.com/adlio/trello"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/tyler-cromwell/forage/clients"
)

var mongoClient *clients.Mongo
var trelloClient clients.Trello
var twilioClient *clients.Twilio

func getExpiring(response http.ResponseWriter, request *http.Request) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{"at": "api.getExpiring"})

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
	documents, err := mongoClient.FindDocuments(context.Background(), filter, opts)
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
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getOneDocument",
		"method": "GET",
	})

	// Extract route parameter
	vars := mux.Vars(request)
	id := vars["id"]

	// Parse document id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil && err.Error() == "the provided hex string is not a valid ObjectID" {
		// Invalid document id provided
		log.WithFields(logrus.Fields{"id": id, "status": http.StatusBadRequest}).WithError(err).Error("Failed to parse document id")
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
	document, err := mongoClient.FindOneDocument(request.Context(), filter)
	if err != nil && err.Error() == "mongo: no documents in result" {
		// Get completed but no document was found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).WithError(err).Warn("Failed to get document")
		RespondWithError(response, log, http.StatusNotFound, err.Error())
	} else if err != nil {
		// Get failed
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
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getManyDocuments",
		"method": "GET",
	})

	// Extract query parameters
	queryParams := request.URL.Query()
	qpFrom := queryParams.Get("from")
	qpHaveStocked := queryParams.Get("haveStocked")
	qpName := queryParams.Get("name")
	qpType := queryParams.Get("type")
	qpTo := queryParams.Get("to")

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
			log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse from date")
			RespondWithError(response, log, http.StatusBadRequest, err.Error())
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
			log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse to date")
			RespondWithError(response, log, http.StatusBadRequest, err.Error())
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

	if qpHaveStocked != "" {
		b, err := strconv.ParseBool(qpHaveStocked)
		if err != nil {
			// Invalid query parameter value provided
			log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse haveStocked")
			RespondWithError(response, log, http.StatusBadRequest, err.Error())
			return
		} else {
			filterHaveStocked = bson.M{
				"haveStocked": bson.M{
					"$eq": b,
				},
			}
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
	documents, err := mongoClient.FindDocuments(request.Context(), filter, nil)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to get documents")
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

func postManyDocuments(response http.ResponseWriter, request *http.Request) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.postManyDocuments",
		"method": "POST",
	})

	// Get documents from body
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse request body")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
		return
	}

	// Parse documents
	var data []interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil && strings.HasPrefix(err.Error(), "invalid character") {
		// Invalid request body
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to decode documents")
		RespondWithError(response, log, http.StatusBadRequest, err.Error())
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to decode documents")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
		return
	}

	// Construct insert instructions
	documents := []interface{}{}
	for _, e := range data {
		documents = append(documents, e)
	}

	// Attempt to put the document
	err = mongoClient.InsertManyDocuments(request.Context(), documents)
	if err != nil {
		// Post failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to post document")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
	} else {
		log.WithFields(logrus.Fields{"status": http.StatusCreated}).Debug("Success")
		response.WriteHeader(http.StatusCreated)
	}
}

func putOneDocument(response http.ResponseWriter, request *http.Request) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.putOneDocument",
		"method": "PUT",
	})

	// Extract route parameter
	vars := mux.Vars(request)
	id := vars["id"]

	// Parse document id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil && err.Error() == "the provided hex string is not a valid ObjectID" {
		// Invalid document id provided
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse document id")
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
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to decode update fields")
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

	// Ignore _id field since it's immutable and will error
	_, ok := interim["_id"]
	if ok {
		delete(interim, "_id")
	}

	update := bson.M{"$set": interim}

	// Attempt to put the document
	matched, _, err := mongoClient.UpdateOneDocument(request.Context(), filter, update)
	if matched == 0 {
		// Put completed but no document was found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).WithError(err).Warn("Failed to put document")
		RespondWithError(response, log, http.StatusNotFound, err.Error())
	} else if err != nil {
		// Put failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to put document")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
	} else {
		log.WithFields(logrus.Fields{"status": http.StatusOK}).Debug("Success")
		response.WriteHeader(http.StatusOK)
	}
}

func deleteOneDocument(response http.ResponseWriter, request *http.Request) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.deleteOneDocument",
		"method": "DELETE",
	})

	// Extract route parameter
	vars := mux.Vars(request)
	id := vars["id"]

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
		// Delete failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to delete document")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
	} else {
		log.WithFields(logrus.Fields{"status": http.StatusOK}).Debug("Success")
		response.WriteHeader(http.StatusOK)
	}
}

func deleteManyDocuments(response http.ResponseWriter, request *http.Request) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.deleteManyDocuments",
		"method": "DELETE",
	})

	// Delete by list of IDs (for now)
	// Get document fields from body
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse request body")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
		return
	}

	// Parse update fields
	var ids []string
	err = json.Unmarshal(bytes, &ids)
	if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to decode delete fields")
		RespondWithError(response, log, http.StatusBadRequest, err.Error())
		return
	}

	// Create filter
	interim := []primitive.ObjectID{}
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse id")
			RespondWithError(response, log, http.StatusBadRequest, err.Error())
			return
		} else {
			interim = append(interim, oid)
		}
	}
	filter := bson.M{"_id": bson.M{"$in": interim}}
	log = log.WithFields(logrus.Fields{"filter": filter})

	// Attempt to delete the documents
	deleted, err := mongoClient.DeleteManyDocuments(request.Context(), filter)
	if deleted == 0 {
		// Delete completed but no documents were found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).WithError(err).Warn("Failed to delete documents")
		RespondWithError(response, log, http.StatusNotFound, fmt.Sprint("no documents found"))
	} else if err != nil {
		// Delete failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to delete documents")
		RespondWithError(response, log, http.StatusInternalServerError, err.Error())
	} else {
		log.WithFields(logrus.Fields{"quantity": deleted, "status": http.StatusOK}).Debug("Success")
		response.WriteHeader(http.StatusOK)
	}
}

func ListenAndServe(tcpSocket string) {
	// Get environment variables
	trelloMemberID := os.Getenv("TRELLO_MEMBER")
	trelloBoardName := os.Getenv("TRELLO_BOARD")
	trelloListName := os.Getenv("TRELLO_LIST")
	trelloLabels := os.Getenv("TRELLO_LABELS")
	trelloApiKey := os.Getenv("TRELLO_API_KEY")
	trelloApiToken := os.Getenv("TRELLO_API_TOKEN")
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	phoneFrom := os.Getenv("TWILIO_PHONE_FROM")
	phoneTo := os.Getenv("TWILIO_PHONE_TO")
	uri := os.Getenv("DATABASE_URI")

	// Initialize context/timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		logrus.WithFields(logrus.Fields{"uri": uri}).WithError(err).Fatal("Failure initialize MongoDB client")
	}

	// Initialize Trello client
	trelloClient = clients.Trello{
		Key:       trelloApiKey,
		Token:     trelloApiToken,
		MemberID:  trelloMemberID,
		BoardName: trelloBoardName,
		ListName:  trelloListName,
		Client:    trello.NewClient(trelloApiKey, trelloApiToken),
	}

	// Initialize Twilio client
	twilioClient = &clients.Twilio{
		From: phoneFrom,
		To:   phoneTo,
		Client: twilio.NewRestClientWithParams(twilio.RestClientParams{
			Username: accountSid,
			Password: authToken,
		})}

	// Connect to database instance
	err = client.Connect(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"uri": uri}).WithError(err).Fatal("Failed to connect to MongoDB instance")
	}
	defer client.Disconnect(ctx)

	// Specify database & collection
	database := client.Database("forage")
	collection := database.Collection("data")
	mongoClient = &clients.Mongo{Collection: collection}

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
				documents, err := mongoClient.FindDocuments(context.Background(), filter, nil)
				if err != nil {
					log.WithError(err).Error("Failed to identify expiring items")
				} else {
					quantity := len(documents)
					log.WithFields(logrus.Fields{"quantity": quantity}).Debug("Items expiring")
					// if > 0, push an event (SMS via Twilio? Email? Schedule shopping in Google Calendar?, Prepare a Peapod order?)

					// Skip if nothing is expiring
					if quantity == 0 {
						continue
					}

					// Construct list of names of items to shop for
					var groceries []string
					for _, document := range documents {
						v, keyFound := document["name"]
						if keyFound {
							groceries = append(groceries, v.(string))
						}
					}

					// Construct shopping list due date
					now := time.Now()
					rounded := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
					dueDate := rounded.Add(lookahead + (time.Hour * 24))

					// Create shopping list card on Trello
					url, err := trelloClient.CreateShoppingList(&dueDate, groceries)
					if err != nil {
						log.WithError(err).Error("Failed to create Trello card")
					} else {
						log.WithFields(logrus.Fields{"url": url}).Debug("Created Trello card")
					}

					// Compose message
					var message string
					if quantity == 1 {
						message = fmt.Sprintf("%d item expiring soon! View shopping list: %s", quantity, url)
					} else {
						message = fmt.Sprintf("%d items expiring soon! View shopping list: %s", quantity, url)
					}

					// Send the message
					_, err = twilioClient.SendMessage(phoneFrom, phoneTo, message)
					if err != nil {
						log.WithFields(logrus.Fields{"from": phoneFrom, "to": phoneTo}).WithError(err).Error("Failed to send Twilio message")
					} else {
						log.WithFields(logrus.Fields{"from": phoneFrom, "to": phoneTo}).Debug("Sent Twilio message")
					}
				}
			case <-quit:
				ticker.Stop()
				logrus.WithFields(logrus.Fields{"interval": interval, "lookahead": lookahead}).Info("Expiration watch job stopped")
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
	router.HandleFunc("/documents", postManyDocuments).Methods("POST")
	router.HandleFunc("/documents", deleteManyDocuments).Methods("DELETE")
	router.HandleFunc("/expiring", getExpiring).Methods("GET")

	logrus.WithFields(logrus.Fields{"socket": tcpSocket}).Info("Listening for HTTP requests")
	err = http.ListenAndServe(tcpSocket, router)
	if err != nil {
		logrus.WithFields(logrus.Fields{"socket": tcpSocket}).WithError(err).Fatal("Failed to listen for and serve requests")
	}
}
