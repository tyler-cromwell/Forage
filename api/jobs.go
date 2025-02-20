package api

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tyler-cromwell/forage/config"
	"go.mongodb.org/mongo-driver/bson"
)

func checkExpirations() {
	// Setup
	ctx := context.Background()
	log := logrus.WithFields(logrus.Fields{
		"at":        "api.checkExpirations",
		"lookahead": configuration.Lookahead,
	})

	// Log diagnostic information
	log.Trace("Begin function")
	defer log.Trace("End function")

	// Filter by food expired already
	now := int64(time.Now().UTC().UnixNano()) / int64(time.Millisecond)
	later := int64(time.Now().Add(configuration.Lookahead).UTC().UnixNano()) / int64(time.Millisecond)
	filterExpired := bson.M{"$and": []bson.M{
		{
			"expirationDate": bson.M{
				"$lte": now,
			},
		},
		{
			"haveStocked": bson.M{
				"$eq": true,
			},
		},
	}}
	log.WithFields(logrus.Fields{"type": "expired", "value": filterExpired}).Debug("Filter data")

	// Filter by food expiring within the given search window
	filterExpiring := bson.M{"$and": []bson.M{
		{
			"expirationDate": bson.M{
				"$gte": now,
				"$lte": later,
			},
		},
		{
			"haveStocked": bson.M{
				"$eq": true,
			},
		},
	}}
	log.WithFields(logrus.Fields{"type": "expiring", "value": filterExpiring}).Debug("Filter data")

	// Grab the documents
	documentsExpired, err := configuration.Mongo.FindDocuments(ctx, config.MongoCollectionIngredients, filterExpired, nil)
	if err != nil {
		log.WithError(err).Error("Failed to identify expired items")
		return
	}

	documentsExpiring, err := configuration.Mongo.FindDocuments(ctx, config.MongoCollectionIngredients, filterExpiring, nil)
	if err != nil {
		log.WithError(err).Error("Failed to identify expiring items")
	} else {
		quantityExpired := len(documentsExpired)
		quantityExpiring := len(documentsExpiring)

		// Skip if nothing is expiring
		if quantityExpiring == 0 && quantityExpired == 0 {
			log.WithFields(logrus.Fields{"expiring": quantityExpiring, "expired": quantityExpired}).Info("Restocking not required")
			return
		} else {
			log.WithFields(logrus.Fields{"expiring": quantityExpiring, "expired": quantityExpired}).Info("Restocking required")
		}

		log.WithFields(logrus.Fields{"quantity": quantityExpired, "value": documentsExpired}).Debug("Expired items")
		log.WithFields(logrus.Fields{"quantity": quantityExpiring, "value": documentsExpiring}).Debug("Expiring items")

		// Construct list of names of items to shop for
		var groceries []string
		for _, document := range documentsExpired {
			name := document["name"]
			stage := "expired"
			text := fmt.Sprintf("%s (%s)", name, stage)

			if _, ok := document["attributes"]; ok {
				brand := document["attributes"].(map[string]string)["brand"]
				flavor := document["attributes"].(map[string]string)["flavor"]
				text = fmt.Sprintf("%s (%s, %s, %s)", name, brand, flavor, stage)
			}

			groceries = append(groceries, text)
		}
		for _, document := range documentsExpiring {
			name := document["name"]
			stage := "expiring"
			text := fmt.Sprintf("%s (%s)", name, stage)

			if _, ok := document["attributes"]; ok {
				brand := document["attributes"].(map[string]string)["brand"]
				flavor := document["attributes"].(map[string]string)["flavor"]
				text = fmt.Sprintf("%s (%s, %s, %s)", name, brand, flavor, stage)
			}

			groceries = append(groceries, text)
		}
		log.WithFields(logrus.Fields{"quantity": len(groceries), "value": groceries}).Debug("Groceries")

		// Construct shopping list due date
		now := time.Now()
		rounded := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		dueDate := rounded.Add(configuration.Lookahead + (time.Hour * 24))
		log.WithFields(logrus.Fields{"value": dueDate}).Debug("Card due date")

		var url string
		shoppingListCard, err := configuration.Trello.GetShoppingList()
		if err != nil {
			log.WithError(err).Error("Failed to get Trello card")
		} else if shoppingListCard != nil {
			// Add to shopping list card on Trello
			url, err = configuration.Trello.AddToShoppingList(groceries)
			if err != nil {
				log.WithError(err).Error("Failed to add to Trello card")
			} else {
				log.WithFields(logrus.Fields{"url": url}).Info("Added to Trello card")
			}
		} else {
			// Create shopping list card on Trello
			innerTrello := reflect.ValueOf(configuration.Trello).Elem()
			labelsStr := *innerTrello.FieldByName("LabelsStr").Addr().Interface().(*string)
			labels := strings.Split(labelsStr, ",")
			url, err = configuration.Trello.CreateShoppingList(&dueDate, labels, groceries)
			if err != nil {
				log.WithError(err).Error("Failed to create Trello card")
			} else {
				log.WithFields(logrus.Fields{"url": url}).Info("Created Trello card")
			}
		}

		// Compose Twilio message
		var message = configuration.Twilio.ComposeMessage(quantityExpiring, quantityExpired, url)

		// Send the Twilio message
		if !configuration.Silence {
			innerTwilio := reflect.ValueOf(configuration.Twilio).Elem()
			from := *innerTwilio.FieldByName("From").Addr().Interface().(*string)
			to := *innerTwilio.FieldByName("To").Addr().Interface().(*string)
			_, err = configuration.Twilio.SendMessage(from, to, message)
			if err != nil {
				log.WithFields(logrus.Fields{"from": from, "to": to}).WithError(err).Error("Failed to send Twilio message")
			} else {
				log.WithFields(logrus.Fields{"from": from, "to": to}).Info("Sent Twilio message")
			}
		} else {
			log.WithFields(logrus.Fields{"silence": configuration.Silence}).Info("Skipped Twilio message")
		}
	}
}
