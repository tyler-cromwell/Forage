package api

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
	"github.com/tyler-cromwell/forage/tests/mocks"
)

func TestJobs(t *testing.T) {
	// Discard logging output
	logrus.SetOutput(io.Discard)

	subtests2 := []struct {
		name         string
		mongoClient  mocks.MockMongo
		mocksClient  mocks.MockTrello
		twilioClient mocks.MockTwilio
		logLevels    []logrus.Level
		logMessages  []string
	}{
		{
			// Error #1, Could not obtain expired items
			"checkExpirationsError#1",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsErrorBasic},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.ErrorLevel},
			[]string{"Failed to identify expired items"},
		},
		{
			// Error #2, Could not obtain expiring items
			"checkExpirationsError#2",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsCheckExpirations2},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.ErrorLevel},
			[]string{"Failed to identify expiring items"},
		},
		{
			// Success #1, No expired/expiring items, no need to proceed.
			"checkExpirationsSuccess#1",
			mocks.MockMongo{},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel},
			[]string{"Restocking not required"},
		},
		{
			// Success #2, items expired/expiring added to existing Trello card and SMS message sent.
			"checkExpirationsSuccess#2",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.InfoLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Added to Trello card", "Sent Twilio message"},
		},
		{
			// Error #3, items expired/expiring but could not obtain Trello card, SMS message still sent.
			"checkExpirationsError#3",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideGetShoppingList: OverrideGetShoppingListErrorBasic},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Failed to get Trello card", "Sent Twilio message"},
		},
		{
			// Success #3, items expired/expiring added to new Trello card and SMS message sent.
			"checkExpirationsSuccess#3",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideGetShoppingList: OverrideGetShoppingListNil},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.InfoLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Created Trello card", "Sent Twilio message"},
		},
		{
			// Error #4, items expired/expiring but could not add to existing card Trello card, SMS message still sent.
			"checkExpirationsError#4",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideAddToShoppingList: OverrideAddToShoppingListErroBasic},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Failed to add to Trello card", "Sent Twilio message"},
		},
		{
			// Error #5, items expired/expiring but could not create new card Trello card, SMS message still sent.
			"checkExpirationsError#5",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideGetShoppingList: OverrideGetShoppingListNil, OverrideCreateShoppingList: OverrideCreateShoppingListErrorBasic},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Failed to create Trello card", "Sent Twilio message"},
		},
		{
			// Error #6, items expired/expiring but could not create new card Trello card or send SMS message.
			"checkExpirationsError#6",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideGetShoppingList: OverrideGetShoppingListNil, OverrideCreateShoppingList: OverrideCreateShoppingListErrorBasic},
			mocks.MockTwilio{OverrideComposeMessage: OverrideComposeMessageEmpty, OverrideSendMessage: OverrideSendMessageErrorBasic},
			[]logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.ErrorLevel},
			[]string{"Restocking required", "Failed to create Trello card", "Failed to send Twilio message"},
		},
		{
			// Success #4, items expired/expiring added to new Trello card and SMS message skipped.
			"checkExpirationsSuccess#4",
			mocks.MockMongo{OverrideFindManyDocuments: OverrideFindManyDocumentsSuccess},
			mocks.MockTrello{OverrideGetShoppingList: OverrideGetShoppingListNil},
			mocks.MockTwilio{},
			[]logrus.Level{logrus.InfoLevel, logrus.InfoLevel, logrus.InfoLevel},
			[]string{"Restocking required", "Created Trello card", "Skipped Twilio message"},
		},
	}

	t.Run("checkExpirations", func(t *testing.T) {
		// Capture logrus output so we can assert
		_, hook := test.NewNullLogger()
		logrus.AddHook(hook)
		base := 0

		for _, st := range subtests2 {
			// Arrange
			configuration.Mongo = &st.mongoClient
			configuration.Trello = &st.mocksClient
			configuration.Twilio = &st.twilioClient

			if st.name == "checkExpirationsSuccess#4" {
				configuration.Silence = true
			} else {
				configuration.Silence = false
			}

			// Act
			checkExpirations()

			// Assert (preliminary)
			require.Equal(t, len(st.logLevels), len(st.logMessages))

			// Assert (primary)
			for i, _ := range st.logLevels {
				index := base + i
				require.Equal(t, st.logLevels[i], hook.AllEntries()[index].Level)
				require.Equal(t, st.logMessages[i], hook.AllEntries()[index].Message)
			}

			base += len(st.logLevels)
		}

		// Rever logrus output change
		logrus.SetOutput(io.Discard)
	})
}
