package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
