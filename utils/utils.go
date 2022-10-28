package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ErrorNoMatchedDocuments = "no document matching filter"
const ErrorInvalidObjectID = "the provided hex string is not a valid ObjectID"
const ErrorMongoNoDocuments = "mongo: no documents in result"

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ParseDatetimeFromMongoID(id string) (time.Time, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	return oid.Timestamp(), err
}

func StringSliceFromBsonM(documents []primitive.M, key string) []string {
	var slice []string
	for _, document := range documents {
		value, keyFound := document[key]
		if keyFound {
			slice = append(slice, value.(string))
		}
	}
	return slice
}
