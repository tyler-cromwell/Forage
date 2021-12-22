package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tyler-cromwell/forage/config"
	"github.com/tyler-cromwell/forage/utils"
)

var configuration *config.Configuration

func getExpired(response http.ResponseWriter, request *http.Request) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{"at": "api.getExpired"})

	// Filter by food expired already
	filter := bson.M{"$and": []bson.M{
		{
			"expirationDate": bson.M{
				"$lte": time.Now().UnixNano() / int64(time.Millisecond),
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
	opts.SetSort(bson.D{{"expirationDate", 1}})

	// Grab the documents
	documents, err := configuration.Mongo.FindDocuments(context.Background(), filter, opts)
	if err != nil {
		log.WithError(err).Error("Failed to identify expired items")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else {
		// Prepare to respond with documents
		marshalled, err := json.Marshal(documents)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode documents")
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(err.Error()))
		} else {
			log.WithFields(logrus.Fields{"quantity": len(documents), "size": len(marshalled), "status": http.StatusOK}).Debug("Success")
			response.WriteHeader(http.StatusOK)
			response.Write(marshalled)
		}
	}
}

func getExpiring(response http.ResponseWriter, request *http.Request) {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{"at": "api.getExpiring"})

	// Extract query parameters
	queryParams := request.URL.Query()
	qpFrom := queryParams.Get("from")
	qpTo := queryParams.Get("to")

	// Check if query parameters are present
	var timeFrom time.Time = time.Now()
	var timeTo time.Time = time.Now().Add(configuration.Lookahead)
	filterExpires := bson.M{
		"expirationDate": bson.M{
			"$gte": timeFrom.UnixNano() / int64(time.Millisecond),
			"$lte": timeTo.UnixNano() / int64(time.Millisecond),
		},
	}

	if qpFrom != "" {
		from, err := strconv.ParseInt(qpFrom, 10, 64)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse from date")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(err.Error()))
			return
		}
		timeFrom = time.Unix(0, from*int64(time.Millisecond))
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$gte": timeFrom.UnixNano() / int64(time.Millisecond),
			},
		}
	}
	if qpTo != "" {
		to, err := strconv.ParseInt(qpTo, 10, 64)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse to date")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(err.Error()))
			return
		}
		timeTo = time.Unix(0, to*int64(time.Millisecond))
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$lte": timeTo.UnixNano() / int64(time.Millisecond),
			},
		}
	}

	if qpFrom != "" && qpTo != "" {
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$gte": timeFrom.UnixNano() / int64(time.Millisecond),
				"$lte": timeTo.UnixNano() / int64(time.Millisecond),
			},
		}
	}

	// Filter by food expiring within the given search window
	filter := bson.M{"$and": []bson.M{
		filterExpires,
		{
			"haveStocked": bson.M{
				"$eq": true,
			},
		},
	}}

	// Define sorting criteria
	opts := options.Find()
	opts.SetSort(bson.D{{"expirationDate", 1}})

	// Grab the documents
	documents, err := configuration.Mongo.FindDocuments(context.Background(), filter, opts)
	if err != nil {
		log.WithError(err).Error("Failed to identify expiring items")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else {
		// Prepare to respond with documents
		marshalled, err := json.Marshal(documents)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode documents")
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(err.Error()))
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
	if err != nil && err.Error() == utils.ErrInvalidObjectID {
		// Invalid document id provided
		log.WithFields(logrus.Fields{"id": id, "status": http.StatusBadRequest}).WithError(err).Warn("Failed to parse document id")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"id": id, "status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse document id")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	// Create filter
	filter := bson.D{{"_id", oid}}
	log = log.WithFields(logrus.Fields{"filter": filter})

	// Attempt to get the document
	document, err := configuration.Mongo.FindOneDocument(request.Context(), filter)
	if err != nil && err.Error() == utils.ErrMongoNoDocuments {
		// Get completed but no document was found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).WithError(err).Warn("Failed to get document")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(err.Error()))
	} else if err != nil {
		// Get failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to get document")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else {
		// Prepare to respond with document
		marshalled, err := json.Marshal(document)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode document")
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(err.Error()))
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
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(err.Error()))
			return
		}
		timeFrom = time.Unix(0, from*int64(time.Millisecond))
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$gte": timeFrom.UnixNano() / int64(time.Millisecond),
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
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(err.Error()))
			return
		}
		timeTo = time.Unix(0, to*int64(time.Millisecond))
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$lte": timeTo.UnixNano() / int64(time.Millisecond),
			},
		}
	}

	if qpFrom != "" && qpTo != "" {
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$gte": timeFrom.UnixNano() / int64(time.Millisecond),
				"$lte": timeTo.UnixNano() / int64(time.Millisecond),
			},
		}
	}

	if qpHaveStocked != "" {
		b, err := strconv.ParseBool(qpHaveStocked)
		if err != nil {
			// Invalid query parameter value provided
			log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse haveStocked")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(err.Error()))
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
	documents, err := configuration.Mongo.FindDocuments(request.Context(), filter, nil)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to get documents")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else {
		// Prepare to respond with documents
		marshalled, err := json.Marshal(documents)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode documents")
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(err.Error()))
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
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	// Parse documents
	var data []interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil && strings.HasPrefix(err.Error(), "invalid character") {
		// Invalid request body
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to decode documents")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to decode documents")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	// Construct insert instructions
	documents := []interface{}{}
	for _, e := range data {
		documents = append(documents, e)
	}

	// Attempt to put the document
	err = configuration.Mongo.InsertManyDocuments(request.Context(), documents)
	if err != nil {
		// Post failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to post documents")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else {
		log.WithFields(logrus.Fields{"quantity": len(documents), "status": http.StatusCreated}).Debug("Success")
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
	if err != nil && err.Error() == utils.ErrInvalidObjectID {
		// Invalid document id provided
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Warn("Failed to parse document id")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse document id")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	// Construct filter
	filter := bson.D{{"_id", oid}}
	log = log.WithFields(logrus.Fields{"filter": filter})

	// Get document fields from body
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse request body")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	// Parse update fields
	var fields map[string]interface{}
	err = json.Unmarshal(bytes, &fields)
	if err != nil && strings.HasPrefix(err.Error(), "invalid character") {
		// Invalid request body
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to decode update fields")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to decode update fields")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
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
	matched, _, err := configuration.Mongo.UpdateOneDocument(request.Context(), filter, update)
	if matched == 0 {
		// Put completed but no document was found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).WithError(err).Warn("Failed to put document")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(err.Error()))
	} else if err != nil {
		// Put failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to put document")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
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
	if err != nil && err.Error() == utils.ErrInvalidObjectID {
		// Invalid document id provided
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Warn("Failed to parse document id")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse document id")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	// Create filter
	filter := bson.D{{"_id", oid}}
	log = log.WithFields(logrus.Fields{"filter": filter})

	// Attempt to delete the document
	err = configuration.Mongo.DeleteOneDocument(request.Context(), filter)
	if err != nil && err.Error() == utils.ErrMongoNoDocuments {
		// Get completed but no document was found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).WithError(err).Warn("Failed to get document")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(err.Error()))
	} else if err != nil {
		// Get failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to get document")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
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
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	// Parse update fields
	var ids []string
	err = json.Unmarshal(bytes, &ids)
	if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to decode delete fields")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}

	// Create filter
	interim := []primitive.ObjectID{}
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse id")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(err.Error()))
			return
		} else {
			interim = append(interim, oid)
		}
	}
	filter := bson.M{"_id": bson.M{"$in": interim}}
	log = log.WithFields(logrus.Fields{"filter": filter})

	// Attempt to delete the documents
	deleted, err := configuration.Mongo.DeleteManyDocuments(request.Context(), filter)
	if deleted == 0 {
		// Delete completed but no documents were found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).WithError(err).Warn("Failed to delete documents")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte("no documents found"))
	} else if err != nil {
		// Delete failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to delete documents")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else {
		log.WithFields(logrus.Fields{"quantity": deleted, "status": http.StatusOK}).Debug("Success")
		response.WriteHeader(http.StatusOK)
	}
}

func checkExpirations() {
	// Define a context
	ctx := context.Background()

	// Specify common fields
	log := logrus.WithFields(logrus.Fields{"at": "api.checkExpirations"})

	// Filter by food expired already
	now := time.Now().UnixNano() / int64(time.Millisecond)
	later := time.Now().Add(configuration.Lookahead).UnixNano() / int64(time.Millisecond)
	filterExpired := bson.M{"$and": []bson.M{
		{
			"expirationDate": bson.M{
				"$lte": now,
			},
		},
		{
			"haveStocked": bson.M{
				"$eq": true,
			},
		},
	}}

	// Filter by food expiring within the given search window
	filter := bson.M{"$and": []bson.M{
		{
			"expirationDate": bson.M{
				"$gte": now,
				"$lte": later,
			},
		},
		{
			"haveStocked": bson.M{
				"$eq": true,
			},
		},
	}}

	// Grab the documents
	documentsExpired, err := configuration.Mongo.FindDocuments(ctx, filterExpired, nil)
	if err != nil {
		log.WithError(err).Error("Failed to identify expired items")
		return
	}

	documents, err := configuration.Mongo.FindDocuments(ctx, filter, nil)
	if err != nil {
		log.WithError(err).Error("Failed to identify expiring items")
	} else {
		quantityExpired := len(documentsExpired)
		quantity := len(documents)

		// Skip if nothing is expiring
		if quantity == 0 && quantityExpired == 0 {
			log.WithFields(logrus.Fields{"expiring": quantity, "expired": quantityExpired}).Info("Restocking not required")
			return
		} else {
			log.WithFields(logrus.Fields{"expiring": quantity, "expired": quantityExpired}).Info("Restocking required")
		}

		// Construct list of names of items to shop for
		var groceries []string
		groceries = append(groceries, utils.StringSliceFromBsonM(documentsExpired, "name")...)
		groceries = append(groceries, utils.StringSliceFromBsonM(documents, "name")...)

		// Construct shopping list due date
		now := time.Now()
		rounded := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		dueDate := rounded.Add(configuration.Lookahead + (time.Hour * 24))

		var url string
		shoppingListCard, err := configuration.Trello.GetShoppingList()
		if err != nil {
			log.WithError(err).Error("Failed to get Trello card")
		} else if shoppingListCard != nil {
			// Add to shopping list card on Trello
			url, err = configuration.Trello.AddToShoppingList(groceries)
			if err != nil {
				log.WithError(err).Error("Failed to add to Trello card")
			} else {
				log.WithFields(logrus.Fields{"url": url}).Info("Added to Trello card")
			}
		} else {
			// Create shopping list card on Trello
			labels := strings.Split(configuration.Trello.LabelsStr, ",")
			url, err = configuration.Trello.CreateShoppingList(&dueDate, labels, groceries)
			if err != nil {
				log.WithError(err).Error("Failed to create Trello card")
			} else {
				log.WithFields(logrus.Fields{"url": url}).Info("Created Trello card")
			}
		}

		// Compose Twilio message
		var message = configuration.Twilio.ComposeMessage(quantity, quantityExpired, url)

		// Send the Twilio message
		_, err = configuration.Twilio.SendMessage(configuration.Twilio.From, configuration.Twilio.To, message)
		if err != nil {
			log.WithFields(logrus.Fields{"from": configuration.Twilio.From, "to": configuration.Twilio.To}).WithError(err).Error("Failed to send Twilio message")
		} else {
			log.WithFields(logrus.Fields{"from": configuration.Twilio.From, "to": configuration.Twilio.To}).Info("Sent Twilio message")
		}
	}
}

func ListenAndServe(ctx context.Context, c *config.Configuration) {
	configuration = c
	// Launch job to periodically check for expiring food
	ticker := time.NewTicker(configuration.Interval)
	quit := make(chan struct{})
	checkExpirations() // Run once before first ticker tick
	go func() {
		// Specify common fields
		log := logrus.WithFields(logrus.Fields{
			"at":        "expirationJob",
			"interval":  configuration.Interval,
			"lookahead": configuration.Lookahead,
		})

		// Wait for ticker ticks
		log.Info("Expiration watch job started")
		for {
			select {
			case <-ticker.C:
				checkExpirations()
			case <-quit:
				ticker.Stop()
				log.Info("Expiration watch job stopped")
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
	router.HandleFunc("/expired", getExpired).Methods("GET")

	// Specify common fields
	log := logrus.WithFields(logrus.Fields{"socket": configuration.ListenSocket})

	// Listen for HTTP requests
	log.Info("Listening for HTTP requests")
	err := http.ListenAndServe(configuration.ListenSocket, router)
	if err != nil {
		log.WithError(err).Fatal("Failed to listen for and serve requests")
	}
}
