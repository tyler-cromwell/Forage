package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
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

func getConfiguration(response http.ResponseWriter, request *http.Request) {
	// Setup
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getConfiguration",
		"method": "GET",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Prepare the response data
	marshalled, _ := json.Marshal(struct {
		Lookahead time.Duration `json:"lookahead"`
		Silence   bool          `json:"silence"`
		Time      string        `json:"time"`
	}{
		configuration.Lookahead,
		configuration.Silence,
		configuration.Time,
	})

	// Log & Respond
	log.WithFields(logrus.Fields{"size": len(marshalled), "state": "marshalled", "value": string(marshalled)}).Debug("Response body")
	log.WithFields(logrus.Fields{"status": http.StatusOK}).Info("Succeeded")
	response.WriteHeader(http.StatusOK)
	response.Write(marshalled)
}

func putConfiguration(response http.ResponseWriter, request *http.Request) {
	// Setup
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.putConfiguration",
		"method": "PUT",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Read in request body
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to read request body")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"size": len(bytes), "state": "marshalled", "value": string(bytes)}).Debug("Request body")
	}

	// Parse request body
	var body struct {
		Lookahead time.Duration `json:"lookahead"`
		Silence   bool          `json:"silence"`
		Time      string        `json:"time"`
	}
	err = json.Unmarshal(bytes, &body)
	if err != nil && strings.HasPrefix(err.Error(), "invalid character") {
		// Invalid request body
		log.WithFields(logrus.Fields{"case": 1}).Trace("Invalid request body")
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Warn("Failed to decode update fields")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	} else if err != nil && err.Error() == "unexpected end of JSON input" {
		// Invalid request body
		log.WithFields(logrus.Fields{"case": 2}).Trace("Invalid request body")
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Warn("Failed to decode update fields")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to decode update fields")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"state": "unmarshalled", "value": body}).Debug("Request body")
		log.WithFields(logrus.Fields{"lookahead": body.Lookahead, "silence": body.Silence, "time": body.Time}).Debug("Parsed data")

		if body.Time != "" && len(strings.Split(body.Time, ":")) != 2 {
			log.WithFields(logrus.Fields{"status": http.StatusBadRequest, "value": body.Time}).Warn("Invalid time format")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("Invalid time format: " + body.Time))
			return
		}

		if len(body.Time) != 0 && configuration.Time != body.Time {
			// Re-schedule expiration job
			configuration.Scheduler.Clear()
			log.Info("Scheduler cleared")
			_, err = configuration.Scheduler.Every(1).Day().At(body.Time).Do(checkExpirations)
			if err != nil {
				log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to schedule expriation watch job")
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte(err.Error()))
				return
			} else {
				configuration.Scheduler.StartAsync()
				log.Info("Expiration watch job scheduled")
			}
		}

		configuration.Lookahead = body.Lookahead
		configuration.Silence = body.Silence
		configuration.Time = body.Time

		log.WithFields(logrus.Fields{"status": http.StatusOK}).Info("Succeeded")
		response.WriteHeader(http.StatusOK)
	}
}

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

