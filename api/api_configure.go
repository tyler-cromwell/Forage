package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/config"
)

var configuration *config.Configuration

func getConfiguration(response http.ResponseWriter, request *http.Request) {
	// Setup
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.getConfiguration",
		"method": "GET",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Prepare the response data
	marshalled, _ := json.Marshal(struct {
		Lookahead time.Duration `json:"lookahead"`
		Silence   bool          `json:"silence"`
		Time      string        `json:"time"`
	}{
		configuration.Lookahead,
		configuration.Silence,
		configuration.Time,
	})

	// Log & Respond
	log.WithFields(logrus.Fields{"size": len(marshalled), "state": "marshalled", "value": string(marshalled)}).Debug("Response body")
	log.WithFields(logrus.Fields{"status": http.StatusOK}).Info("Succeeded")
	response.WriteHeader(http.StatusOK)
	response.Write(marshalled)
}

func putConfiguration(response http.ResponseWriter, request *http.Request) {
	// Setup
	log := logrus.WithFields(logrus.Fields{
		"at":     "api.putConfiguration",
		"method": "PUT",
	})

	// Log diagnostic information
	log.Trace("Begin function")
	log.WithFields(logrus.Fields{"value": request}).Debug("Request data")
	defer log.Trace("End function")

	// Read in request body
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to read request body")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"size": len(bytes), "state": "marshalled", "value": string(bytes)}).Debug("Request body")
	}

	// Parse request body
	var body struct {
		Lookahead time.Duration `json:"lookahead"`
		Silence   bool          `json:"silence"`
		Time      string        `json:"time"`
	}
	err = json.Unmarshal(bytes, &body)
	if err != nil && strings.HasPrefix(err.Error(), "invalid character") {
		// Invalid request body
		log.WithFields(logrus.Fields{"case": 1}).Trace("Invalid request body")
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Warn("Failed to decode update fields")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	} else if err != nil && err.Error() == "unexpected end of JSON input" {
		// Invalid request body
		log.WithFields(logrus.Fields{"case": 2}).Trace("Invalid request body")
		log.WithFields(logrus.Fields{"status": http.StatusBadRequest}).WithError(err).Warn("Failed to decode update fields")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	} else if err != nil {
		// Something else failed
		log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to decode update fields")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	} else {
		log.WithFields(logrus.Fields{"state": "unmarshalled", "value": body}).Debug("Request body")
		log.WithFields(logrus.Fields{"lookahead": body.Lookahead, "silence": body.Silence, "time": body.Time}).Debug("Parsed data")

		if body.Time != "" && len(strings.Split(body.Time, ":")) != 2 {
			log.WithFields(logrus.Fields{"status": http.StatusBadRequest, "value": body.Time}).Warn("Invalid time format")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("Invalid time format: " + body.Time))
			return
		}

		if len(body.Time) != 0 && configuration.Time != body.Time {
			// Re-schedule expiration job
			configuration.Scheduler.Clear()
			log.Info("Scheduler cleared")
			_, err = configuration.Scheduler.Every(1).Day().At(body.Time).Do(checkExpirations)
			if err != nil {
				log.WithFields(logrus.Fields{"status": http.StatusInternalServerError}).WithError(err).Error("Failed to schedule expriation watch job")
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte(err.Error()))
				return
			} else {
				configuration.Scheduler.StartAsync()
				log.Info("Expiration watch job scheduled")
			}
		}

		configuration.Lookahead = body.Lookahead
		configuration.Silence = body.Silence
		configuration.Time = body.Time

		log.WithFields(logrus.Fields{"status": http.StatusOK}).Info("Succeeded")
		response.WriteHeader(http.StatusOK)
	}
}
