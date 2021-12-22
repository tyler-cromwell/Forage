package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tyler-cromwell/forage/clients"
)

func TestTwilioClient(t *testing.T) {
	t.Run("NewTwilioClientWrapper", func(t *testing.T) {
		// Case: "Success"
		client := clients.NewTwilioClientWrapper("", "", "", "")
		require.NotNil(t, client)
	})
}