func getExpired(response http.ResponseWriter, request *http.Request) {
	// Setup
	ctx := request.Context()
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getExpired",
		"method": "GET",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Filter by food expired already
	filter := bson.M{"$and": []bson.M{
		{
			"expirationDate": bson.M{
				"$lte": int64(time.Now().UTC().UnixNano()) / int64(time.Millisecond),
			},
		},
		{
			"haveStocked": bson.M{
				"$eq": true,
			},
		},
	}}
	log.WithFields(logrus.Fields{"value": filter}).Debug("Filter data")

	// Define sorting criteria
	opts := options.Find()
	opts.SetSort(bson.D{{"expirationDate", 1}})
	log.WithFields(logrus.Fields{"value": opts.Sort}).Debug("Sorting criteria")

	// Grab the documents
	documents, err := configuration.Mongo.FindDocuments(ctx, config.MongoCollectionIngredients, filter, opts)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to identify expired items")
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

func getExpiring(response http.ResponseWriter, request *http.Request) {
	// Setup
	ctx := request.Context()
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getExpiring",
		"method": "GET",
	})
	qpNameFrom := "from"
	qpNameTo := "to"

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Extract query parameters
	queryParams := request.URL.Query()
	qpFrom := queryParams.Get(qpNameFrom)
	qpTo := queryParams.Get(qpNameTo)
	log.WithFields(logrus.Fields{"value": queryParams}).Debug("Query parameters")

	// Check if query parameters are present
	var timeFrom time.Time = time.Now()
	var timeTo time.Time = time.Now().Add(configuration.Lookahead)
	filterExpires := bson.M{
		"expirationDate": bson.M{
			"$gte": int64(timeFrom.UTC().UnixNano()) / int64(time.Millisecond),
			"$lte": int64(timeTo.UTC().UnixNano()) / int64(time.Millisecond),
		},
	}

	if qpFrom != "" {
		l := log.WithFields(logrus.Fields{"name": qpNameFrom, "value": qpFrom})
		l.Trace("Query parameter handling")
		from, err := strconv.ParseInt(qpFrom, 10, 64)
		if err != nil {
			l.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse number")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(err.Error()))
			return
		}
		timeFrom = time.Unix(0, from*int64(time.Millisecond))
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$gte": int64(timeFrom.UTC().UnixNano()) / int64(time.Millisecond),
			},
		}
	}
	if qpTo != "" {
		l := log.WithFields(logrus.Fields{"name": qpNameTo, "value": qpTo})
		l.Trace("Query parameter handling")
		to, err := strconv.ParseInt(qpTo, 10, 64)
		if err != nil {
			l.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse number")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(err.Error()))
			return
		}
		timeTo = time.Unix(0, to*int64(time.Millisecond))
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$lte": int64(timeTo.UTC().UnixNano()) / int64(time.Millisecond),
			},
		}
	}

	// Is this even necessary?
	if qpFrom != "" && qpTo != "" {
		log.WithFields(logrus.Fields{"values": []string{qpNameFrom, qpNameTo}}).Trace("Query parameters handling")
		filterExpires = bson.M{
			"expirationDate": bson.M{
				"$gte": int64(timeFrom.UTC().UnixNano()) / int64(time.Millisecond),
				"$lte": int64(timeTo.UTC().UnixNano()) / int64(time.Millisecond),
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
	log.WithFields(logrus.Fields{"value": filter}).Debug("Filter data")

	// Define sorting criteria
	opts := options.Find()
	opts.SetSort(bson.D{{"expirationDate", 1}})
	log.WithFields(logrus.Fields{"value": opts.Sort}).Debug("Sorting criteria")

	// Grab the documents
	documents, err := configuration.Mongo.FindDocuments(ctx, config.MongoCollectionIngredients, filter, opts)
	if err != nil {
		log.WithError(err).Error("Failed to identify expiring items")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else {
		l := log.WithFields(logrus.Fields{"quantity": len(documents)})
		l.WithFields(logrus.Fields{"value": documents}).Debug("Documents found")

		// Prepare to respond with documents
		marshalled, err := json.Marshal(documents)
		if err != nil {
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode documents")
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(err.Error()))
		} else {
			l.WithFields(logrus.Fields{"size": len(marshalled), "status": http.StatusOK}).Info("Succeeded")
			response.WriteHeader(http.StatusOK)
			response.Write(marshalled)
		}
	}
}

func getOneDocument(response http.ResponseWriter, request *http.Request) {
	// Setup
	ctx := request.Context()
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getOneDocument",
		"method": "GET",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Extract route parameters
	vars := mux.Vars(request)
	collection := vars["collection"]
	id := vars["id"]
	log.WithFields(logrus.Fields{"value": vars}).Debug("Route variables")
	log = log.WithFields(logrus.Fields{"collection": collection})

	// Validate collection route variable
	collections, err := configuration.Mongo.Collections(ctx)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed list collection names")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else if !utils.Contains(collections, collection) {
		err := fmt.Errorf("collection not found: %s", collection)
		log.WithFields(logrus.Fields{"collections": collections, "status": http.StatusNotFound}).WithError(err).Warn("Failed to find collection")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(err.Error()))
		return
	}

	// Parse document id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil && err.Error() == utils.ErrorInvalidObjectID {
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
	log.WithFields(logrus.Fields{"value": filter}).Debug("Filter data")
	log = log.WithFields(logrus.Fields{"id": id})

	// Attempt to get the document
	document, err := configuration.Mongo.FindOneDocument(ctx, collection, filter)
	if err != nil && err.Error() == utils.ErrorMongoNoDocuments {
		// Get completed but no document was found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).WithError(err).Warn("Failed to get document")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(err.Error()))
		return
	} else if err != nil {
		// Get failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to get document")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"value": document}).Debug("Document found")
	}

	// Check if document is a recipe
	if collection == config.MongoCollectionRecipes {
		log.Trace("Begin recipe scan")
		// Check if recipe can be made (i.e. associated ingredients are stocked and not expiring)
		originalCanMake := (*document)["isCookable"].(bool)
		isCookable, err := isCookable(ctx, document)
		if err != nil {
			// Something broke
			log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to determine cookable")
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(err.Error()))
			return
		} else if isCookable != originalCanMake {
			// Update if different
			log.WithFields(logrus.Fields{"original": originalCanMake, "updated": isCookable}).Debug("Updating isCookable")
			(*document)["isCookable"] = isCookable

			// Create filter
			l := log.WithFields(logrus.Fields{"method": "PUT"})
			l.WithFields(logrus.Fields{"value": filter}).Debug("Filter data")

			// Define update instructions
			update := bson.M{"$set": bson.M{"isCookable": isCookable}}
			l.WithFields(logrus.Fields{"value": update}).Debug("Update instructions")

			// Attempt to put the document
			matched, _, err := configuration.Mongo.UpdateOneDocument(ctx, collection, filter, update)
			if err != nil {
				// Put failed
				l.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to put document")
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte(err.Error()))
				return
			} else {
				l.WithFields(logrus.Fields{"method": "PUT", "quantity": matched, "status": http.StatusOK}).Info("Succeeded")
			}
		}
		log.Trace("End recipe scan")
	}

	// Prepare to respond with document
	marshalled, err := json.Marshal(document)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to encode document")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else {
		log.WithFields(logrus.Fields{"size": len(marshalled), "status": http.StatusOK}).Info("Succeeded")
		response.WriteHeader(http.StatusOK)
		response.Write(marshalled)
	}
}

