package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/api"
)

func main() {
	// Configure logrus logging
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	//    logrus.SetReportCaller(true)

	// Decide what port to listen on
	tcpSocket := os.Getenv("LISTEN_SOCKET")
	api.ListenAndServe(tcpSocket)
}
