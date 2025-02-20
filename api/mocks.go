package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adlio/trello"
	"github.com/tyler-cromwell/forage/config"
	"github.com/tyler-cromwell/forage/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Frequently used test values/functions
const bodyEmpty = "[]"
const bodyExpiring = "[{\"_id\":1337,\"expirationDate\":25,\"haveStocked\":\"false\",\"name\":\"hello\",\"type\":\"thing\"}]"

const collectionIdInvalid = "dfhsrgaweg"

const documentId = "6187e576abc057dac3e7d5dc"
const documentIds = "[\"" + documentId + "\"]"
const documentEmpty = "{}"
const documentBasic = "{\"_id\": \"" + documentId + "\", \"name\": \"Document\"}"
const documentsBasic = "[" + documentBasic + "]"
const documentIdInvalid = "hello"
const documentIdEncodeFail = "xxxxxxxxxxxxxxxxxxxxxxxx"

const errorBasic = "failure"
const errorCollectionIdInvalid = "collection not found: " + collectionIdInvalid
const errorDecodeFail = "json: unsupported type: chan int"
const errorDocumentIdInvalid = "the provided hex string is not a valid ObjectID"
const errorDocumentIdEncodeFail = "encoding/hex: invalid byte: U+0078 'x'"
const errorEmptyUpdateInstructions = "write exception: write errors: ['$set' is empty. You must specify a field like so: {$set: {<field>: ...}}]"
const errorIoReadAll = "ioutil.ReadAll error"
const errorJsonEnd = "unexpected end of JSON input"
const errorJsonUndecodable = "invalid character ':' looking for beginning of object key string"
const errorStrconvX = "strconv.ParseInt: parsing \"x\": invalid syntax"
const errorStrconvY = "strconv.ParseInt: parsing \"y\": invalid syntax"

const http200 = http.StatusOK
const http201 = http.StatusCreated
const http400 = http.StatusBadRequest
const http404 = http.StatusNotFound
const http500 = http.StatusInternalServerError

var queryParams1020 = map[string]string{"from": "10", "to": "20"}
var queryParams2030 = map[string]string{"from": "20", "to": "30"}
var queryParamsAll1020 = map[string]string{"name": "hello", "type": "thing", "haveStocked": "false", "from": "10", "to": "20"}
var queryParamsAll2030 = map[string]string{"name": "hello", "type": "thing", "haveStocked": "false", "from": "20", "to": "30"}

var routeVarsInvalid = map[string]string{"collection": collectionIdInvalid}
var routeVarsInvalidDoc = map[string]string{"collection": collectionIdInvalid, "id": documentId}
var routeVarsIngredients = map[string]string{"collection": config.MongoCollectionIngredients}
var routeVarsIngredientsDoc = map[string]string{"collection": config.MongoCollectionIngredients, "id": documentId}
var routeVarsIngredientsDocInvalid = map[string]string{"collection": config.MongoCollectionIngredients, "id": documentIdInvalid}
var routeVarsIngredientsDocEncodeFail = map[string]string{"collection": config.MongoCollectionIngredients, "id": documentIdEncodeFail}
var routeVarsRecipes = map[string]string{"collection": config.MongoCollectionRecipes}
var routeVarsRecipesDoc = map[string]string{"collection": config.MongoCollectionRecipes, "id": documentId}

func OverrideCollectionsErrorBasic(ctx context.Context) ([]string, error) {
	return nil, fmt.Errorf(errorBasic)
}

func OverrideFindOneDocumentNone(ctx context.Context, collection string, filter bson.D) (*bson.M, error) {
	return nil, fmt.Errorf(utils.ErrorMongoNoDocuments)
}

func OverrideFindOneDocumentRecipe(ctx context.Context, collection string, filter bson.D) (*bson.M, error) {
	var doc bson.M = bson.M{
		"ingredients": primitive.A{"hello"},
		"isCookable":  false,
	}
	return &doc, nil
}

func OverrideFindOneDocumentErrorBasic(ctx context.Context, collection string, filter bson.D) (*bson.M, error) {
	return nil, fmt.Errorf(errorBasic)
}

func OverrideFindOneDocumentErrorDecodeFail(ctx context.Context, collection string, filter bson.D) (*bson.M, error) {
	var doc bson.M = map[string]interface{}{"key": make(chan int)}
	return &doc, nil
}

func OverrideFindManyDocumentsSuccess(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	return []bson.M{
		map[string]interface{}{"name": "value1"},
		map[string]interface{}{"name": "value2", "attributes": map[string]string{}},
	}, nil
}

func OverrideFindManyDocumentsIngredient(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	// For isCookable
	//	and := filter["$and"].([]bson.M)
	//	expirationDate := and[0]["expirationDate"].(bson.M)
	//	value := expirationDate["$gt"].(int64)
	current := int64(time.Now().UTC().UnixNano()) / int64(time.Millisecond)
	//	if current >= value {
	return []bson.M{map[string]interface{}{"_id": 1337, "expirationDate": current, "haveStocked": "true", "name": "hello", "type": "thing"}}, nil
	//	} else {
	//		return []bson.M{}, nil
	//	}
}