func getManyDocuments(response http.ResponseWriter, request *http.Request) {
	// Setup
	ctx := request.Context()
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getManyDocuments",
		"method": "GET",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Extract route parameters
	vars := mux.Vars(request)
	collection := vars["collection"]
	log.WithFields(logrus.Fields{"value": vars}).Debug("Route variables")
	log = log.WithFields(logrus.Fields{"collection": collection})

	// Validate collection route variable
	collections, err := configuration.Mongo.Collections(ctx)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed list collection names")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else if !utils.Contains(collections, collection) {
		err := fmt.Errorf("collection not found: %s", collection)
		log.WithFields(logrus.Fields{"collections": collections, "status": http.StatusNotFound}).WithError(err).Warn("Failed to find collection")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(err.Error()))
		return
	}

	// Extract query parameters
	qpNameCookable := "isCookable"
	qpNameFrom := "from"
	qpNameHaveStocked := "haveStocked"
	qpNameName := "name"
	qpNameTo := "to"
	queryParams := request.URL.Query()
	qpName := queryParams.Get(qpNameName)
	log.WithFields(logrus.Fields{"value": queryParams}).Debug("Query parameters")

	// Check if query parameters are present
	var filter bson.M
	filterName := bson.M{}
	if qpName != "" {
		l := log.WithFields(logrus.Fields{"name": qpNameName, "value": qpName})
		l.Trace("Query parameter handling")
		filterName = bson.M{"name": qpName}
	}

	if collection == config.MongoCollectionRecipes {
		qpCookable := queryParams.Get(qpNameCookable)
		filterIsCookable := bson.M{}

		if qpCookable != "" {
			l := log.WithFields(logrus.Fields{"name": qpNameCookable, "value": qpCookable})
			l.Trace("Query parameter handling")
			b, err := strconv.ParseBool(qpCookable)
			if err != nil {
				// Invalid query parameter value provided
				l.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse isCookable")
				response.WriteHeader(http.StatusBadRequest)
				response.Write([]byte(err.Error()))
				return
			} else {
				filterIsCookable = bson.M{
					"isCookable": bson.M{
						"$eq": b,
					},
				}
			}
		}

		// Create filter
		filter = bson.M{"$and": []bson.M{
			filterName,
			filterIsCookable,
		}}
	} else {
		var timeFrom time.Time
		var timeTo time.Time

		qpFrom := queryParams.Get(qpNameFrom)
		qpHaveStocked := queryParams.Get(qpNameHaveStocked)
		qpTo := queryParams.Get(qpNameTo)

		filterExpires := bson.M{}
		filterType := bson.M{}
		filterHaveStocked := bson.M{
			"haveStocked": bson.M{
				"$eq": true,
			},
		}

		if qpFrom != "" {
			l := log.WithFields(logrus.Fields{"name": qpNameFrom, "value": qpFrom})
			l.Trace("Query parameter handling")
			from, err := strconv.ParseInt(qpFrom, 10, 64)
			if err != nil {
				l.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse from date")
				response.WriteHeader(http.StatusBadRequest)
				response.Write([]byte(err.Error()))
				return
			}
			l.WithFields(logrus.Fields{"old": timeFrom.UnixNano() / int64(time.Millisecond), "new": int64(timeFrom.UTC().UnixNano()) / int64(time.Millisecond)}).Debug("Something")
			timeFrom = time.Unix(0, from*int64(time.Millisecond))
			filterExpires = bson.M{
				"expirationDate": bson.M{
					"$gte": int64(timeFrom.UTC().UnixNano()) / int64(time.Millisecond),
				},
			}
		}
		if qpTo != "" {
			l := log.WithFields(logrus.Fields{"name": qpNameTo, "value": qpTo})
			l.Trace("Query parameter handling")
			to, err := strconv.ParseInt(qpTo, 10, 64)
			if err != nil {
				l.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse to date")
				response.WriteHeader(http.StatusBadRequest)
				response.Write([]byte(err.Error()))
				return
			}
			timeTo = time.Unix(0, to*int64(time.Millisecond))
			filterExpires = bson.M{
				"expirationDate": bson.M{
					"$lte": int64(timeTo.UTC().UnixNano()) / int64(time.Millisecond),
				},
			}
		}

		if qpFrom != "" && qpTo != "" {
			log.WithFields(logrus.Fields{"values": []string{qpNameFrom, qpNameTo}}).Trace("Query parameters handling")
			filterExpires = bson.M{
				"expirationDate": bson.M{
					"$gte": int64(timeFrom.UTC().UnixNano()) / int64(time.Millisecond),
					"$lte": int64(timeTo.UTC().UnixNano()) / int64(time.Millisecond),
				},
			}
		}

		if qpHaveStocked != "" {
			l := log.WithFields(logrus.Fields{"name": qpNameHaveStocked, "value": qpHaveStocked})
			l.Trace("Query parameter handling")
			b, err := strconv.ParseBool(qpHaveStocked)
			if err != nil {
				// Invalid query parameter value provided
				l.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Error("Failed to parse haveStocked")
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
		filter = bson.M{"$and": []bson.M{
			filterExpires,
			filterHaveStocked,
			filterName,
			filterType,
		}}
	}
	log.WithFields(logrus.Fields{"value": filter}).Debug("Filter data")

	// Attempt to get the documents
	documents, err := configuration.Mongo.FindDocuments(ctx, collection, filter, nil)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to get documents")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"quantity": len(documents), "value": documents}).Debug("Documents found")
	}

	// Check if document is a recipe
	if collection == config.MongoCollectionRecipes {
		log.Trace("Begin recipe scan")
		for _, document := range documents {
			// Check if recipe can be made (i.e. associated ingredients are stocked and not expiring)
			id := document["_id"]
			l := log.WithFields(logrus.Fields{"recipe": id})
			originalCanMake := document["isCookable"].(bool)
			isCookable, err := isCookable(ctx, &document)
			if err != nil {
				// Something broke
				l.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to determine cookable")
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte(err.Error()))
				return
			} else if isCookable != originalCanMake {
				// Update if different
				l.WithFields(logrus.Fields{"original": originalCanMake, "updated": isCookable}).Debug("Updating isCookable")
				document["isCookable"] = isCookable

				// Create filter
				filter := bson.D{{"_id", id}}
				l = l.WithFields(logrus.Fields{"method": "PUT"})
				l.WithFields(logrus.Fields{"value": filter}).Debug("Filter data")

				// Define update instructions
				update := bson.M{"$set": bson.M{"isCookable": isCookable}}
				l.WithFields(logrus.Fields{"value": update}).Debug("Update instructions")

				// Attempt to put the document
				matched, _, err := configuration.Mongo.UpdateOneDocument(ctx, collection, filter, update)
				if err != nil {
					// Put failed
					l.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to update recipe")
					response.WriteHeader(http.StatusInternalServerError)
					response.Write([]byte(err.Error()))
					return
				} else {
					l.WithFields(logrus.Fields{"quantity": matched}).Info("Updated recipe")
				}
			}
		}
		log.Trace("End recipe scan")
	}

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

