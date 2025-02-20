package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
	"github.com/tyler-cromwell/forage/config"
	"github.com/tyler-cromwell/forage/tests/mocks"
	"github.com/tyler-cromwell/forage/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testRequest struct {
	m  string
	e  string
	rv map[string]string
	qp map[string]string
	b  io.ReadCloser
}

type testResponse struct {
	s int
	b string
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New(errorIoutilReadAll)
}

func TestAPI(t *testing.T) {
	// Discard logging output
	logrus.SetOutput(io.Discard)

	forageContextTimeout, _ := time.ParseDuration("5s")
	forageLookahead, _ := time.ParseDuration("48h")
	listenSocket := ":8001"
	loc, _ := time.LoadLocation("America/New_York")

	configuration = &config.Configuration{
		ContextTimeout: forageContextTimeout,
		Lookahead:      forageLookahead,
		LogrusLevel:    logrus.DebugLevel,
		ListenSocket:   listenSocket,
		Scheduler:      gocron.NewScheduler(loc),
	}

	subtests := []struct {
		name        string
		handler     func(http.ResponseWriter, *http.Request)
		request     testRequest
		response    testResponse
		mongoClient mocks.MockMongo
	}{
		{"getConfiguration200#1", getConfiguration, testRequest{m: "GET", e: "/configure"}, testResponse{s: http200, b: "{\"lookahead\":172800000000000,\"silence\":false,\"time\":\"\"}"}, mocks.MockMongo{}},
		{"putConfiguration200#1", putConfiguration, testRequest{m: "PUT", e: "/configure", b: io.NopCloser(strings.NewReader("{\"lookahead\": 172800000000000, \"time\": \"19:00\"}"))}, testResponse{s: http200}, mocks.MockMongo{}},
		{"putConfiguration400#1", putConfiguration, testRequest{m: "PUT", e: "/configure", b: io.NopCloser(strings.NewReader("{:}"))}, testResponse{s: http400, b: errorJsonUndecodable}, mocks.MockMongo{}},
		{"putConfiguration400#2", putConfiguration, testRequest{m: "PUT", e: "/configure", b: io.NopCloser(strings.NewReader(""))}, testResponse{s: http400, b: errorJsonEnd}, mocks.MockMongo{}},
		{"putConfiguration400#3", putConfiguration, testRequest{m: "PUT", e: "/configure", b: io.NopCloser(strings.NewReader("{\"lookahead\": 172800000000000, \"time\": \"19/00\"}"))}, testResponse{s: http400, b: "Invalid time format: 19/00"}, mocks.MockMongo{}},
		{"putConfiguration500#1", putConfiguration, testRequest{m: "PUT", e: "/configure", b: io.NopCloser(errReader(0))}, testResponse{s: http500, b: errorIoutilReadAll}, mocks.MockMongo{}},
		{"putConfiguration500#2", putConfiguration, testRequest{m: "PUT", e: "/configure", b: io.NopCloser(strings.NewReader("{\"lookahead\": \"172800000000000\", \"silence\": false, \"time\": \"19:00\"}"))}, testResponse{s: http500, b: "json: cannot unmarshal string into Go struct field .lookahead of type time.Duration"}, mocks.MockMongo{}},
		{"putConfiguration500#3", putConfiguration, testRequest{m: "PUT", e: "/configure", b: io.NopCloser(strings.NewReader("{\"lookahead\":172800000000000,\"silence\":false,\"time\":\"18:0z\"}"))}, testResponse{s: http500, b: "the given time format is not supported"}, mocks.MockMongo{}},
		{"getCookable200#1", getCookable, testRequest{m: "GET", e: "/getCookable"}, testResponse{s: http200, b: bodyEmpty}, mocks.MockMongo{}},
		{"getCookable500#1", getCookable, testRequest{m: "GET", e: "/getCookable"}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"getCookable500#2", getCookable, testRequest{m: "GET", e: "/getCookable"}, testResponse{s: http500, b: errorDecodeFail}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail}},
		{"getExpired200#1", getExpired, testRequest{m: "GET", e: "/expired"}, testResponse{s: http200, b: bodyEmpty}, mocks.MockMongo{}},
		{"getExpired500#1", getExpired, testRequest{m: "GET", e: "/expired"}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"getExpired500#2", getExpired, testRequest{m: "GET", e: "/expired"}, testResponse{s: http500, b: errorDecodeFail}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail}},
		{"getExpiring200#1", getExpiring, testRequest{m: "GET", e: "/expiring", qp: queryParams1020}, testResponse{s: http200, b: bodyEmpty}, mocks.MockMongo{}},
		{"getExpiring200#2", getExpiring, testRequest{m: "GET", e: "/expiring", qp: queryParams2030}, testResponse{s: http200, b: bodyExpiring}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsIngredientRange}},
		{"getExpiring400#1", getExpiring, testRequest{m: "GET", e: "/expiring", qp: map[string]string{"from": "x", "to": "20"}}, testResponse{s: http400, b: errorStrconvX}, mocks.MockMongo{}},
		{"getExpiring400#2", getExpiring, testRequest{m: "GET", e: "/expiring", qp: map[string]string{"from": "10", "to": "y"}}, testResponse{s: http400, b: errorStrconvY}, mocks.MockMongo{}},
		{"getExpiring500#1", getExpiring, testRequest{m: "GET", e: "/expiring"}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"getExpiring500#2", getExpiring, testRequest{m: "GET", e: "/expiring"}, testResponse{s: http500, b: errorDecodeFail}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail}},
		{"getOneDocument200#1a", getOneDocument, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredientsDoc}, testResponse{s: http200, b: "null"}, mocks.MockMongo{}},
		{"getOneDocument200#1b", getOneDocument, testRequest{m: "GET", e: "/documents", rv: routeVarsRecipesDoc}, testResponse{s: http200, b: "{\"ingredients\":[\"hello\"],\"isCookable\":true}"}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentRecipe, OverrideFindManyDocuments: OverrideFindManyDocumentsIngredient}},
		{"getOneDocument400#1", getOneDocument, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredientsDocInvalid}, testResponse{s: http400, b: errorDocumentIdInvalid}, mocks.MockMongo{}},
		{"getOneDocument404#1", getOneDocument, testRequest{m: "GET", e: "/documents", rv: routeVarsInvalidDoc}, testResponse{s: http404, b: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"getOneDocument404#2", getOneDocument, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredientsDoc}, testResponse{s: http404, b: utils.ErrorMongoNoDocuments}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentNone}},
		{"getOneDocument500#1", getOneDocument, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredientsDoc}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"getOneDocument500#2", getOneDocument, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredientsDocEncodeFail}, testResponse{s: http500, b: errorDocumentIdEncodeFail}, mocks.MockMongo{}},
		{"getOneDocument500#3", getOneDocument, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredientsDoc}, testResponse{s: http500, b: errorDecodeFail}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentErrorDecodeFail}},
		{"getOneDocument500#4", getOneDocument, testRequest{m: "GET", e: "/documents", rv: routeVarsRecipesDoc}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentRecipe, OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"getOneDocument500#5", getOneDocument, testRequest{m: "GET", e: "/documents", rv: routeVarsRecipesDoc}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentRecipe, OverrideFindManyDocuments: OverrideFindManyDocumentsIngredient, OverrideUpdateOneDocument: OverrideUpdateOneDocumentErrorBasic}},
		{"getOneDocument500#6", getOneDocument, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredientsDoc}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentErrorBasic}},
		{"getManyDocuments200#1a", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredients, qp: queryParamsAll2030}, testResponse{s: http200, b: bodyExpiring}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsIngredientRange}},
		{"getManyDocuments200#1b", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsRecipes, qp: queryParamsAll2030}, testResponse{s: http200, b: "[{\"ingredients\":[\"hello\"],\"isCookable\":true}]"}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuper}},
		{"getManyDocuments200#2b", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsRecipes, qp: map[string]string{"name": "hello", "isCookable": "true"}}, testResponse{s: http200, b: "[]"}, mocks.MockMongo{}},
		{"getManyDocuments400#1", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsRecipes, qp: map[string]string{"name": "hello", "isCookable": "lol"}}, testResponse{s: http400, b: "strconv.ParseBool: parsing \"lol\": invalid syntax"}, mocks.MockMongo{}},
		{"getManyDocuments400#2", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredients, qp: map[string]string{"name": "hello", "haveStocked": "false", "from": "x", "to": ""}}, testResponse{s: http400, b: errorStrconvX}, mocks.MockMongo{}},
		{"getManyDocuments400#3", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredients, qp: map[string]string{"name": "hello", "haveStocked": "false", "from": "10", "to": "y"}}, testResponse{s: http400, b: errorStrconvY}, mocks.MockMongo{}},
		{"getManyDocuments400#4", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredients, qp: map[string]string{"name": "hello", "haveStocked": "lol", "from": "10", "to": "20"}}, testResponse{s: http400, b: "strconv.ParseBool: parsing \"lol\": invalid syntax"}, mocks.MockMongo{}},
		{"getManyDocuments404#1", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsInvalid, qp: queryParamsAll1020}, testResponse{s: http404, b: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"getManyDocuments500#1", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredients, qp: queryParamsAll1020}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"getManyDocuments500#2", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredients, qp: queryParamsAll1020}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"getManyDocuments500#3", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsRecipes, qp: queryParamsAll2030}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsRecipeOrErrorBasic}},
		{"getManyDocuments500#4", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsRecipes, qp: queryParamsAll2030}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuper, OverrideUpdateOneDocument: OverrideUpdateOneDocumentErrorBasic}},
		{"getManyDocuments500#5", getManyDocuments, testRequest{m: "GET", e: "/documents", rv: routeVarsIngredients, qp: queryParamsAll1020}, testResponse{s: http500, b: errorDecodeFail}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail}},
		{"postManyDocuments200#1a", postManyDocuments, testRequest{m: "POST", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(strings.NewReader(documentsBasic))}, testResponse{s: http201}, mocks.MockMongo{}},
		{"postManyDocuments200#1b", postManyDocuments, testRequest{m: "POST", e: "/documents", rv: routeVarsRecipes, b: io.NopCloser(strings.NewReader("[{\"name\": \"Document\", \"ingredients\": []}]"))}, testResponse{s: http201}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsIngredient}},
		{"postManyDocuments400#1", postManyDocuments, testRequest{m: "POST", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(strings.NewReader("{:}"))}, testResponse{s: http400, b: errorJsonUndecodable}, mocks.MockMongo{}},
		{"postManyDocuments404#1", postManyDocuments, testRequest{m: "POST", e: "/documents", rv: routeVarsInvalid, b: io.NopCloser(strings.NewReader(documentsBasic))}, testResponse{s: http404, b: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"postManyDocuments500#1", postManyDocuments, testRequest{m: "POST", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(strings.NewReader(documentsBasic))}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"postManyDocuments500#2", postManyDocuments, testRequest{m: "POST", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(errReader(0))}, testResponse{s: http500, b: errorIoutilReadAll}, mocks.MockMongo{}},
		{"postManyDocuments500#3", postManyDocuments, testRequest{m: "POST", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(strings.NewReader(documentEmpty))}, testResponse{s: http500, b: "json: cannot unmarshal object into Go value of type []primitive.M"}, mocks.MockMongo{}},
		{"postManyDocuments500#4", postManyDocuments, testRequest{m: "POST", e: "/documents", rv: routeVarsRecipes, b: io.NopCloser(strings.NewReader("[{\"name\": \"Document\", \"ingredients\": [\"hello\"]}]"))}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"postManyDocuments500#5", postManyDocuments, testRequest{m: "POST", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(strings.NewReader(documentsBasic))}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideInsertManyDocuments: OverrideInsertManyDocumentsErrorBasic}},
		{"putOneDocument200#1", putOneDocument, testRequest{m: "PUT", e: "/documents", rv: routeVarsIngredientsDoc, b: io.NopCloser(strings.NewReader(documentBasic))}, testResponse{s: http200}, mocks.MockMongo{}},
		{"putOneDocument400#1", putOneDocument, testRequest{m: "PUT", e: "/documents", rv: routeVarsIngredientsDocInvalid}, testResponse{s: http400, b: errorDocumentIdInvalid}, mocks.MockMongo{}},
		{"putOneDocument400#2", putOneDocument, testRequest{m: "PUT", e: "/documents", rv: routeVarsIngredientsDoc, b: io.NopCloser(strings.NewReader("{:}"))}, testResponse{s: http400, b: errorJsonUndecodable}, mocks.MockMongo{}},
		{"putOneDocument400#3", putOneDocument, testRequest{m: "PUT", e: "/documents", rv: routeVarsIngredientsDoc, b: io.NopCloser(strings.NewReader(documentEmpty))}, testResponse{s: http400, b: errorEmptyUpdateInstructions}, mocks.MockMongo{OverrideUpdateOneDocument: OverrideUpdateOneDocumentEmptyUpdate}},
		{"putOneDocument404#1", putOneDocument, testRequest{m: "PUT", e: "/documents", rv: routeVarsInvalidDoc, b: io.NopCloser(strings.NewReader(documentBasic))}, testResponse{s: http404, b: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"putOneDocument404#2", putOneDocument, testRequest{m: "PUT", e: "/documents", rv: routeVarsIngredientsDoc, b: io.NopCloser(strings.NewReader(documentEmpty))}, testResponse{s: http404}, mocks.MockMongo{OverrideUpdateOneDocument: OverrideUpdateOneDocumentZero}},
		{"putOneDocument500#1", putOneDocument, testRequest{m: "PUT", e: "/documents", rv: routeVarsIngredientsDoc, b: io.NopCloser(strings.NewReader(documentBasic))}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"putOneDocument500#2", putOneDocument, testRequest{m: "PUT", e: "/documents", rv: routeVarsIngredientsDocEncodeFail}, testResponse{s: http500, b: errorDocumentIdEncodeFail}, mocks.MockMongo{}},
		{"putOneDocument500#3", putOneDocument, testRequest{m: "PUT", e: "/documents", rv: routeVarsIngredientsDoc, b: io.NopCloser(errReader(0))}, testResponse{s: http500, b: errorIoutilReadAll}, mocks.MockMongo{}},
		{"putOneDocument500#4", putOneDocument, testRequest{m: "PUT", e: "/documents", rv: routeVarsIngredientsDoc, b: io.NopCloser(strings.NewReader("[{}"))}, testResponse{s: http500, b: errorJsonEnd}, mocks.MockMongo{}},
		{"putOneDocument500#5", putOneDocument, testRequest{m: "PUT", e: "/documents", rv: routeVarsIngredientsDoc, b: io.NopCloser(strings.NewReader(documentBasic))}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideUpdateOneDocument: OverrideUpdateOneDocumentErrorBasic}},
		{"deleteOneDocument200#1", deleteOneDocument, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredientsDoc}, testResponse{s: http200}, mocks.MockMongo{}},
		{"deleteOneDocument400#1", deleteOneDocument, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredientsDocInvalid}, testResponse{s: http400, b: errorDocumentIdInvalid}, mocks.MockMongo{}},
		{"deleteOneDocument404#1", deleteOneDocument, testRequest{m: "DELETE", e: "/documents", rv: routeVarsInvalidDoc}, testResponse{s: http404, b: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"deleteOneDocument404#2", deleteOneDocument, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredientsDoc}, testResponse{s: http404, b: utils.ErrorMongoNoDocuments}, mocks.MockMongo{OverrideDeleteOneDocument: OverrideDeleteOneDocumentNone}},
		{"deleteOneDocument500#1", deleteOneDocument, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredientsDoc}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"deleteOneDocument500#2", deleteOneDocument, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredientsDocEncodeFail}, testResponse{s: http500, b: errorDocumentIdEncodeFail}, mocks.MockMongo{}},
		{"deleteOneDocument500#3", deleteOneDocument, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredientsDoc}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideDeleteOneDocument: OverrideDeleteOneDocumentErrorBasic}},
		{"deleteManyDocuments200#1", deleteManyDocuments, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(strings.NewReader(documentIds))}, testResponse{s: http200}, mocks.MockMongo{}},
		{"deleteManyDocuments400#1", deleteManyDocuments, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(strings.NewReader("{:}"))}, testResponse{s: http400, b: errorJsonUndecodable}, mocks.MockMongo{}},
		{"deleteManyDocuments400#2", deleteManyDocuments, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(strings.NewReader("[\"hello\"]"))}, testResponse{s: http400, b: errorDocumentIdInvalid}, mocks.MockMongo{}},
		{"deleteManyDocuments404#1", deleteManyDocuments, testRequest{m: "DELETE", e: "/documents", rv: routeVarsInvalid, b: io.NopCloser(strings.NewReader(documentIds))}, testResponse{s: http404, b: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"deleteManyDocuments404#2", deleteManyDocuments, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(strings.NewReader(documentIds))}, testResponse{s: http404, b: "no documents found"}, mocks.MockMongo{OverrideDeleteManyDocuments: OverrideDeleteManyDocumentsZero}},
		{"deleteManyDocuments500#1", deleteManyDocuments, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(strings.NewReader(documentIds))}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"deleteManyDocuments500#2", deleteManyDocuments, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredientsDoc, b: io.NopCloser(errReader(0))}, testResponse{s: http500, b: errorIoutilReadAll}, mocks.MockMongo{}},
		{"deleteManyDocuments500#3", deleteManyDocuments, testRequest{m: "DELETE", e: "/documents", rv: routeVarsIngredients, b: io.NopCloser(strings.NewReader(documentIds))}, testResponse{s: http500, b: errorBasic}, mocks.MockMongo{OverrideDeleteManyDocuments: OverrideDeleteManyDocumentsErrorBasic}},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			configuration.Mongo = &st.mongoClient

			req, err := http.NewRequest(st.request.m, st.request.e, st.request.b)
			if err != nil {
				t.Fatal(err)
			}

			if st.request.qp != nil {
				q := req.URL.Query()
				for k, v := range st.request.qp {
					q.Add(k, v)
				}
				req.URL.RawQuery = q.Encode()
			}

			if st.request.rv != nil {
				req = mux.SetURLVars(req, st.request.rv)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(st.handler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != st.response.s {
				t.Errorf("handler returned wrong status code: got %v want %v", status, st.response.s)
			}

			if rr.Body.String() != st.response.b {
				t.Errorf("handler returned unexpected body: got \"%v\" want \"%v\"", rr.Body.String(), st.response.b)
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
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.ErrorLevel},
			[]string{"Failed to identify expired items"},
		},
		{
			// Error #2, Could not obtain expiring items
			"checkExpirationsError#2",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsCheckExpirations2},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.ErrorLevel},
			[]string{"Failed to identify expiring items"},
		},
		{
			// Success #1, No expired/expiring items, no need to proceed.
			"checkExpirationsSuccess#1",
			mocks.MockMongo{},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel},
			[]string{"Restocking not required"},
		},
		{
			// Success #2, items expired/expiring added to existing Trello card and SMS message sent.
			"checkExpirationsSuccess#2",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.InfoLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Added to Trello card", "Sent Twilio message"},
		},
		{
			// Error #3, items expired/expiring but could not obtain Trello card, SMS message still sent.
			"checkExpirationsError#3",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideGetShoppingList: OverrideGetShoppingListErrorBasic},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Failed to get Trello card", "Sent Twilio message"},
		},
		{
			// Success #3, items expired/expiring added to new Trello card and SMS message sent.
			"checkExpirationsSuccess#3",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideGetShoppingList: OverrideGetShoppingListNil},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.InfoLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Created Trello card", "Sent Twilio message"},
		},
		{
			// Error #4, items expired/expiring but could not add to existing card Trello card, SMS message still sent.
			"checkExpirationsError#4",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideAddToShoppingList: OverrideAddToShoppingListErroBasic},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Failed to add to Trello card", "Sent Twilio message"},
		},
		{
			// Error #5, items expired/expiring but could not create new card Trello card, SMS message still sent.
			"checkExpirationsError#5",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideGetShoppingList: OverrideGetShoppingListNil, OverrideCreateShoppingList: OverrideCreateShoppingListErrorBasic},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Failed to create Trello card", "Sent Twilio message"},
		},
		{
			// Error #6, items expired/expiring but could not create new card Trello card or send SMS message.
			"checkExpirationsError#6",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideGetShoppingList: OverrideGetShoppingListNil, OverrideCreateShoppingList: OverrideCreateShoppingListErrorBasic},
			mocks.MockTwilio{OverrideComposeMessage: OverrideComposeMessageEmpty, OverrideSendMessage: OverrideSendMessageErrorBasic},
			[]logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.ErrorLevel},
			[]string{"Restocking required", "Failed to create Trello card", "Failed to send Twilio message"},
		},
		{
			// Success #4, items expired/expiring added to new Trello card and SMS message skipped.
			"checkExpirationsSuccess#4",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideGetShoppingList: OverrideGetShoppingListNil},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.InfoLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Created Trello card", "Skipped Twilio message"},
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

			if st.name == "checkExpirationsSuccess#4" {
				configuration.Silence = true
			} else {
				configuration.Silence = false
			}

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
		logrus.SetOutput(io.Discard)
	})

	t.Run("isCookable", func(t *testing.T) {
		ctx := context.Background()

		mcErr := mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}
		mc := mocks.MockMongo{
			OverrideFindManyDocuments: func(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
				return []primitive.M{}, nil
			},
		}

		cases := []struct {
			mock   mocks.MockMongo
			recipe primitive.M
			want   bool
			err    error
		}{
			{mc, primitive.M{"_id": "hello"}, false, nil},
			{mcErr, primitive.M{"_id": "hello", "ingredients": []interface{}{}}, false, fmt.Errorf(errorBasic)},
			{mc, primitive.M{"_id": "hello", "ingredients": []interface{}{}}, true, nil},
		}
		for _, c := range cases {
			configuration.Mongo = &c.mock
			got, err := isCookable(ctx, &c.recipe)
			if got != c.want {
				t.Errorf("isCookable(\"%+v\"), got (\"%t\", \"%s\"), want (\"%t\", \"%s\")", c.recipe, got, err, c.want, c.err)
			}
		}
	})
}