func OverrideFindManyDocumentsIngredientRange(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	//	m1 := filter["$and"].([]bson.M)
	//	m2 := m1[0]["expirationDate"].(bson.M)
	//	low := m2["$gte"].(int64)
	//	high := m2["$lte"].(int64)
	expirationDate := int64(25)
	//	if expirationDate >= low && expirationDate <= high {
	return []bson.M{map[string]interface{}{"_id": 1337, "expirationDate": expirationDate, "haveStocked": "false", "name": "hello", "type": "thing"}}, nil
	//	} else {
	//		return []bson.M{}, nil
	//	}
}

func OverrideFindManyDocumentsRecipe(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	return []bson.M{map[string]interface{}{"isCookable": false, "ingredients": primitive.A{"hello"}}}, nil
}

func OverrideFindManyDocumentsSuper(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	if collection == config.MongoCollectionRecipes {
		return OverrideFindManyDocumentsRecipe(ctx, collection, filter, opts)
	} else {
		return OverrideFindManyDocumentsIngredient(ctx, collection, filter, opts)
	}
}

func OverrideFindManyDocumentsErrorBasic(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	return nil, fmt.Errorf(errorBasic)
}

func OverrideFindManyDocumentsRecipeOrErrorBasic(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	if collection == config.MongoCollectionRecipes {
		return OverrideFindManyDocumentsRecipe(ctx, collection, filter, opts)
	} else {
		return OverrideFindManyDocumentsErrorBasic(ctx, collection, filter, opts)
	}
}

func OverrideFindManyDocumentsDecodeFail(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	return []bson.M{map[string]interface{}{"key": make(chan int)}}, nil
}

func OverrideFindManyDocumentsCheckExpirations2(ctx context.Context, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	expectation := bson.M{"$and": []bson.M{
		{
			"expirationDate": bson.M{
				"$gte": int64(time.Now().UTC().UnixNano()) / int64(time.Millisecond),
				"$lte": int64(time.Now().Add(configuration.Lookahead).UTC().UnixNano()) / int64(time.Millisecond),
			},
		},
		{
			"haveStocked": bson.M{
				"$eq": true,
			},
		},
	}}
	e, _ := bson.Marshal(expectation)
	f, _ := bson.Marshal(filter)

	if bytes.Equal(f, e) {
		return nil, fmt.Errorf(errorBasic)
	} else {
		return []bson.M{
			map[string]interface{}{"name": "value1"},
			map[string]interface{}{"name": "value2", "attributes": map[string]string{}},
		}, nil
	}
}

func OverrideInsertManyDocumentsErrorBasic(ctx context.Context, collection string, docs []interface{}) error {
	return fmt.Errorf(errorBasic)
}

func OverrideUpdateOneDocumentEmptyUpdate(ctx context.Context, collection string, filter bson.D, update interface{}) (int64, int64, error) {
	return 0, 0, fmt.Errorf(errorEmptyUpdateInstructions)
}

func OverrideUpdateOneDocumentZero(ctx context.Context, collection string, filter bson.D, update interface{}) (int64, int64, error) {
	return 0, 0, nil
}

func OverrideUpdateOneDocumentErrorBasic(ctx context.Context, collection string, filter bson.D, update interface{}) (int64, int64, error) {
	return 0, 0, fmt.Errorf(errorBasic)
}

func OverrideDeleteOneDocumentNone(ctx context.Context, collection string, filter bson.D) error {
	return fmt.Errorf(utils.ErrorMongoNoDocuments)
}

func OverrideDeleteOneDocumentErrorBasic(ctx context.Context, collection string, filter bson.D) error {
	return fmt.Errorf(errorBasic)
}

func OverrideDeleteManyDocumentsZero(ctx context.Context, collection string, filter bson.M) (int64, error) {
	return 0, nil
}

func OverrideDeleteManyDocumentsErrorBasic(ctx context.Context, collection string, filter bson.M) (int64, error) {
	return 0, fmt.Errorf(errorBasic)
}

func OverrideGetShoppingListNil() (*trello.Card, error) {
	return nil, nil
}

func OverrideGetShoppingListErrorBasic() (*trello.Card, error) {
	return nil, fmt.Errorf(errorBasic)
}

func OverrideCreateShoppingListErrorBasic(dueDate *time.Time, applyLabels []string, listItems []string) (string, error) {
	return "", fmt.Errorf(errorBasic)
}

func OverrideAddToShoppingListErroBasic(itemNames []string) (string, error) {
	return "", fmt.Errorf(errorBasic)
}

func OverrideComposeMessageEmpty(quantity, quantityExpired int, url string) string {
	return ""
}

func OverrideSendMessageErrorBasic(phoneFrom, phoneTo, message string) (string, error) {
	return "", fmt.Errorf(errorBasic)
}
