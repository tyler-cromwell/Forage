package api

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/config"
	"github.com/tyler-cromwell/forage/tests/mocks"
	"github.com/tyler-cromwell/forage/utils"
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
	return 0, errors.New(errorIoReadAll)
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
		{
			/*
			 */
			"getConfiguration200#1",
			getConfiguration,
			testRequest{
				method:   "GET",
				endpoint: "/configure",
			},
			testResponse{
				status: http.StatusOK,
				body:   "{\"lookahead\":172800000000000,\"silence\":false,\"time\":\"\"}",
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putConfiguration200#1",
			putConfiguration,
			testRequest{
				method:   "PUT",
				endpoint: "/configure",
				body:     io.NopCloser(strings.NewReader("{\"lookahead\": 172800000000000, \"time\": \"19:00\"}")),
			},
			testResponse{
				status: http.StatusOK,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putConfiguration400#1",
			putConfiguration,
			testRequest{
				method:   "PUT",
				endpoint: "/configure",
				body:     io.NopCloser(strings.NewReader("{:}")),
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorJsonUndecodable,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putConfiguration400#2",
			putConfiguration,
			testRequest{
				method:   "PUT",
				endpoint: "/configure",
				body:     io.NopCloser(strings.NewReader("")),
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorJsonEnd,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putConfiguration400#3",
			putConfiguration,
			testRequest{
				method:   "PUT",
				endpoint: "/configure",
				body:     io.NopCloser(strings.NewReader("{\"lookahead\": 172800000000000, \"time\": \"19/00\"}")),
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   "Invalid time format: 19/00",
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putConfiguration500#1",
			putConfiguration,
			testRequest{
				method:   "PUT",
				endpoint: "/configure",
				body:     io.NopCloser(errReader(0)),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorIoReadAll,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putConfiguration500#2",
			putConfiguration,
			testRequest{
				method:   "PUT",
				endpoint: "/configure",
				body:     io.NopCloser(strings.NewReader("{\"lookahead\": \"172800000000000\", \"silence\": false, \"time\": \"19:00\"}")),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   "json: cannot unmarshal string into Go struct field .lookahead of type time.Duration",
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putConfiguration500#3",
			putConfiguration,
			testRequest{
				method:   "PUT",
				endpoint: "/configure",
				body:     io.NopCloser(strings.NewReader("{\"lookahead\":172800000000000,\"silence\":false,\"time\":\"18:0z\"}")),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   "the given time format is not supported",
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getCookable200#1",
			getCookable,
			testRequest{
				method:   "GET",
				endpoint: "/getCookable",
			},
			testResponse{
				status: http.StatusOK,
				body:   bodyEmpty,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getCookable500#1",
			getCookable,
			testRequest{
				method:   "GET",
				endpoint: "/getCookable",
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic,
			},
		},
		{
			/*
			 */
			"getCookable500#2",
			getCookable,
			testRequest{
				method:   "GET",
				endpoint: "/getCookable",
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorDecodeFail,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail,
			},
		},
		{
			/*
			 */
			"getExpired200#1",
			getExpired,
			testRequest{
				method:   "GET",
				endpoint: "/expired",
			},
			testResponse{
				status: http.StatusOK,
				body:   bodyEmpty,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getExpired500#1",
			getExpired,
			testRequest{
				method:   "GET",
				endpoint: "/expired",
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic,
			},
		},
		{
			/*
			 */
			"getExpired500#2",
			getExpired,
			testRequest{
				method:   "GET",
				endpoint: "/expired",
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorDecodeFail,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail,
			},
		},
		{
			/*
			 */
			"getExpiring200#1",
			getExpiring,
			testRequest{
				method:          "GET",
				endpoint:        "/expiring",
				queryParameters: queryParams1020,
			},
			testResponse{
				status: http.StatusOK,
				body:   bodyEmpty,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getExpiring200#2",
			getExpiring,
			testRequest{
				method:          "GET",
				endpoint:        "/expiring",
				queryParameters: queryParams2030,
			},
			testResponse{
				status: http.StatusOK,
				body:   bodyExpiring,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsIngredientRange,
			},
		},
		{
			/*
			 */
			"getExpiring400#1",
			getExpiring,
			testRequest{
				method:   "GET",
				endpoint: "/expiring",
				queryParameters: map[string]string{
					"from": "x",
					"to":   "20",
				},
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorStrconvX,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getExpiring400#2",
			getExpiring,
			testRequest{
				method:   "GET",
				endpoint: "/expiring",
				queryParameters: map[string]string{
					"from": "10",
					"to":   "y",
				},
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorStrconvY,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getExpiring500#1",
			getExpiring,
			testRequest{
				method:   "GET",
				endpoint: "/expiring",
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic,
			},
		},
		{
			/*
			 */
			"getExpiring500#2",
			getExpiring,
			testRequest{
				method:   "GET",
				endpoint: "/expiring",
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorDecodeFail,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail,
			},
		},
		{
			/*
			 */
			"getOneDocument200#1a",
			getOneDocument,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
			},
			testResponse{
				status: http.StatusOK,
				body:   "null",
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getOneDocument200#1b",
			getOneDocument,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsRecipesDoc,
			},
			testResponse{
				status: http.StatusOK,
				body:   bodyCookable,
			},
			mocks.MockMongo{
				OverrideFindOneDocument:   OverrideFindOneDocumentRecipe,
				OverrideFindManyDocuments: OverrideFindManyDocumentsIngredient,
			},
		},
		{
			/*
			 */
			"getOneDocument400#1",
			getOneDocument,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDocInvalid,
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorDocumentIdInvalid,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getOneDocument404#1",
			getOneDocument,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsInvalidDoc,
			},
			testResponse{
				status: http.StatusNotFound,
				body:   errorCollectionIdInvalid,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getOneDocument404#2",
			getOneDocument,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
			},
			testResponse{
				status: http.StatusNotFound,
				body:   utils.ErrorMongoNoDocuments,
			},
			mocks.MockMongo{
				OverrideFindOneDocument: OverrideFindOneDocumentNone,
			},
		},
		{
			/*
			 */
			"getOneDocument500#1",
			getOneDocument,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideCollections: OverrideCollectionsErrorBasic,
			},
		},
		{
			/*
			 */
			"getOneDocument500#2",
			getOneDocument,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDocEncodeFail,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorDocumentIdEncodeFail,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getOneDocument500#3",
			getOneDocument,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorDecodeFail,
			},
			mocks.MockMongo{
				OverrideFindOneDocument: OverrideFindOneDocumentErrorDecodeFail,
			},
		},
		{
			/*
			 */
			"getOneDocument500#4",
			getOneDocument,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsRecipesDoc,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideFindOneDocument:   OverrideFindOneDocumentRecipe,
				OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic,
			},
		},
		{
			/*
			 */
			"getOneDocument500#5",
			getOneDocument,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsRecipesDoc,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideFindOneDocument:   OverrideFindOneDocumentRecipe,
				OverrideFindManyDocuments: OverrideFindManyDocumentsIngredient,
				OverrideUpdateOneDocument: OverrideUpdateOneDocumentErrorBasic,
			},
		},
		{
			/*
			 */
			"getOneDocument500#6",
			getOneDocument,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideFindOneDocument: OverrideFindOneDocumentErrorBasic,
			},
		},
		{
			/*
			 */
			"getManyDocuments200#1a",
			getManyDocuments,
			testRequest{
				method:          "GET",
				endpoint:        "/documents",
				routeVariables:  routeVarsIngredients,
				queryParameters: queryParamsAll2030,
			},
			testResponse{
				status: http.StatusOK,
				body:   bodyExpiring,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsIngredientRange,
			},
		},
		{
			/*
			 */
			"getManyDocuments200#1b",
			getManyDocuments,
			testRequest{
				method:          "GET",
				endpoint:        "/documents",
				routeVariables:  routeVarsRecipes,
				queryParameters: queryParamsAll2030,
			},
			testResponse{
				status: http.StatusOK,
				body:   bodyCookables,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsSuper,
			},
		},
		{
			/*
			 */
			"getManyDocuments200#2b",
			getManyDocuments,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsRecipes,
				queryParameters: map[string]string{
					"name":       "hello",
					"isCookable": "true",
				},
			},
			testResponse{
				status: http.StatusOK,
				body:   bodyEmpty,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getManyDocuments400#1",
			getManyDocuments,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsRecipes,
				queryParameters: map[string]string{
					"name":       "hello",
					"isCookable": "lol",
				},
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorStrconvLol,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getManyDocuments400#2",
			getManyDocuments,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				queryParameters: map[string]string{
					"name":        "hello",
					"haveStocked": "false",
					"from":        "x",
					"to":          "",
				},
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorStrconvX,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getManyDocuments400#3",
			getManyDocuments,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				queryParameters: map[string]string{
					"name":        "hello",
					"haveStocked": "false",
					"from":        "10",
					"to":          "y",
				},
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorStrconvY,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getManyDocuments400#4",
			getManyDocuments,
			testRequest{
				method:         "GET",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				queryParameters: map[string]string{
					"name":        "hello",
					"haveStocked": "lol",
					"from":        "10",
					"to":          "20",
				},
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorStrconvLol,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getManyDocuments404#1",
			getManyDocuments,
			testRequest{
				method:          "GET",
				endpoint:        "/documents",
				routeVariables:  routeVarsInvalid,
				queryParameters: queryParamsAll1020,
			},
			testResponse{
				status: http.StatusNotFound,
				body:   errorCollectionIdInvalid,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"getManyDocuments500#1",
			getManyDocuments,
			testRequest{
				method:          "GET",
				endpoint:        "/documents",
				routeVariables:  routeVarsIngredients,
				queryParameters: queryParamsAll1020,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideCollections: OverrideCollectionsErrorBasic,
			},
		},
		{
			/*
			 */
			"getManyDocuments500#2",
			getManyDocuments,
			testRequest{
				method:          "GET",
				endpoint:        "/documents",
				routeVariables:  routeVarsIngredients,
				queryParameters: queryParamsAll1020,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic,
			},
		},
		{
			/*
			 */
			"getManyDocuments500#3",
			getManyDocuments,
			testRequest{
				method:          "GET",
				endpoint:        "/documents",
				routeVariables:  routeVarsRecipes,
				queryParameters: queryParamsAll2030,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsRecipeOrErrorBasic,
			},
		},
		{
			/*
			 */
			"getManyDocuments500#4",
			getManyDocuments,
			testRequest{
				method:          "GET",
				endpoint:        "/documents",
				routeVariables:  routeVarsRecipes,
				queryParameters: queryParamsAll2030,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsSuper,
				OverrideUpdateOneDocument: OverrideUpdateOneDocumentErrorBasic,
			},
		},
		{
			/*
			 */
			"getManyDocuments500#5",
			getManyDocuments,
			testRequest{
				method:          "GET",
				endpoint:        "/documents",
				routeVariables:  routeVarsIngredients,
				queryParameters: queryParamsAll1020,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorDecodeFail,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsDecodeFail,
			},
		},
		{
			/*
			 */
			"postManyDocuments200#1a",
			postManyDocuments,
			testRequest{
				method:         "POST",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(strings.NewReader(documentsBasic)),
			},
			testResponse{
				status: http.StatusCreated,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"postManyDocuments200#1b",
			postManyDocuments,
			testRequest{
				method:         "POST",
				endpoint:       "/documents",
				routeVariables: routeVarsRecipes,
				body:           io.NopCloser(strings.NewReader("[{\"name\": \"Document\", \"ingredients\": []}]")),
			},
			testResponse{
				status: http.StatusCreated,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsIngredient,
			},
		},
		{
			/*
			 */
			"postManyDocuments400#1",
			postManyDocuments,
			testRequest{
				method:         "POST",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(strings.NewReader("{:}")),
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorJsonUndecodable,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"postManyDocuments404#1",
			postManyDocuments,
			testRequest{
				method:         "POST",
				endpoint:       "/documents",
				routeVariables: routeVarsInvalid,
				body:           io.NopCloser(strings.NewReader(documentsBasic)),
			},
			testResponse{
				status: http.StatusNotFound,
				body:   errorCollectionIdInvalid,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"postManyDocuments500#1",
			postManyDocuments,
			testRequest{
				method:         "POST",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(strings.NewReader(documentsBasic)),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideCollections: OverrideCollectionsErrorBasic,
			},
		},
		{
			/*
			 */
			"postManyDocuments500#2",
			postManyDocuments,
			testRequest{
				method:         "POST",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(errReader(0)),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorIoReadAll,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"postManyDocuments500#3",
			postManyDocuments,
			testRequest{
				method:         "POST",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(strings.NewReader(documentEmpty)),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   "json: cannot unmarshal object into Go value of type []primitive.M",
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"postManyDocuments500#4",
			postManyDocuments,
			testRequest{
				method:         "POST",
				endpoint:       "/documents",
				routeVariables: routeVarsRecipes,
				body:           io.NopCloser(strings.NewReader("[{\"name\": \"Document\", \"ingredients\": [\"hello\"]}]")),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic,
			},
		},
		{
			/*
			 */
			"postManyDocuments500#5",
			postManyDocuments,
			testRequest{
				method:         "POST",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(strings.NewReader(documentsBasic)),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideInsertManyDocuments: OverrideInsertManyDocumentsErrorBasic,
			},
		},
		{
			/*
			 */
			"putOneDocument200#1",
			putOneDocument,
			testRequest{
				method:         "PUT",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
				body:           io.NopCloser(strings.NewReader(documentBasic)),
			},
			testResponse{
				status: http.StatusOK,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putOneDocument400#1",
			putOneDocument,
			testRequest{
				method:         "PUT",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDocInvalid,
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorDocumentIdInvalid,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putOneDocument400#2",
			putOneDocument,
			testRequest{
				method:         "PUT",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
				body:           io.NopCloser(strings.NewReader("{:}")),
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorJsonUndecodable,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putOneDocument400#3",
			putOneDocument,
			testRequest{
				method:         "PUT",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
				body:           io.NopCloser(strings.NewReader(documentEmpty)),
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorEmptyUpdateInstructions,
			},
			mocks.MockMongo{
				OverrideUpdateOneDocument: OverrideUpdateOneDocumentEmptyUpdate,
			},
		},
		{
			/*
			 */
			"putOneDocument404#1",
			putOneDocument,
			testRequest{
				method:         "PUT",
				endpoint:       "/documents",
				routeVariables: routeVarsInvalidDoc,
				body:           io.NopCloser(strings.NewReader(documentBasic)),
			},
			testResponse{
				status: http.StatusNotFound,
				body:   errorCollectionIdInvalid,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putOneDocument404#2",
			putOneDocument,
			testRequest{
				method:         "PUT",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
				body:           io.NopCloser(strings.NewReader(documentEmpty)),
			},
			testResponse{
				status: http.StatusNotFound,
			},
			mocks.MockMongo{
				OverrideUpdateOneDocument: OverrideUpdateOneDocumentZero,
			},
		},
		{
			/*
			 */
			"putOneDocument500#1",
			putOneDocument,
			testRequest{
				method:         "PUT",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
				body:           io.NopCloser(strings.NewReader(documentBasic)),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideCollections: OverrideCollectionsErrorBasic,
			},
		},
		{
			/*
			 */
			"putOneDocument500#2",
			putOneDocument,
			testRequest{
				method:         "PUT",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDocEncodeFail,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorDocumentIdEncodeFail,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putOneDocument500#3",
			putOneDocument,
			testRequest{
				method:         "PUT",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
				body:           io.NopCloser(errReader(0)),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorIoReadAll,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putOneDocument500#4",
			putOneDocument,
			testRequest{
				method:         "PUT",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
				body:           io.NopCloser(strings.NewReader("[{}")),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorJsonEnd,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"putOneDocument500#5",
			putOneDocument,
			testRequest{
				method:         "PUT",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
				body:           io.NopCloser(strings.NewReader(documentBasic)),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideUpdateOneDocument: OverrideUpdateOneDocumentErrorBasic,
			},
		},
		{
			/*
			 */
			"deleteOneDocument200#1",
			deleteOneDocument,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
			},
			testResponse{
				status: http.StatusOK,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"deleteOneDocument400#1",
			deleteOneDocument,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDocInvalid,
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorDocumentIdInvalid,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"deleteOneDocument404#1",
			deleteOneDocument,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsInvalidDoc,
			},
			testResponse{
				status: http.StatusNotFound,
				body:   errorCollectionIdInvalid,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"deleteOneDocument404#2",
			deleteOneDocument,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
			},
			testResponse{
				status: http.StatusNotFound,
				body:   utils.ErrorMongoNoDocuments,
			},
			mocks.MockMongo{
				OverrideDeleteOneDocument: OverrideDeleteOneDocumentNone,
			},
		},
		{
			/*
			 */
			"deleteOneDocument500#1",
			deleteOneDocument,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideCollections: OverrideCollectionsErrorBasic,
			},
		},
		{
			/*
			 */
			"deleteOneDocument500#2",
			deleteOneDocument,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDocEncodeFail,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorDocumentIdEncodeFail,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"deleteOneDocument500#3",
			deleteOneDocument,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideDeleteOneDocument: OverrideDeleteOneDocumentErrorBasic,
			},
		},
		{
			/*
			 */
			"deleteManyDocuments200#1",
			deleteManyDocuments,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(strings.NewReader(documentIds)),
			},
			testResponse{
				status: http.StatusOK,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"deleteManyDocuments400#1",
			deleteManyDocuments,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(strings.NewReader("{:}")),
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorJsonUndecodable,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"deleteManyDocuments400#2",
			deleteManyDocuments,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(strings.NewReader("[\"hello\"]")),
			},
			testResponse{
				status: http.StatusBadRequest,
				body:   errorDocumentIdInvalid,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"deleteManyDocuments404#1",
			deleteManyDocuments,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsInvalid,
				body:           io.NopCloser(strings.NewReader(documentIds)),
			},
			testResponse{
				status: http.StatusNotFound,
				body:   errorCollectionIdInvalid,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"deleteManyDocuments404#2",
			deleteManyDocuments,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(strings.NewReader(documentIds)),
			},
			testResponse{
				status: http.StatusNotFound,
				body:   errorNoDocuments,
			},
			mocks.MockMongo{
				OverrideDeleteManyDocuments: OverrideDeleteManyDocumentsZero,
			},
		},
		{
			/*
			 */
			"deleteManyDocuments500#1",
			deleteManyDocuments,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(strings.NewReader(documentIds)),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideCollections: OverrideCollectionsErrorBasic,
			},
		},
		{
			/*
			 */
			"deleteManyDocuments500#2",
			deleteManyDocuments,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredientsDoc,
				body:           io.NopCloser(errReader(0)),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorIoReadAll,
			},
			mocks.MockMongo{},
		},
		{
			/*
			 */
			"deleteManyDocuments500#3",
			deleteManyDocuments,
			testRequest{
				method:         "DELETE",
				endpoint:       "/documents",
				routeVariables: routeVarsIngredients,
				body:           io.NopCloser(strings.NewReader(documentIds)),
			},
			testResponse{
				status: http.StatusInternalServerError,
				body:   errorBasic,
			},
			mocks.MockMongo{
				OverrideDeleteManyDocuments: OverrideDeleteManyDocumentsErrorBasic,
			},
		},
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

	// Reverse logrus output change
	log.SetOutput(os.Stdout)
}
