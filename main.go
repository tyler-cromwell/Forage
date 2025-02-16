//go:build !test

package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/api"
	"github.com/tyler-cromwell/forage/clients"
	"github.com/tyler-cromwell/forage/config"
)

func main() {
	// Specify common fields
	log := logrus.WithFields(logrus.Fields{"at": "main.main"})

	// Configure logrus logging
	levelStr := os.Getenv("LOGRUS_LEVEL")
	if levelStr == "" {
		log.Fatal("Logging level not specified")
	}

	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		log.WithFields(logrus.Fields{"level": levelStr}).WithError(err).Fatal("Failed to parse logging level")
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(level)
	logrus.SetOutput(os.Stdout)
	//    logrus.SetReportCaller(true)
	log.WithFields(logrus.Fields{"level": levelStr}).Info("Logging configured")

	// Get other environment variables
	contextTimeoutStr := os.Getenv("FORAGE_CONTEXT_TIMEOUT")
	if contextTimeoutStr == "" {
		// Default case
		contextTimeoutStr = "5s"
		log.WithFields(logrus.Fields{"timeout": contextTimeoutStr}).Info("Using default context timeout")
	}

	forageContextTimeout, err := time.ParseDuration(contextTimeoutStr)
	if err != nil {
		log.WithFields(logrus.Fields{"timeout": contextTimeoutStr}).WithError(err).Fatal("Failed to parse context timeout")
	}

	intervalStr := os.Getenv("FORAGE_INTERVAL")
	if intervalStr == "" {
		// Default case
		intervalStr = "1" // Days
		log.WithFields(logrus.Fields{"interval": intervalStr}).Info("Using default expiration interval")
	}

	lookaheadStr := os.Getenv("FORAGE_LOOKAHEAD")
	if lookaheadStr == "" {
		// Default case
		lookaheadStr = "48h"
		log.WithFields(logrus.Fields{"lookahead": lookaheadStr}).Info("Using default expiration lookahead")
	}

	forageTime := os.Getenv("FORAGE_TIME")
	if forageTime == "" {
		// Default case
		forageTime = "19:00"
		log.WithFields(logrus.Fields{"timezone": forageTime}).Info("Using default time")
	}

	forageTimezone := os.Getenv("FORAGE_TIMEZONE")
	if forageTimezone == "" {
		// Default case
		forageTimezone = "America/New_York"
		log.WithFields(logrus.Fields{"timezone": forageTimezone}).Info("Using default timezone")
	}

	forageLookahead, err := time.ParseDuration(lookaheadStr)
	if err != nil {
		log.WithFields(logrus.Fields{"lookahead": lookaheadStr}).WithError(err).Fatal("Failed to parse expiration lookahead")
	}

	forageInterval, err := strconv.Atoi(intervalStr)
	if err != nil {
		log.WithFields(logrus.Fields{"interval": intervalStr}).WithError(err).Fatal("Failed to parse expiration interval")
	}

	mongoUri := os.Getenv("MONGO_URI")
	listenSocket := os.Getenv("LISTEN_SOCKET")
	trelloMemberID := os.Getenv("TRELLO_MEMBER")
	trelloBoardName := os.Getenv("TRELLO_BOARD")
	trelloListName := os.Getenv("TRELLO_LIST")
	trelloLabels := os.Getenv("TRELLO_LABELS")
	trelloApiKey := os.Getenv("TRELLO_API_KEY")
	trelloApiToken := os.Getenv("TRELLO_API_TOKEN")
	twilioAccountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhoneFrom := os.Getenv("TWILIO_PHONE_FROM")
	twilioPhoneTo := os.Getenv("TWILIO_PHONE_TO")

	// Initialize context/timeout
	ctx, cancel := context.WithTimeout(context.Background(), forageContextTimeout)
	log.WithFields(logrus.Fields{"timeout": forageContextTimeout}).Info("Initialized context")
	defer cancel()

	// Initialize clients
	mongoClient, err := clients.NewMongoClientWrapper(ctx, mongoUri)
	if err != nil {
		log.WithError(err).Fatal("Failed to create MongoDB client wrapper")
	} else {
		defer mongoClient.Client.Disconnect(ctx)
	}
	trelloClient := clients.NewTrelloClientWrapper(trelloApiKey, trelloApiToken, trelloMemberID, trelloBoardName, trelloListName, trelloLabels)
	twilioClient := clients.NewTwilioClientWrapper(twilioAccountSid, twilioAuthToken, twilioPhoneFrom, twilioPhoneTo)

	config := config.Configuration{
		ContextTimeout: forageContextTimeout,
		Lookahead:      forageLookahead,
		Interval:       forageInterval,
		Time:           forageTime,
		Timezone:       forageTimezone,
		LogrusLevel:    level,
		ListenSocket:   listenSocket,
		Mongo:          mongoClient,
		Trello:         trelloClient,
		Twilio:         twilioClient,
	}

	api.ListenAndServe(ctx, &config)
}
