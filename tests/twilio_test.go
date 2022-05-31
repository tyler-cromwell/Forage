package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tyler-cromwell/forage/clients"
	"github.com/tyler-cromwell/forage/tests/mocks"
)

func TestTwilioClient(t *testing.T) {
	t.Run("NewTwilioClientWrapper", func(t *testing.T) {
		client := clients.NewTwilioClientWrapper("", "", "", "")
		require.NotNil(t, client)
	})

	t.Run("InnerStruct", func(t *testing.T) {
		client := clients.NewTwilioClientWrapper("", "", "", "")
		require.NotNil(t, client)

		self := client.InnerStruct()
		require.NotNil(t, self)
	})

	t.Run("ComposeMessage", func(t *testing.T) {
		client := mocks.NewTwilioClientWrapper("", "", "", "")
		require.NotNil(t, client)

		cases := []struct {
			quantity        int
			quantityExpired int
			url             string
			want            string
		}{
			{1, 0, "http://nothing.com", fmt.Sprintf("%d item expiring soon and %d already expired! View shopping list: %s", 1, 0, "http://nothing.com")},
			{2, 0, "http://nothing.com", fmt.Sprintf("%d items expiring soon and %d already expired! View shopping list: %s", 2, 0, "http://nothing.com")},
			{0, 1, "http://nothing.com", fmt.Sprintf("%d item expired! View shopping list: %s", 1, "http://nothing.com")},
			{0, 2, "http://nothing.com", fmt.Sprintf("%d items expired! View shopping list: %s", 2, "http://nothing.com")},
			{0, 0, "http://nothing.com", ""},
		}

		for _, c := range cases {
			got := client.ComposeMessage(c.quantity, c.quantityExpired, c.url)
			if got != c.want {
				t.Errorf("ComposeMessage(%d, %d, \"%s\"), got (\"%s\"), want (\"%s\")", c.quantity, c.quantityExpired, c.url, got, c.want)
			}
		}
	})

	t.Run("SendMessage", func(t *testing.T) {
		client := mocks.NewTwilioClientWrapper("", "", "", "")
		require.NotNil(t, client)

		// Case: "Success"
		_, err := client.SendMessage("+11111111111", "+11111111112", "Hello")
		require.Nil(t, err)

		// Case: "Error"
		_, err = client.SendMessage("+11111111111", "+11111111112", "Error")
		require.NotNil(t, err)
	})
}