func postManyDocuments(response http.ResponseWriter, request *http.Request) {
	// Setup
	ctx := request.Context()
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.postManyDocuments",
		"method": "POST",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Extract route parameters
	vars := mux.Vars(request)
	collection := vars["collection"]
	log.WithFields(logrus.Fields{"value": vars}).Debug("Route variables")
	log = log.WithFields(logrus.Fields{"collection": collection})

	// Validate collection route variable
	collections, err := configuration.Mongo.Collections(ctx)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed list collection names")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else if !utils.Contains(collections, collection) {
		err := fmt.Errorf("collection not found: %s", collection)
		log.WithFields(logrus.Fields{"collections": collections, "status": http.StatusNotFound}).WithError(err).Warn("Failed to find collection")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(err.Error()))
		return
	}

	// Get documents from body
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse request body")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"size": len(bytes), "state": "marshalled", "value": string(bytes)}).Debug("Request body")
	}

	// Parse documents
	var body []primitive.M
	err = json.Unmarshal(bytes, &body)
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
	} else {
		log.WithFields(logrus.Fields{"state": "unmarshalled", "value": body}).Debug("Request body")
	}

	log = log.WithFields(logrus.Fields{"quantity": len(body)})
	log.WithFields(logrus.Fields{"value": body}).Debug("Documents received")

	// Check if document is a recipe & prepare for insertion
	documents := []interface{}{}
	if collection == config.MongoCollectionRecipes {
		log.Trace("Begin recipe scan")
		for _, document := range body {
			// Check if recipe can be made (i.e. associated ingredients are stocked and not expiring)
			isCookable, err := isCookable(ctx, &document)
			l := log.WithFields(logrus.Fields{"recipe": document["_id"]})
			if err != nil {
				// Something broke
				l.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to determine cookable")
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte(err.Error()))
				return
			} else {
				l.WithFields(logrus.Fields{"isCookable": isCookable}).Debug("Updating isCookable")
				document["isCookable"] = isCookable
			}

			documents = append(documents, document)
		}
		log.Trace("End recipe scan")
	} else {
		for _, document := range body {
			documents = append(documents, document)
		}
	}

	// Attempt to put the document
	err = configuration.Mongo.InsertManyDocuments(ctx, collection, documents)
	if err != nil {
		// Post failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to post documents")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"status": http.StatusCreated}).Info("Succeeded")
		response.WriteHeader(http.StatusCreated)
	}
}

