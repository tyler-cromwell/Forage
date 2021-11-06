package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/api"
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

	// Decide what port to listen on
	tcpSocket := os.Getenv("LISTEN_SOCKET")
	api.ListenAndServe(tcpSocket)
}
