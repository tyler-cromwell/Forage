package utils

import (
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ErrNoMatchedDocuments = "no document matching filter"
const ErrInvalidObjectID = "the provided hex string is not a valid ObjectID"
const ErrMongoNoDocuments = "mongo: no documents in result"

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

func ParseDatetimeFromMongoID(id string) (time.Time, error) {
	parsed, err := strconv.ParseInt(id[0:8]+"", 16, 64)
	if err != nil {
		return time.Time{}, nil
	}
	timestamp := time.Unix(parsed, 0)
	return timestamp, nil
}