func putOneDocument(response http.ResponseWriter, request *http.Request) {
	// Setup
	ctx := request.Context()
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.putOneDocument",
		"method": "PUT",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Extract route parameter
	vars := mux.Vars(request)
	collection := vars["collection"]
	id := vars["id"]
	log.WithFields(logrus.Fields{"value": vars}).Debug("Route variables")
	log = log.WithFields(logrus.Fields{"collection": collection})

	// Validate collection route variable
	collections, err := configuration.Mongo.Collections(ctx)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed list collection names")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else if !utils.Contains(collections, collection) {
		err := fmt.Errorf("collection not found: %s", collection)
		log.WithFields(logrus.Fields{"collections": collections, "status": http.StatusNotFound}).WithError(err).Warn("Failed to find collection")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(err.Error()))
		return
	}

	// Parse document id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil && err.Error() == utils.ErrorInvalidObjectID {
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

	// Construct filter
	filter := bson.D{{"_id", oid}}
	log.WithFields(logrus.Fields{"value": filter}).Debug("Filter data")
	log = log.WithFields(logrus.Fields{"id": id})

	// Get document fields from body
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse request body")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"size": len(bytes), "state": "marshalled", "value": string(bytes)}).Debug("Request body")
	}

	// Parse update fields
	var fields map[string]interface{}
	err = json.Unmarshal(bytes, &fields)
	if err != nil && strings.HasPrefix(err.Error(), "invalid character") {
		// Invalid request body
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Warn("Failed to decode update fields")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to decode update fields")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"value": fields}).Trace("Update fields")
	}

	// Construct update instructions
	interim := bson.M{}
	for k, v := range fields {
		interim[k] = v
	}
	log.WithFields(logrus.Fields{"step": 1, "value": interim}).Trace("Interim update instructions")

	// Ignore _id field since it's immutable and will error
	_, found := interim["_id"]
	if found {
		delete(interim, "_id")
	}
	log.WithFields(logrus.Fields{"step": 2, "value": interim}).Trace("Interim update instructions")

	update := bson.M{"$set": interim}
	log.WithFields(logrus.Fields{"value": update}).Debug("Update instructions")

	// Attempt to put the document
	matched, _, err := configuration.Mongo.UpdateOneDocument(ctx, collection, filter, update)
	if err != nil && strings.Contains(err.Error(), "You must specify a field like so") {
		// Empty request body
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Warn("Failed to put document")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
	} else if err != nil {
		// Put failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to put document")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else if matched == 0 {
		// Put completed but no document was found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).Warn("Failed to put document")
		response.WriteHeader(http.StatusNotFound)
	} else {
		log.WithFields(logrus.Fields{"quantity": matched, "status": http.StatusOK}).Info("Succeeded")
		response.WriteHeader(http.StatusOK)
	}
}

