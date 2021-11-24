package config

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/clients"
)

type Configuration struct {
	ContextTimeout time.Duration
	Interval       time.Duration
	Lookahead      time.Duration
	LogrusLevel    logrus.Level
	ListenSocket   string
	Mongo          *clients.Mongo
	Trello         *clients.Trello
	Twilio         *clients.Twilio
}
