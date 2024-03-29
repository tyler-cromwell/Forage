package api

import (
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
	return 0, errors.New(errorIoutilReadAll)
}

func TestAPI(t *testing.T) {
	// Discard logging output
	logrus.SetOutput(ioutil.Discard)

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
		{"getConfiguration200#1", getConfiguration, testRequest{method: "GET", endpoint: "/configure"}, testResponse{status: http.StatusOK, body: "{\"lookahead\":172800000000000,\"silence\":false,\"time\":\"\"}"}, mocks.MockMongo{}},
		{"putConfiguration200#1", putConfiguration, testRequest{method: "PUT", endpoint: "/configure", body: io.NopCloser(strings.NewReader("{\"lookahead\": 172800000000000, \"time\": \"19:00\"}"))}, testResponse{status: http.StatusOK}, mocks.MockMongo{}},
		{"putConfiguration400#1", putConfiguration, testRequest{method: "PUT", endpoint: "/configure", body: io.NopCloser(strings.NewReader("{:}"))}, testResponse{status: http.StatusBadRequest, body: errorJsonUndecodable}, mocks.MockMongo{}},
		{"putConfiguration400#2", putConfiguration, testRequest{method: "PUT", endpoint: "/configure", body: io.NopCloser(strings.NewReader(""))}, testResponse{status: http.StatusBadRequest, body: errorJsonEnd}, mocks.MockMongo{}},
		{"putConfiguration400#3", putConfiguration, testRequest{method: "PUT", endpoint: "/configure", body: io.NopCloser(strings.NewReader("{\"lookahead\": 172800000000000, \"time\": \"19/00\"}"))}, testResponse{status: http.StatusBadRequest, body: "Invalid time format: 19/00"}, mocks.MockMongo{}},
		{"putConfiguration500#1", putConfiguration, testRequest{method: "PUT", endpoint: "/configure", body: io.NopCloser(errReader(0))}, testResponse{status: http.StatusInternalServerError, body: errorIoutilReadAll}, mocks.MockMongo{}},
		{"putConfiguration500#2", putConfiguration, testRequest{method: "PUT", endpoint: "/configure", body: io.NopCloser(strings.NewReader("{\"lookahead\": \"172800000000000\", \"silence\": false, \"time\": \"19:00\"}"))}, testResponse{status: http.StatusInternalServerError, body: "json: cannot unmarshal string into Go struct field .lookahead of type time.Duration"}, mocks.MockMongo{}},
		{"putConfiguration500#3", putConfiguration, testRequest{method: "PUT", endpoint: "/configure", body: io.NopCloser(strings.NewReader("{\"lookahead\":172800000000000,\"silence\":false,\"time\":\"18:0z\"}"))}, testResponse{status: http.StatusInternalServerError, body: "the given time format is not supported"}, mocks.MockMongo{}},
		{"getCookable200#1", getCookable, testRequest{method: "GET", endpoint: "/getCookable"}, testResponse{status: http.StatusOK, body: bodyEmpty}, mocks.MockMongo{}},
		{"getCookable500#1", getCookable, testRequest{method: "GET", endpoint: "/getCookable"}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"getCookable500#2", getCookable, testRequest{method: "GET", endpoint: "/getCookable"}, testResponse{status: http.StatusInternalServerError, body: errorDecodeFail}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail}},
		{"getExpired200#1", getExpired, testRequest{method: "GET", endpoint: "/expired"}, testResponse{status: http.StatusOK, body: bodyEmpty}, mocks.MockMongo{}},
		{"getExpired500#1", getExpired, testRequest{method: "GET", endpoint: "/expired"}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"getExpired500#2", getExpired, testRequest{method: "GET", endpoint: "/expired"}, testResponse{status: http.StatusInternalServerError, body: errorDecodeFail}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail}},
		{"getExpiring200#1", getExpiring, testRequest{method: "GET", endpoint: "/expiring", queryParameters: queryParameters1020}, testResponse{status: http.StatusOK, body: bodyEmpty}, mocks.MockMongo{}},
		{"getExpiring200#2", getExpiring, testRequest{method: "GET", endpoint: "/expiring", queryParameters: queryParameters2030}, testResponse{status: http.StatusOK, body: bodyExpiring}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsIngredientRange}},
		{"getExpiring400#1", getExpiring, testRequest{method: "GET", endpoint: "/expiring", queryParameters: map[string]string{"from": "x", "to": "20"}}, testResponse{status: http.StatusBadRequest, body: errorStrconvX}, mocks.MockMongo{}},
		{"getExpiring400#2", getExpiring, testRequest{method: "GET", endpoint: "/expiring", queryParameters: map[string]string{"from": "10", "to": "y"}}, testResponse{status: http.StatusBadRequest, body: errorStrconvY}, mocks.MockMongo{}},
		{"getExpiring500#1", getExpiring, testRequest{method: "GET", endpoint: "/expiring"}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"getExpiring500#2", getExpiring, testRequest{method: "GET", endpoint: "/expiring"}, testResponse{status: http.StatusInternalServerError, body: errorDecodeFail}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail}},
		{"getOneDocument200#1a", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}}, testResponse{status: http.StatusOK, body: "null"}, mocks.MockMongo{}},
		{"getOneDocument200#1b", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionRecipes, "id": documentId}}, testResponse{status: http.StatusOK, body: "{\"ingredients\":[\"hello\"],\"isCookable\":true}"}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentRecipe, OverrideFindManyDocuments: OverrideFindManyDocumentsIngredient}},
		{"getOneDocument400#1", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentIdInvalid}}, testResponse{status: http.StatusBadRequest, body: errorDocumentIdInvalid}, mocks.MockMongo{}},
		{"getOneDocument404#1", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": collectionIdInvalid, "id": documentId}}, testResponse{status: http.StatusNotFound, body: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"getOneDocument404#2", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}}, testResponse{status: http.StatusNotFound, body: utils.ErrorMongoNoDocuments}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentNone}},
		{"getOneDocument500#1", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"getOneDocument500#2", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentIdEncodeFail}}, testResponse{status: http.StatusInternalServerError, body: errorDocumentIdEncodeFail}, mocks.MockMongo{}},
		{"getOneDocument500#3", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}}, testResponse{status: http.StatusInternalServerError, body: errorDecodeFail}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentErrorDecodeFail}},
		{"getOneDocument500#4", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionRecipes, "id": documentId}}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentRecipe, OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"getOneDocument500#5", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionRecipes, "id": documentId}}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentRecipe, OverrideFindManyDocuments: OverrideFindManyDocumentsIngredient, OverrideUpdateOneDocument: OverrideUpdateOneDocumentErrorBasic}},
		{"getOneDocument500#6", getOneDocument, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideFindOneDocument: OverrideFindOneDocumentErrorBasic}},
		{"getManyDocuments200#1a", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, queryParameters: queryParametersAll2030}, testResponse{status: http.StatusOK, body: bodyExpiring}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsIngredientRange}},
		{"getManyDocuments200#1b", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionRecipes}, queryParameters: queryParametersAll2030}, testResponse{status: http.StatusOK, body: "[{\"ingredients\":[\"hello\"],\"isCookable\":true}]"}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuper}},
		{"getManyDocuments200#2b", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionRecipes}, queryParameters: map[string]string{"name": "hello", "isCookable": "true"}}, testResponse{status: http.StatusOK, body: "[]"}, mocks.MockMongo{}},
		{"getManyDocuments400#1", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionRecipes}, queryParameters: map[string]string{"name": "hello", "isCookable": "lol"}}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseBool: parsing \"lol\": invalid syntax"}, mocks.MockMongo{}},
		{"getManyDocuments400#2", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, queryParameters: map[string]string{"name": "hello", "haveStocked": "false", "from": "x", "to": ""}}, testResponse{status: http.StatusBadRequest, body: errorStrconvX}, mocks.MockMongo{}},
		{"getManyDocuments400#3", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, queryParameters: map[string]string{"name": "hello", "haveStocked": "false", "from": "10", "to": "y"}}, testResponse{status: http.StatusBadRequest, body: errorStrconvY}, mocks.MockMongo{}},
		{"getManyDocuments400#4", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, queryParameters: map[string]string{"name": "hello", "haveStocked": "lol", "from": "10", "to": "20"}}, testResponse{status: http.StatusBadRequest, body: "strconv.ParseBool: parsing \"lol\": invalid syntax"}, mocks.MockMongo{}},
		{"getManyDocuments404#1", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": collectionIdInvalid}, queryParameters: queryParametersAll1020}, testResponse{status: http.StatusNotFound, body: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"getManyDocuments500#1", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, queryParameters: queryParametersAll1020}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"getManyDocuments500#2", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, queryParameters: queryParametersAll1020}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"getManyDocuments500#3", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionRecipes}, queryParameters: queryParametersAll2030}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsRecipeOrErrorBasic}},
		{"getManyDocuments500#4", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionRecipes}, queryParameters: queryParametersAll2030}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuper, OverrideUpdateOneDocument: OverrideUpdateOneDocumentErrorBasic}},
		{"getManyDocuments500#5", getManyDocuments, testRequest{method: "GET", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, queryParameters: queryParametersAll1020}, testResponse{status: http.StatusInternalServerError, body: errorDecodeFail}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail}},
		{"postManyDocuments200#1a", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(strings.NewReader(documentsBasic))}, testResponse{status: http.StatusCreated}, mocks.MockMongo{}},
		{"postManyDocuments200#1b", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionRecipes}, body: io.NopCloser(strings.NewReader("[{\"name\": \"Document\", \"ingredients\": []}]"))}, testResponse{status: http.StatusCreated}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsIngredient}},
		{"postManyDocuments400#1", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(strings.NewReader("{:}"))}, testResponse{status: http.StatusBadRequest, body: errorJsonUndecodable}, mocks.MockMongo{}},
		{"postManyDocuments404#1", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: map[string]string{"collection": collectionIdInvalid}, body: io.NopCloser(strings.NewReader(documentsBasic))}, testResponse{status: http.StatusNotFound, body: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"postManyDocuments500#1", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(strings.NewReader(documentsBasic))}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"postManyDocuments500#2", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(errReader(0))}, testResponse{status: http.StatusInternalServerError, body: errorIoutilReadAll}, mocks.MockMongo{}},
		{"postManyDocuments500#3", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(strings.NewReader(documentEmpty))}, testResponse{status: http.StatusInternalServerError, body: "json: cannot unmarshal object into Go value of type []primitive.M"}, mocks.MockMongo{}},
		{"postManyDocuments500#4", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionRecipes}, body: io.NopCloser(strings.NewReader("[{\"name\": \"Document\", \"ingredients\": [\"hello\"]}]"))}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic}},
		{"postManyDocuments500#5", postManyDocuments, testRequest{method: "POST", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(strings.NewReader(documentsBasic))}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideInsertManyDocuments: OverrideInsertManyDocumentsErrorBasic}},
		{"putOneDocument200#1", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}, body: io.NopCloser(strings.NewReader(documentBasic))}, testResponse{status: http.StatusOK}, mocks.MockMongo{}},
		{"putOneDocument400#1", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentIdInvalid}}, testResponse{status: http.StatusBadRequest, body: errorDocumentIdInvalid}, mocks.MockMongo{}},
		{"putOneDocument400#2", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}, body: io.NopCloser(strings.NewReader("{:}"))}, testResponse{status: http.StatusBadRequest, body: errorJsonUndecodable}, mocks.MockMongo{}},
		{"putOneDocument400#3", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}, body: io.NopCloser(strings.NewReader(documentEmpty))}, testResponse{status: http.StatusBadRequest, body: errorEmptyUpdateInstructions}, mocks.MockMongo{OverrideUpdateOneDocument: OverrideUpdateOneDocumentEmptyUpdate}},
		{"putOneDocument404#1", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"collection": collectionIdInvalid, "id": documentId}, body: io.NopCloser(strings.NewReader(documentBasic))}, testResponse{status: http.StatusNotFound, body: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"putOneDocument404#2", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}, body: io.NopCloser(strings.NewReader(documentEmpty))}, testResponse{status: http.StatusNotFound}, mocks.MockMongo{OverrideUpdateOneDocument: OverrideUpdateOneDocumentZero}},
		{"putOneDocument500#1", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}, body: io.NopCloser(strings.NewReader(documentBasic))}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"putOneDocument500#2", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentIdEncodeFail}}, testResponse{status: http.StatusInternalServerError, body: errorDocumentIdEncodeFail}, mocks.MockMongo{}},
		{"putOneDocument500#3", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}, body: io.NopCloser(errReader(0))}, testResponse{status: http.StatusInternalServerError, body: errorIoutilReadAll}, mocks.MockMongo{}},
		{"putOneDocument500#4", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}, body: io.NopCloser(strings.NewReader("[{}"))}, testResponse{status: http.StatusInternalServerError, body: errorJsonEnd}, mocks.MockMongo{}},
		{"putOneDocument500#5", putOneDocument, testRequest{method: "PUT", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}, body: io.NopCloser(strings.NewReader(documentBasic))}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideUpdateOneDocument: OverrideUpdateOneDocumentErrorBasic}},
		{"deleteOneDocument200#1", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}}, testResponse{status: http.StatusOK}, mocks.MockMongo{}},
		{"deleteOneDocument400#1", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentIdInvalid}}, testResponse{status: http.StatusBadRequest, body: errorDocumentIdInvalid}, mocks.MockMongo{}},
		{"deleteOneDocument404#1", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": collectionIdInvalid, "id": documentId}}, testResponse{status: http.StatusNotFound, body: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"deleteOneDocument404#2", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}}, testResponse{status: http.StatusNotFound, body: utils.ErrorMongoNoDocuments}, mocks.MockMongo{OverrideDeleteOneDocument: OverrideDeleteOneDocumentNone}},
		{"deleteOneDocument500#1", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"deleteOneDocument500#2", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentIdEncodeFail}}, testResponse{status: http.StatusInternalServerError, body: errorDocumentIdEncodeFail}, mocks.MockMongo{}},
		{"deleteOneDocument500#3", deleteOneDocument, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideDeleteOneDocument: OverrideDeleteOneDocumentErrorBasic}},
		{"deleteManyDocuments200#1", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(strings.NewReader(documentIds))}, testResponse{status: http.StatusOK}, mocks.MockMongo{}},
		{"deleteManyDocuments400#1", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(strings.NewReader("{:}"))}, testResponse{status: http.StatusBadRequest, body: errorJsonUndecodable}, mocks.MockMongo{}},
		{"deleteManyDocuments400#2", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(strings.NewReader("[\"hello\"]"))}, testResponse{status: http.StatusBadRequest, body: errorDocumentIdInvalid}, mocks.MockMongo{}},
		{"deleteManyDocuments404#1", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": collectionIdInvalid}, body: io.NopCloser(strings.NewReader(documentIds))}, testResponse{status: http.StatusNotFound, body: errorCollectionIdInvalid}, mocks.MockMongo{}},
		{"deleteManyDocuments404#2", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(strings.NewReader(documentIds))}, testResponse{status: http.StatusNotFound, body: "no documents found"}, mocks.MockMongo{OverrideDeleteManyDocuments: OverrideDeleteManyDocumentsZero}},
		{"deleteManyDocuments500#1", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(strings.NewReader(documentIds))}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideCollections: OverrideCollectionsErrorBasic}},
		{"deleteManyDocuments500#2", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}, body: io.NopCloser(errReader(0))}, testResponse{status: http.StatusInternalServerError, body: errorIoutilReadAll}, mocks.MockMongo{}},
		{"deleteManyDocuments500#3", deleteManyDocuments, testRequest{method: "DELETE", endpoint: "/documents", routeVariables: map[string]string{"collection": config.MongoCollectionIngredients}, body: io.NopCloser(strings.NewReader(documentIds))}, testResponse{status: http.StatusInternalServerError, body: errorBasic}, mocks.MockMongo{OverrideDeleteManyDocuments: OverrideDeleteManyDocumentsErrorBasic}},
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
		logrus.SetOutput(ioutil.Discard)
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