func deleteOneDocument(response http.ResponseWriter, request *http.Request) {
	// Setup
	ctx := request.Context()
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.deleteOneDocument",
		"method": "DELETE",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Extract route parameter
	vars := mux.Vars(request)
	collection := vars["collection"]
	id := vars["id"]
	log.WithFields(logrus.Fields{"value": vars}).Debug("Route variables")
	log = log.WithFields(logrus.Fields{"collection": collection})

	// Validate collection route variable
	collections, err := configuration.Mongo.Collections(ctx)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed list collection names")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else if !utils.Contains(collections, collection) {
		err := fmt.Errorf("collection not found: %s", collection)
		log.WithFields(logrus.Fields{"collections": collections, "status": http.StatusNotFound}).WithError(err).Warn("Failed to find collection")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(err.Error()))
		return
	}

	// Parse document id
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil && err.Error() == utils.ErrorInvalidObjectID {
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
	log.WithFields(logrus.Fields{"value": filter}).Debug("Filter data")
	log = log.WithFields(logrus.Fields{"id": id})

	// Attempt to delete the document
	err = configuration.Mongo.DeleteOneDocument(ctx, collection, filter)
	if err != nil && err.Error() == utils.ErrorMongoNoDocuments {
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
		log.WithFields(logrus.Fields{"status": http.StatusOK}).Info("Succeeded")
		response.WriteHeader(http.StatusOK)
	}
}

