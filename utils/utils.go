package utils

import (
	"strconv"
	"time"
)

func ParseDatetimeFromMongoID(id string) (time.Time, error) {
	parsed, err := strconv.ParseInt(id[0:8]+"", 16, 64)
	if err != nil {
		return time.Time{}, nil
	}
	timestamp := time.Unix(parsed, 0)
	return timestamp, nil
}
