package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDBConfig_GetURL(t *testing.T) {
	t.Run("Returns DB.URL if set", func(t *testing.T) {
		url := "url"
		v := viper.New()
		v.Set(dbConfigURL, url)
		v.Set(dbConfigURLTemplate, "template %v")
		v.Set(dbConfigName, "name")

		dbConfig := DBConfig{v}
		assert.Equal(t, url, dbConfig.GetURL())
	})

	t.Run("Returns Create DB URL from template if no DB.URL", func(t *testing.T) {
		v := viper.New()
		v.Set(dbConfigURLTemplate, "template %v")
		v.Set(dbConfigName, "name")

		dbConfig := DBConfig{v}
		assert.Equal(t, "template name", dbConfig.GetURL())
	})
}