func deleteManyDocuments(response http.ResponseWriter, request *http.Request) {
	// Setup
	ctx := request.Context()
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.deleteManyDocuments",
		"method": "DELETE",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Extract route parameters
	vars := mux.Vars(request)
	collection := vars["collection"]
	log.WithFields(logrus.Fields{"value": vars}).Debug("Route variables")
	log = log.WithFields(logrus.Fields{"collection": collection})

	// Validate collection route variable
	collections, err := configuration.Mongo.Collections(ctx)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed list collection names")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else if !utils.Contains(collections, collection) {
		err := fmt.Errorf("collection not found: %s", collection)
		log.WithFields(logrus.Fields{"collections": collections, "status": http.StatusNotFound}).WithError(err).Warn("Failed to find collection")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(err.Error()))
		return
	}

	// Delete by list of IDs (for now)
	// Get document fields from body
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to parse request body")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"size": len(bytes), "state": "marshalled", "value": string(bytes)}).Debug("Request body")
	}

	// Parse update fields
	var ids []string
	err = json.Unmarshal(bytes, &ids)
	if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Warn("Failed to decode delete fields")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"value": ids}).Debug("Document IDs")
	}

	// Create filter
	interim := []primitive.ObjectID{}
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.WithFields(logrus.Fields{"id": id, "status": http.StatusBadRequest}).WithError(err).Warn("Failed to parse id")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(err.Error()))
			return
		} else {
			interim = append(interim, oid)
		}
	}
	log.WithFields(logrus.Fields{"step": 1, "value": interim}).Trace("Interim update instructions")
	filter := bson.M{"_id": bson.M{"$in": interim}}
	log.WithFields(logrus.Fields{"value": filter}).Debug("Filter data")

	// Attempt to delete the documents
	deleted, err := configuration.Mongo.DeleteManyDocuments(ctx, collection, filter)
	if err != nil {
		// Delete failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to delete documents")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else if deleted == 0 {
		// Delete completed but no documents were found
		log.WithFields(logrus.Fields{"status": http.StatusNotFound}).WithError(err).Warn("Failed to delete documents")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte("no documents found"))
	} else {
		log.WithFields(logrus.Fields{"quantity": deleted, "status": http.StatusOK}).Info("Succeeded")
		response.WriteHeader(http.StatusOK)
	}
}

