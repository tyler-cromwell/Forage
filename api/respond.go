package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func RespondWithError(response http.ResponseWriter, log *logrus.Entry, status int, message string) {
	er := ErrorResponse{
		Status:  status,
		Message: message,
	}

	mr, err := json.Marshal(er)
	if err != nil {
		log.WithError(err).Error("Failed to encode response")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Internal Server Error"))
	} else {
		response.WriteHeader(er.Status)
		response.Write(mr)
	}
}
