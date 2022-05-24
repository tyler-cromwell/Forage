package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/clients"
	"github.com/tyler-cromwell/forage/config"
	"github.com/tyler-cromwell/forage/tests/mocks/mongo"
	"github.com/tyler-cromwell/forage/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testRequest struct {
	method          string
	endpoint        string
	routeVariables  map[string]string
	queryParameters map[string]string
	body            string
}

type testResponse struct {
	status int
	body   string
}

func TestAPI(t *testing.T) {
	// Discard logging output
	logrus.SetOutput(ioutil.Discard)

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
		name        string
		handler     func(http.ResponseWriter, *http.Request)
		request     testRequest
		response    testResponse
		mongoClient mongo.MockMongo
	}{
		{"getExpired200", getExpired, testRequest{method: "GET", endpoint: "/expired", routeVariables: nil, queryParameters: nil, body: ""}, testResponse{status: http.StatusOK, body: "[]"}, mongo.MockMongo{}},
		{"getExpired500", getExpired, testRequest{method: "GET", endpoint: "/expired", routeVariables: nil, queryParameters: nil, body: ""}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mongo.MockMongo{
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
		{"getExpiring200", getExpiring, testRequest{method: "GET", endpoint: "/expiring", routeVariables: nil, queryParameters: map[string]string{"from": "10", "to": "20"}, body: ""}, testResponse{status: http.StatusOK, body: "[]"}, mongo.MockMongo{}},
		{"getExpiring400#1", getExpiring, testRequest{method: "GET", endpoint: "/expiring", routeVariables: nil, queryParameters: map[string]string{"from": "x", "to": "20"}, body: ""}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseInt: parsing \"x\": invalid syntax"}, mongo.MockMongo{}},
		{"getExpiring400#2", getExpiring, testRequest{method: "GET", endpoint: "/expiring", routeVariables: nil, queryParameters: map[string]string{"from": "10", "to": "y"}, body: ""}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseInt: parsing \"y\": invalid syntax"}, mongo.MockMongo{}},
		{"getExpiring500", getExpiring, testRequest{method: "GET", endpoint: "/expiring", routeVariables: nil, queryParameters: nil, body: ""}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mongo.MockMongo{
			OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return nil, fmt.Errorf("failure")
			},
		}},
		{"getOneDocument200", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: ""}, testResponse{status: http.StatusOK, body: "null"}, mongo.MockMongo{}},
		{"getOneDocument400", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"id": "hello"}, queryParameters: nil, body: ""}, testResponse{status: http.StatusBadRequest, body: "the provided hex string is not a valid ObjectID"}, mongo.MockMongo{}},
		{"getOneDocument404", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: ""}, testResponse{status: http.StatusNotFound, body: utils.ErrMongoNoDocuments}, mongo.MockMongo{
			OverrideFindOneDocument: func(ctx context.Context, filter bson.D) (*bson.M, error) {
				return nil, fmt.Errorf(utils.ErrMongoNoDocuments)
			},
		}},
		{"getOneDocument500", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: ""}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mongo.MockMongo{
			OverrideFindOneDocument: func(ctx context.Context, filter bson.D) (*bson.M, error) {
				return nil, fmt.Errorf("failure")
			},
		}},
		{"getManyDocuments200", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: nil, queryParameters: map[string]string{"name": "hello", "type": "thing", "haveStocked": "false", "from": "10", "to": "20"}, body: ""}, testResponse{status: http.StatusOK, body: "[]"}, mongo.MockMongo{}},
		{"getManyDocuments400#1", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: nil, queryParameters: map[string]string{"name": "hello", "type": "thing", "haveStocked": "false", "from": "x", "to": ""}, body: ""}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseInt: parsing \"x\": invalid syntax"}, mongo.MockMongo{}},
		{"getManyDocuments400#2", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: nil, queryParameters: map[string]string{"name": "hello", "type": "thing", "haveStocked": "false", "from": "10", "to": "y"}, body: ""}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseInt: parsing \"y\": invalid syntax"}, mongo.MockMongo{}},
		{"getManyDocuments400#3", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: nil, queryParameters: map[string]string{"name": "hello", "type": "thing", "haveStocked": "lol", "from": "10", "to": "20"}, body: ""}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseBool: parsing \"lol\": invalid syntax"}, mongo.MockMongo{}},
		{"getManyDocuments500", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: nil, queryParameters: map[string]string{"name": "hello", "type": "thing", "haveStocked": "false", "from": "10", "to": "20"}, body: ""}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mongo.MockMongo{
			OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return nil, fmt.Errorf("failure")
			},
		}},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			configuration.Mongo = &st.mongoClient

			req, err := http.NewRequest(st.request.method, st.request.endpoint, nil)
			if err != nil {
				t.Fatal(err)
			}

			if st.request.queryParameters != nil {
				q := req.URL.Query()
				for k, v := range st.request.queryParameters {
					q.Add(k, v)
				}
				req.URL.RawQuery = q.Encode()
			}

			rr := httptest.NewRecorder()
			if st.request.routeVariables != nil {
				req = mux.SetURLVars(req, st.request.routeVariables)
			}
			handler := http.HandlerFunc(st.handler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != st.response.status {
				t.Errorf("handler returned wrong status code: got %v want %v", status, st.response.status)
			}

			if rr.Body.String() != st.response.body {
				t.Errorf("handler returned unexpected body: got \"%v\" want \"%v\"", rr.Body.String(), st.response.body)
			}
		})
	}
}
