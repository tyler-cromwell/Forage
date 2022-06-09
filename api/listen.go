//go:build !test

package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/config"
)

func ListenAndServe(ctx context.Context, c *config.Configuration) {
	configuration = c
	// Launch job to periodically check for expiring food
	ticker := time.NewTicker(configuration.Interval)
	quit := make(chan struct{})
	checkExpirations() // Run once before first ticker tick
	go func() {
		// Specify common fields
		log := logrus.WithFields(logrus.Fields{
			"at":        "expirationJob",
			"interval":  configuration.Interval,
			"lookahead": configuration.Lookahead,
		})

		// Wait for ticker ticks
		log.Info("Expiration watch job started")
		for {
			select {
			case <-ticker.C:
				checkExpirations()
			case <-quit:
				ticker.Stop()
				log.Info("Expiration watch job stopped")
				return
			}
		}
	}()

	// Define route actions/methods
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/documents/{id}", getOneDocument).Methods("GET")
	router.HandleFunc("/documents/{id}", putOneDocument).Methods("PUT")
	router.HandleFunc("/documents/{id}", deleteOneDocument).Methods("DELETE")
	router.HandleFunc("/documents", getManyDocuments).Methods("GET")
	router.HandleFunc("/documents", postManyDocuments).Methods("POST")
	router.HandleFunc("/documents", deleteManyDocuments).Methods("DELETE")
	router.HandleFunc("/expiring", getExpiring).Methods("GET")
	router.HandleFunc("/expired", getExpired).Methods("GET")

	// Specify common fields
	log := logrus.WithFields(logrus.Fields{"socket": configuration.ListenSocket})

	// Listen for HTTP requests
	log.Info("Listening for HTTP requests")
	err := http.ListenAndServe(configuration.ListenSocket, router)
	if err != nil {
		log.WithError(err).Fatal("Failed to listen for and serve requests")
	}
}