func checkExpirations() {
	// Setup
	ctx := context.Background()
	log := logrus.WithFields(logrus.Fields{
		"at":        "api.checkExpirations",
		"lookahead": configuration.Lookahead,
	})

	// Log diagnostic information
	log.Trace("Begin function")
	defer log.Trace("End function")

	// Filter by food expired already
	now := int64(time.Now().UTC().UnixNano()) / int64(time.Millisecond)
	later := int64(time.Now().Add(configuration.Lookahead).UTC().UnixNano()) / int64(time.Millisecond)
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
	log.WithFields(logrus.Fields{"type": "expired", "value": filterExpired}).Debug("Filter data")

	// Filter by food expiring within the given search window
	filterExpiring := bson.M{"$and": []bson.M{
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
	log.WithFields(logrus.Fields{"type": "expiring", "value": filterExpiring}).Debug("Filter data")

	// Grab the documents
	documentsExpired, err := configuration.Mongo.FindDocuments(ctx, config.MongoCollectionIngredients, filterExpired, nil)
	if err != nil {
		log.WithError(err).Error("Failed to identify expired items")
		return
	}

	documentsExpiring, err := configuration.Mongo.FindDocuments(ctx, config.MongoCollectionIngredients, filterExpiring, nil)
	if err != nil {
		log.WithError(err).Error("Failed to identify expiring items")
	} else {
		quantityExpired := len(documentsExpired)
		quantityExpiring := len(documentsExpiring)

		// Skip if nothing is expiring
		if quantityExpiring == 0 && quantityExpired == 0 {
			log.WithFields(logrus.Fields{"expiring": quantityExpiring, "expired": quantityExpired}).Info("Restocking not required")
			return
		} else {
			log.WithFields(logrus.Fields{"expiring": quantityExpiring, "expired": quantityExpired}).Info("Restocking required")
		}

		log.WithFields(logrus.Fields{"quantity": quantityExpired, "value": documentsExpired}).Debug("Expired items")
		log.WithFields(logrus.Fields{"quantity": quantityExpiring, "value": documentsExpiring}).Debug("Expiring items")

		// Construct list of names of items to shop for
		var groceries []string
		for _, document := range documentsExpired {
			name := document["name"]
			stage := "expired"
			text := fmt.Sprintf("%s (%s)", name, stage)

			if _, ok := document["attributes"]; ok {
				brand := document["attributes"].(map[string]string)["brand"]
				flavor := document["attributes"].(map[string]string)["flavor"]
				text = fmt.Sprintf("%s (%s, %s, %s)", name, brand, flavor, stage)
			}

			groceries = append(groceries, text)
		}
		for _, document := range documentsExpiring {
			name := document["name"]
			stage := "expiring"
			text := fmt.Sprintf("%s (%s)", name, stage)

			if _, ok := document["attributes"]; ok {
				brand := document["attributes"].(map[string]string)["brand"]
				flavor := document["attributes"].(map[string]string)["flavor"]
				text = fmt.Sprintf("%s (%s, %s, %s)", name, brand, flavor, stage)
			}

			groceries = append(groceries, text)
		}
		log.WithFields(logrus.Fields{"quantity": len(groceries), "value": groceries}).Debug("Groceries")

		// Construct shopping list due date
		now := time.Now()
		rounded := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		dueDate := rounded.Add(configuration.Lookahead + (time.Hour * 24))
		log.WithFields(logrus.Fields{"value": dueDate}).Debug("Card due date")

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
			innerTrello := reflect.ValueOf(configuration.Trello).Elem()
			labelsStr := *innerTrello.FieldByName("LabelsStr").Addr().Interface().(*string)
			labels := strings.Split(labelsStr, ",")
			url, err = configuration.Trello.CreateShoppingList(&dueDate, labels, groceries)
			if err != nil {
				log.WithError(err).Error("Failed to create Trello card")
			} else {
				log.WithFields(logrus.Fields{"url": url}).Info("Created Trello card")
			}
		}

		// Compose Twilio message
		var message = configuration.Twilio.ComposeMessage(quantityExpiring, quantityExpired, url)

		// Send the Twilio message
		if !configuration.Silence {
			innerTwilio := reflect.ValueOf(configuration.Twilio).Elem()
			from := *innerTwilio.FieldByName("From").Addr().Interface().(*string)
			to := *innerTwilio.FieldByName("To").Addr().Interface().(*string)
			_, err = configuration.Twilio.SendMessage(from, to, message)
			if err != nil {
				log.WithFields(logrus.Fields{"from": from, "to": to}).WithError(err).Error("Failed to send Twilio message")
			} else {
				log.WithFields(logrus.Fields{"from": from, "to": to}).Info("Sent Twilio message")
			}
		} else {
			log.WithFields(logrus.Fields{"silence": configuration.Silence}).Info("Skipped Twilio message")
		}
	}
}
