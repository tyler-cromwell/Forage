package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/adlio/trello"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
	"github.com/tyler-cromwell/forage/config"
	"github.com/tyler-cromwell/forage/tests/mocks"
	"github.com/tyler-cromwell/forage/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testRequest struct {
	method          string
	endpoint        string
	routeVariables  map[string]string
	queryParameters map[string]string
	body            io.ReadCloser
}

type testResponse struct {
	status int
	body   string
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestAPI(t *testing.T) {
	// Discard logging output
	logrus.SetOutput(ioutil.Discard)

	forageContextTimeout, _ := time.ParseDuration("5s")
	forageLookahead, _ := time.ParseDuration("48h")
	listenSocket := ":8001"

	configuration = &config.Configuration{
		ContextTimeout: forageContextTimeout,
		Lookahead:      forageLookahead,
		LogrusLevel:    logrus.DebugLevel,
		ListenSocket:   listenSocket,
	}

	subtests := []struct {
		name        string
		handler     func(http.ResponseWriter, *http.Request)
		request     testRequest
		response    testResponse
		mongoClient mocks.MockMongo
	}{
		{"getConfiguration", getConfiguration, testRequest{method: "GET", endpoint: "/configure", routeVariables: nil, queryParameters: nil, body: nil}, testResponse{status: http.StatusOK, body: "{\"Lookahead\":172800000000000,\"Time\":\"\"}"}, mocks.MockMongo{}},
		{"putConfiguration200", putConfiguration, testRequest{method: "PUT", endpoint: "/configure", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("{\"lookahead\": 172800000000000, \"time\": \"24h\"}"))}, testResponse{status: http.StatusOK, body: ""}, mocks.MockMongo{}},
		{"putConfiguration400", putConfiguration, testRequest{method: "PUT", endpoint: "/configure", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("{:}"))}, testResponse{status: http.StatusBadRequest, body: "invalid character ':' looking for beginning of object key string"}, mocks.MockMongo{}},
		{"putConfiguration500#1", putConfiguration, testRequest{method: "PUT", endpoint: "/configure", routeVariables: nil, queryParameters: nil, body: io.NopCloser(errReader(0))}, testResponse{status: http.StatusInternalServerError, body: "test error"}, mocks.MockMongo{}},
		{"putConfiguration500#2", putConfiguration, testRequest{method: "PUT", endpoint: "/configure", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("{\"lookahead\": \"172800000000000\", \"time\": \"24h\"}"))}, testResponse{status: http.StatusInternalServerError, body: "json: cannot unmarshal string into Go struct field .lookahead of type time.Duration"}, mocks.MockMongo{}},
		{"getExpired200", getExpired, testRequest{method: "GET", endpoint: "/expired", routeVariables: nil, queryParameters: nil, body: nil}, testResponse{status: http.StatusOK, body: "[]"}, mocks.MockMongo{}},
		{"getExpired500#1", getExpired, testRequest{method: "GET", endpoint: "/expired", routeVariables: nil, queryParameters: nil, body: nil}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mocks.MockMongo{
			OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return nil, fmt.Errorf("failure")
			},
		}},
		{"getExpired500#2", getExpired, testRequest{method: "GET", endpoint: "/expired", routeVariables: nil, queryParameters: nil, body: nil}, testResponse{status: http.StatusInternalServerError, body: "json: unsupported type: chan int"}, mocks.MockMongo{
			OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return []bson.M{map[string]interface{}{"key": make(chan int)}}, nil
			},
		}},
		{"getExpiring200", getExpiring, testRequest{method: "GET", endpoint: "/expiring", routeVariables: nil, queryParameters: map[string]string{"from": "10", "to": "20"}, body: nil}, testResponse{status: http.StatusOK, body: "[]"}, mocks.MockMongo{}},
		{"getExpiring400#1", getExpiring, testRequest{method: "GET", endpoint: "/expiring", routeVariables: nil, queryParameters: map[string]string{"from": "x", "to": "20"}, body: nil}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseInt: parsing \"x\": invalid syntax"}, mocks.MockMongo{}},
		{"getExpiring400#2", getExpiring, testRequest{method: "GET", endpoint: "/expiring", routeVariables: nil, queryParameters: map[string]string{"from": "10", "to": "y"}, body: nil}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseInt: parsing \"y\": invalid syntax"}, mocks.MockMongo{}},
		{"getExpiring500#1", getExpiring, testRequest{method: "GET", endpoint: "/expiring", routeVariables: nil, queryParameters: nil, body: nil}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mocks.MockMongo{
			OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return nil, fmt.Errorf("failure")
			},
		}},
		{"getExpiring500#2", getExpiring, testRequest{method: "GET", endpoint: "/expiring", routeVariables: nil, queryParameters: nil, body: nil}, testResponse{status: http.StatusInternalServerError, body: "json: unsupported type: chan int"}, mocks.MockMongo{
			OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return []bson.M{map[string]interface{}{"key": make(chan int)}}, nil
			},
		}},
		{"getOneDocument200", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusOK, body: "null"}, mocks.MockMongo{}},
		{"getOneDocument400", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"id": "hello"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusBadRequest, body: "the provided hex string is not a valid ObjectID"}, mocks.MockMongo{}},
		{"getOneDocument404", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusNotFound, body: utils.ErrMongoNoDocuments}, mocks.MockMongo{
			OverrideFindOneDocument: func(ctx context.Context, filter bson.D) (*bson.M, error) {
				return nil, fmt.Errorf(utils.ErrMongoNoDocuments)
			},
		}},
		{"getOneDocument500#1", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"id": "xxxxxxxxxxxxxxxxxxxxxxxx"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusInternalServerError, body: "encoding/hex: invalid byte: U+0078 'x'"}, mocks.MockMongo{}},
		{"getOneDocument500#2", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusInternalServerError, body: "json: unsupported type: chan int"}, mocks.MockMongo{
			OverrideFindOneDocument: func(ctx context.Context, filter bson.D) (*bson.M, error) {
				var doc bson.M = map[string]interface{}{"key": make(chan int)}
				return &doc, nil
			},
		}},
		{"getOneDocument500#3", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mocks.MockMongo{
			OverrideFindOneDocument: func(ctx context.Context, filter bson.D) (*bson.M, error) {
				return nil, fmt.Errorf("failure")
			},
		}},
		{"getManyDocuments200", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: nil, queryParameters: map[string]string{"name": "hello", "type": "thing", "haveStocked": "false", "from": "10", "to": "20"}, body: nil}, testResponse{status: http.StatusOK, body: "[]"}, mocks.MockMongo{}},
		{"getManyDocuments400#1", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: nil, queryParameters: map[string]string{"name": "hello", "type": "thing", "haveStocked": "false", "from": "x", "to": ""}, body: nil}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseInt: parsing \"x\": invalid syntax"}, mocks.MockMongo{}},
		{"getManyDocuments400#2", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: nil, queryParameters: map[string]string{"name": "hello", "type": "thing", "haveStocked": "false", "from": "10", "to": "y"}, body: nil}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseInt: parsing \"y\": invalid syntax"}, mocks.MockMongo{}},
		{"getManyDocuments400#3", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: nil, queryParameters: map[string]string{"name": "hello", "type": "thing", "haveStocked": "lol", "from": "10", "to": "20"}, body: nil}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseBool: parsing \"lol\": invalid syntax"}, mocks.MockMongo{}},
		{"getManyDocuments500#1", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: nil, queryParameters: map[string]string{"name": "hello", "type": "thing", "haveStocked": "false", "from": "10", "to": "20"}, body: nil}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mocks.MockMongo{
			OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return nil, fmt.Errorf("failure")
			},
		}},
		{"getManyDocuments500#2", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: nil, queryParameters: map[string]string{"name": "hello", "type": "thing", "haveStocked": "false", "from": "10", "to": "20"}, body: nil}, testResponse{status: http.StatusInternalServerError, body: "json: unsupported type: chan int"}, mocks.MockMongo{
			OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return []bson.M{map[string]interface{}{"key": make(chan int)}}, nil
			},
		}},
		{"postManyDocuments200", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("[{\"name\": \"Document\"}]"))}, testResponse{status: http.StatusCreated, body: ""}, mocks.MockMongo{}},
		{"postManyDocuments400", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("{:}"))}, testResponse{status: http.StatusBadRequest, body: "invalid character ':' looking for beginning of object key string"}, mocks.MockMongo{}},
		{"postManyDocuments500#1", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: nil, queryParameters: nil, body: io.NopCloser(errReader(0))}, testResponse{status: http.StatusInternalServerError, body: "test error"}, mocks.MockMongo{}},
		{"postManyDocuments500#2", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("{}"))}, testResponse{status: http.StatusInternalServerError, body: "json: cannot unmarshal object into Go value of type []interface {}"}, mocks.MockMongo{}},
		{"postManyDocuments500#3", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("[{\"name\": \"Document\"}]"))}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mocks.MockMongo{
			OverrideInsertManyDocuments: func(ctx context.Context, docs []interface{}) error {
				return fmt.Errorf("failure")
			},
		}},
		{"putOneDocument200", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: io.NopCloser(strings.NewReader("{\"_id\": \"6187e576abc057dac3e7d5dc\", \"name\": \"Document\"}"))}, testResponse{status: http.StatusOK, body: ""}, mocks.MockMongo{}},
		{"putOneDocument400#1", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"id": "hello"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusBadRequest, body: "the provided hex string is not a valid ObjectID"}, mocks.MockMongo{}},
		{"putOneDocument400#2", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: io.NopCloser(strings.NewReader("{:}"))}, testResponse{status: http.StatusBadRequest, body: "invalid character ':' looking for beginning of object key string"}, mocks.MockMongo{}},
		{"putOneDocument404", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: io.NopCloser(strings.NewReader("{\"_id\": \"6187e576abc057dac3e7d5dc\", \"name\": \"Document\"}"))}, testResponse{status: http.StatusNotFound, body: "failure"}, mocks.MockMongo{
			OverrideUpdateOneDocument: func(ctx context.Context, filter bson.D, update interface{}) (int64, int64, error) {
				return 0, 0, fmt.Errorf("failure")
			},
		}},
		{"putOneDocument500#1", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"id": "xxxxxxxxxxxxxxxxxxxxxxxx"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusInternalServerError, body: "encoding/hex: invalid byte: U+0078 'x'"}, mocks.MockMongo{}},
		{"putOneDocument500#2", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: io.NopCloser(errReader(0))}, testResponse{status: http.StatusInternalServerError, body: "test error"}, mocks.MockMongo{}},
		{"putOneDocument500#3", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: io.NopCloser(strings.NewReader("[{}"))}, testResponse{status: http.StatusInternalServerError, body: "unexpected end of JSON input"}, mocks.MockMongo{}},
		{"putOneDocument500#4", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: io.NopCloser(strings.NewReader("{\"_id\": \"6187e576abc057dac3e7d5dc\", \"name\": \"Document\"}"))}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mocks.MockMongo{
			OverrideUpdateOneDocument: func(ctx context.Context, filter bson.D, update interface{}) (int64, int64, error) {
				return 1, 1, fmt.Errorf("failure")
			},
		}},
		{"deleteOneDocument200", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusOK, body: ""}, mocks.MockMongo{}},
		{"deleteOneDocument400", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"id": "hello"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusBadRequest, body: "the provided hex string is not a valid ObjectID"}, mocks.MockMongo{}},
		{"deleteOneDocument404", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusNotFound, body: utils.ErrMongoNoDocuments}, mocks.MockMongo{
			OverrideDeleteOneDocument: func(ctx context.Context, filter bson.D) error {
				return fmt.Errorf(utils.ErrMongoNoDocuments)
			},
		}},
		{"deleteOneDocument500#1", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"id": "xxxxxxxxxxxxxxxxxxxxxxxx"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusInternalServerError, body: "encoding/hex: invalid byte: U+0078 'x'"}, mocks.MockMongo{}},
		{"deleteOneDocument500#2", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: nil}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mocks.MockMongo{
			OverrideDeleteOneDocument: func(ctx context.Context, filter bson.D) error {
				return fmt.Errorf("failure")
			},
		}},
		{"deleteManyDocuments200", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("[\"6187e576abc057dac3e7d5dc\"]"))}, testResponse{status: http.StatusOK, body: ""}, mocks.MockMongo{}},
		{"deleteManyDocuments400#1", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("{:}"))}, testResponse{status: http.StatusBadRequest, body: "invalid character ':' looking for beginning of object key string"}, mocks.MockMongo{}},
		{"deleteManyDocuments400#2", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("[\"hello\"]"))}, testResponse{status: http.StatusBadRequest, body: "the provided hex string is not a valid ObjectID"}, mocks.MockMongo{}},
		{"deleteManyDocuments404", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("[\"6187e576abc057dac3e7d5dc\"]"))}, testResponse{status: http.StatusNotFound, body: "no documents found"}, mocks.MockMongo{
			OverrideDeleteManyDocuments: func(ctx context.Context, filter bson.M) (int64, error) {
				return 0, nil
			},
		}},
		{"deleteManyDocuments500#1", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"id": "6187e576abc057dac3e7d5dc"}, queryParameters: nil, body: io.NopCloser(errReader(0))}, testResponse{status: http.StatusInternalServerError, body: "test error"}, mocks.MockMongo{}},
		{"deleteManyDocuments500#2", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: nil, queryParameters: nil, body: io.NopCloser(strings.NewReader("[\"6187e576abc057dac3e7d5dc\"]"))}, testResponse{status: http.StatusInternalServerError, body: "failure"}, mocks.MockMongo{
			OverrideDeleteManyDocuments: func(ctx context.Context, filter bson.M) (int64, error) {
				return 0, fmt.Errorf("failure")
			},
		}},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			configuration.Mongo = &st.mongoClient

			req, err := http.NewRequest(st.request.method, st.request.endpoint, st.request.body)
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

			if st.request.routeVariables != nil {
				req = mux.SetURLVars(req, st.request.routeVariables)
			}

			rr := httptest.NewRecorder()
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

	subtests2 := []struct {
		name         string
		mongoClient  mocks.MockMongo
		mocksClient  mocks.MockTrello
		twilioClient mocks.MockTwilio
		logLevels    []logrus.Level
		logMessages  []string
	}{
		{
			// Error #1, Could not obtain expired items
			"checkExpirationsError#1",
			mocks.MockMongo{
				OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
					expectation := bson.M{"$and": []bson.M{
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
					e, err := bson.Marshal(expectation)
					if err != nil {
						return nil, err
					}
					f, err := bson.Marshal(filter)
					if err != nil {
						return nil, err
					}

					if bytes.Equal(f, e) {
						return nil, fmt.Errorf("failure")
					} else {
						return []bson.M{
							map[string]interface{}{"key1": "value1"},
							map[string]interface{}{"key2": "value2"},
						}, nil
					}
				},
			},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{
				logrus.ErrorLevel,
			},
			[]string{
				"Failed to identify expired items",
			},
		},
		{
			// Error #2, Could not obtain expiring items
			"checkExpirationsError#2",
			mocks.MockMongo{
				OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
					expectation := bson.M{"$and": []bson.M{
						{
							"expirationDate": bson.M{
								"$gte": time.Now().UnixNano() / int64(time.Millisecond),
								"$lte": time.Now().Add(configuration.Lookahead).UnixNano() / int64(time.Millisecond),
							},
						},
						{
							"haveStocked": bson.M{
								"$eq": true,
							},
						},
					}}
					e, err := bson.Marshal(expectation)
					if err != nil {
						return nil, err
					}
					f, err := bson.Marshal(filter)
					if err != nil {
						return nil, err
					}

					if bytes.Equal(f, e) {
						return nil, fmt.Errorf("failure")
					} else {
						return []bson.M{
							map[string]interface{}{"key1": "value1"},
							map[string]interface{}{"key2": "value2"},
						}, nil
					}
				},
			},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{
				logrus.ErrorLevel,
			},
			[]string{
				"Failed to identify expiring items",
			},
		},
		{
			// Success #1, No expired/expiring items, no need to proceed.
			"checkExpirationsSuccess#1",
			mocks.MockMongo{},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{
				logrus.InfoLevel,
			},
			[]string{
				"Restocking not required",
			},
		},
		{
			// Success #2, items expired/expiring added to existing Trello card and SMS message sent.
			"checkExpirationsSuccess#2",
			mocks.MockMongo{
				OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
					return []bson.M{
						map[string]interface{}{"key1": "value1"},
						map[string]interface{}{"key2": "value2"},
					}, nil
				},
			},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{
				logrus.InfoLevel,
				logrus.InfoLevel,
				logrus.InfoLevel,
			},
			[]string{
				"Restocking required",
				"Added to Trello card",
				"Sent Twilio message",
			},
		},
		{
			// Error #3, items expired/expiring but could not obtain Trello card, SMS message still sent.
			"checkExpirationsError#3",
			mocks.MockMongo{
				OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
					return []bson.M{
						map[string]interface{}{"key1": "value1"},
						map[string]interface{}{"key2": "value2"},
					}, nil
				},
			},
			mocks.MockTrello{
				OverrideGetShoppingList: func() (*trello.Card, error) {
					return nil, fmt.Errorf("failure")
				},
			},
			mocks.MockTwilio{},
			[]logrus.Level{
				logrus.InfoLevel,
				logrus.ErrorLevel,
				logrus.InfoLevel,
			},
			[]string{
				"Restocking required",
				"Failed to get Trello card",
				"Sent Twilio message",
			},
		},
		{
			// Success #3, items expired/expiring added to new Trello card and SMS message sent.
			"checkExpirationsSuccess#3",
			mocks.MockMongo{
				OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
					return []bson.M{
						map[string]interface{}{"key1": "value1"},
						map[string]interface{}{"key2": "value2"},
					}, nil
				},
			},
			mocks.MockTrello{
				OverrideGetShoppingList: func() (*trello.Card, error) {
					return nil, nil
				},
			},
			mocks.MockTwilio{},
			[]logrus.Level{
				logrus.InfoLevel,
				logrus.InfoLevel,
				logrus.InfoLevel,
			},
			[]string{
				"Restocking required",
				"Created Trello card",
				"Sent Twilio message",
			},
		},
		{
			// Error #4, items expired/expiring but could not add to existing card Trello card, SMS message still sent.
			"checkExpirationsError#4",
			mocks.MockMongo{
				OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
					return []bson.M{
						map[string]interface{}{"key1": "value1"},
						map[string]interface{}{"key2": "value2"},
					}, nil
				},
			},
			mocks.MockTrello{
				OverrideAddToShoppingList: func(itemNames []string) (string, error) {
					return "", fmt.Errorf("failure")
				},
			},
			mocks.MockTwilio{},
			[]logrus.Level{
				logrus.InfoLevel,
				logrus.ErrorLevel,
				logrus.InfoLevel,
			},
			[]string{
				"Restocking required",
				"Failed to add to Trello card",
				"Sent Twilio message",
			},
		},
		{
			// Error #5, items expired/expiring but could not create new card Trello card, SMS message still sent.
			"checkExpirationsError#5",
			mocks.MockMongo{
				OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
					return []bson.M{
						map[string]interface{}{"key1": "value1"},
						map[string]interface{}{"key2": "value2"},
					}, nil
				},
			},
			mocks.MockTrello{
				OverrideGetShoppingList: func() (*trello.Card, error) {
					return nil, nil
				},
				OverrideCreateShoppingList: func(dueDate *time.Time, applyLabels []string, listItems []string) (string, error) {
					return "", fmt.Errorf("failure")
				},
			},
			mocks.MockTwilio{},
			[]logrus.Level{
				logrus.InfoLevel,
				logrus.ErrorLevel,
				logrus.InfoLevel,
			},
			[]string{
				"Restocking required",
				"Failed to create Trello card",
				"Sent Twilio message",
			},
		},
		{
			// Error #6, items expired/expiring but could not create new card Trello card or send SMS message.
			"checkExpirationsError#6",
			mocks.MockMongo{
				OverrideFindManyDocuments: func(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
					return []bson.M{
						map[string]interface{}{"key1": "value1"},
						map[string]interface{}{"key2": "value2"},
					}, nil
				},
			},
			mocks.MockTrello{
				OverrideGetShoppingList: func() (*trello.Card, error) {
					return nil, nil
				},
				OverrideCreateShoppingList: func(dueDate *time.Time, applyLabels []string, listItems []string) (string, error) {
					return "", fmt.Errorf("failure")
				},
			},
			mocks.MockTwilio{
				OverrideComposeMessage: func(quantity, quantityExpired int, url string) string {
					return ""
				},
				OverrideSendMessage: func(phoneFrom, phoneTo, message string) (string, error) {
					return "", fmt.Errorf("failure")
				},
			},
			[]logrus.Level{
				logrus.InfoLevel,
				logrus.ErrorLevel,
				logrus.ErrorLevel,
			},
			[]string{
				"Restocking required",
				"Failed to create Trello card",
				"Failed to send Twilio message",
			},
		},
	}
	t.Run("checkExpirations", func(t *testing.T) {
		// Capture logrus output so we can assert
		_, hook := test.NewNullLogger()
		logrus.AddHook(hook)
		base := 0

		for _, st := range subtests2 {
			// Arrange
			configuration.Mongo = &st.mongoClient
			configuration.Trello = &st.mocksClient
			configuration.Twilio = &st.twilioClient

			// Act
			checkExpirations()

			// Assert (preliminary)
			require.Equal(t, len(st.logLevels), len(st.logMessages))

			// Assert (primary)
			for i, _ := range st.logLevels {
				index := base + i
				require.Equal(t, st.logLevels[i], hook.AllEntries()[index].Level)
				require.Equal(t, st.logMessages[i], hook.AllEntries()[index].Message)
			}

			base += len(st.logLevels)
		}

		// Rever logrus output change
		logrus.SetOutput(ioutil.Discard)
	})
}
