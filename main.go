//go:build !test

package main

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/api"
	"github.com/tyler-cromwell/forage/clients"
	"github.com/tyler-cromwell/forage/config"
)

func main() {
	// Configure logrus logging
	levelStr := os.Getenv("LOGRUS_LEVEL")
	if levelStr == "" {
		logrus.Fatal("Logging level not specified")
	}

	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{"level": levelStr}).WithError(err).Fatal("Failed to parse logging level")
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(level)
	logrus.SetOutput(os.Stdout)
	//    logrus.SetReportCaller(true)
	logrus.WithFields(logrus.Fields{"level": levelStr}).Info("Logging configured")

	// Get other environment variables
	contextTimeoutStr := os.Getenv("FORAGE_CONTEXT_TIMEOUT")
	if contextTimeoutStr == "" {
		// Default case
		contextTimeoutStr = "5s"
		logrus.WithFields(logrus.Fields{"timeout": contextTimeoutStr}).Debug("Setting context timeout to default")
	}

	forageContextTimeout, err := time.ParseDuration(contextTimeoutStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{"timeout": contextTimeoutStr}).WithError(err).Fatal("Failed to parse context timeout")
	}

	intervalStr := os.Getenv("FORAGE_INTERVAL")
	if intervalStr == "" {
		// Default case
		intervalStr = "24h"
		logrus.WithFields(logrus.Fields{"interval": intervalStr}).Debug("Setting expiration interval to default")
	}

	forageInterval, err := time.ParseDuration(intervalStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{"interval": intervalStr}).WithError(err).Fatal("Failed to parse expiration interval")
	}

	lookaheadStr := os.Getenv("FORAGE_LOOKAHEAD")
	if lookaheadStr == "" {
		// Default case
		lookaheadStr = "48h"
		logrus.WithFields(logrus.Fields{"lookahead": lookaheadStr}).Debug("Setting expiration lookahead to default")
	}

	forageLookahead, err := time.ParseDuration(lookaheadStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{"lookahead": lookaheadStr}).WithError(err).Fatal("Failed to parse expiration lookahead")
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
	logrus.WithFields(logrus.Fields{"timeout": forageContextTimeout}).Info("Initialized context")
	defer cancel()

	// Initialize clients
	mongoClient, err := clients.NewMongoClientWrapper(ctx, mongoUri)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create MongoDB client wrapper")
	} else {
		defer mongoClient.Client.Disconnect(ctx)
	}
	trelloClient := clients.NewTrelloClientWrapper(trelloApiKey, trelloApiToken, trelloMemberID, trelloBoardName, trelloListName, trelloLabels)
	twilioClient := clients.NewTwilioClientWrapper(twilioAccountSid, twilioAuthToken, twilioPhoneFrom, twilioPhoneTo)

	config := config.Configuration{
		ContextTimeout: forageContextTimeout,
		Interval:       forageInterval,
		Lookahead:      forageLookahead,
		LogrusLevel:    level,
		ListenSocket:   listenSocket,
		Mongo:          mongoClient,
		Trello:         trelloClient,
		Twilio:         twilioClient,
	}

	api.ListenAndServe(ctx, &config)
}
