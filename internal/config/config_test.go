package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	t.Run("Overrides Config with Environment Variable", func(t *testing.T) {
		const newAddress = ":69875"

		err := os.Setenv("SERVER_HTTP_ADDRESS", ":69875")
		assert.NoError(t, err)

		c, _ := ReadConfig()
		assert.Equal(t, newAddress, c.ServerConfig.GetHTTPAddress())
	})
}
