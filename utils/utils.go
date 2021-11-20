package utils

import (
	"strconv"
	"time"
)

const ErrInvalidObjectID = "the provided hex string is not a valid ObjectID"
const ErrMongoNoDocuments = "mongo: no documents in result"

func ParseDatetimeFromMongoID(id string) (time.Time, error) {
	parsed, err := strconv.ParseInt(id[0:8]+"", 16, 64)
	if err != nil {
		return time.Time{}, nil
	}
	timestamp := time.Unix(parsed, 0)
	return timestamp, nil
}
