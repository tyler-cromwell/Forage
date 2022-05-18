package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/clients"
	"github.com/tyler-cromwell/forage/config"
	"github.com/tyler-cromwell/forage/tests/mocks/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAPI(t *testing.T) {
	forageContextTimeout, _ := time.ParseDuration("5s")
	forageInterval, _ := time.ParseDuration("24h")
	forageLookahead, _ := time.ParseDuration("48h")

	// Initialize context/timeout
	ctx, cancel := context.WithTimeout(context.Background(), forageContextTimeout)
	logrus.WithFields(logrus.Fields{"timeout": forageContextTimeout}).Info("Initialized context")
	defer cancel()

	// Initialize client
	mongoUri := "mongodb://127.0.0.1:27017"
	listenSocket := ":8001"
	mongoClient, err := clients.NewMongoClientWrapper(ctx, mongoUri)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create MongoDB client wrapper")
	} else {
		defer mongoClient.Client.Disconnect(ctx)
	}

	configuration = &config.Configuration{
		ContextTimeout: forageContextTimeout,
		Interval:       forageInterval,
		Lookahead:      forageLookahead,
		LogrusLevel:    logrus.DebugLevel,
		ListenSocket:   listenSocket,
	}

	subtests := []struct {
		name            string
		requestType     string
		endpoint        string
		handler         func(http.ResponseWriter, *http.Request)
		status          int
		body            string
		urlVars         map[string]string
		queryParameters map[string]string
		mongoClient     mongo.MockMongo
	}{
		{"getExpired200", "GET", "/expired", getExpired, http.StatusOK, "[]", nil, nil, mongo.MockMongo{}},
		{"getExpired500", "GET", "/expired", getExpired, http.StatusInternalServerError, "failure", nil, nil, mongo.MockMongo{
			// Mongo failure
			OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return nil, fmt.Errorf("failure")
			},
		}},
		/*
			{"getExpired500#2", "GET", "/expired", getExpired, http.StatusInternalServerError, `""`, mongo.MockMongo{
				// JSON failure (not sure this case is possible)
				OverrideFindManyDocuments: func(context.Context, bson.M, *options.FindOptions) ([]bson.M, error) {
					docs := []bson.M{}
					return docs, nil
				},
			}},
		*/
		{"getExpiring200", "GET", "/expiring", getExpiring, http.StatusOK, "[]", nil, map[string]string{"from": "10", "to": "20"}, mongo.MockMongo{}},
		{"getExpiring400#1", "GET", "/expiring", getExpiring, http.StatusBadRequest, "strconv.ParseInt: parsing \"x\": invalid syntax", nil, map[string]string{"from": "x", "to": ""}, mongo.MockMongo{}},
		{"getExpiring400#2", "GET", "/expiring", getExpiring, http.StatusBadRequest, "strconv.ParseInt: parsing \"y\": invalid syntax", nil, map[string]string{"from": "10", "to": "y"}, mongo.MockMongo{}},
		{"getExpiring500", "GET", "/expiring", getExpiring, http.StatusInternalServerError, "failure", nil, nil, mongo.MockMongo{
			// Mongo failure
			OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return nil, fmt.Errorf("failure")
			},
		}},
		{"getOneDocument", "GET", "/documents", getOneDocument, http.StatusOK, "null", map[string]string{"id": "6187e576abc057dac3e7d5dc"}, nil, mongo.MockMongo{}},
		{"getManyDocuments200", "GET", "/documents", getManyDocuments, http.StatusOK, "[]", nil, nil, mongo.MockMongo{}},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			configuration.Mongo = &st.mongoClient

			req, err := http.NewRequest(st.requestType, st.endpoint, nil)
			if err != nil {
				t.Fatal(err)
			}

			if st.queryParameters != nil {
				q := req.URL.Query()
				for k, v := range st.queryParameters {
					q.Add(k, v)
				}
				req.URL.RawQuery = q.Encode()
			}

			rr := httptest.NewRecorder()
			if st.urlVars != nil {
				req = mux.SetURLVars(req, st.urlVars)
			}
			handler := http.HandlerFunc(st.handler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != st.status {
				t.Errorf("handler returned wrong status code: got %v want %v", status, st.status)
			}

			if rr.Body.String() != st.body {
				t.Errorf("handler returned unexpected body: got \"%v\" want \"%v\"", rr.Body.String(), st.body)
			}
		})
	}
}
