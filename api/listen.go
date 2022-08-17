//go:build !test

package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/config"
)

func ListenAndServe(ctx context.Context, c *config.Configuration) {
	configuration = c

	// Specify common fields
	log := logrus.WithFields(logrus.Fields{"at": "listen.ListenAndServe"})

	// Log diagnostic information
	log.Trace("Begin function")
	defer log.Trace("End function")

	// Load timezone
	loc, err := time.LoadLocation(configuration.Timezone)
	if err != nil {
		log.WithFields(logrus.Fields{"timezone": configuration.Timezone}).WithError(err).Fatal("Failed to obtain timezone")
	} else {
		log.WithFields(logrus.Fields{"timezone": configuration.Timezone}).Info("Parsed timezone")
	}

	// Launch job to periodically check for expiring food
	s := gocron.NewScheduler(loc)
	s.Every(1).Day().At(configuration.Time).Do(checkExpirations)
	s.StartAsync()
	log.Info("Expiration watch job started")

	// Define route actions/methods
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/configure", getConfiguration).Methods("GET")
	router.HandleFunc("/configure", putConfiguration).Methods("PUT")
	router.HandleFunc("/documents/{id}", getOneDocument).Methods("GET")
	router.HandleFunc("/documents/{id}", putOneDocument).Methods("PUT")
	router.HandleFunc("/documents/{id}", deleteOneDocument).Methods("DELETE")
	router.HandleFunc("/documents", getManyDocuments).Methods("GET")
	router.HandleFunc("/documents", postManyDocuments).Methods("POST")
	router.HandleFunc("/documents", deleteManyDocuments).Methods("DELETE")
	router.HandleFunc("/expiring", getExpiring).Methods("GET")
	router.HandleFunc("/expired", getExpired).Methods("GET")

	// Specify common fields
	log = log.WithFields(logrus.Fields{"socket": configuration.ListenSocket})

	// Listen for HTTP requests
	log.Info("Listening for HTTP requests")
	err = http.ListenAndServe(configuration.ListenSocket, router)
	if err != nil {
		log.WithError(err).Fatal("Failed to listen for and serve requests")
	}
}
