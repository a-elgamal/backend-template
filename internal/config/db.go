package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// DBConfig contains the DB related configuration
type DBConfig struct {
	v *viper.Viper
}

const dbConfigURL = "DB.URL"
const dbConfigURLTemplate = "DB.URL_TEMPLATE"
const dbConfigName = "DB.NAME"

// GetURL the actual DB connection URL that should be used. It automatically replaces the %v in the URL_Template with the Name.
func (c DBConfig) GetURL() string {
	if c.v.IsSet(dbConfigURL) {
		return c.v.GetString(dbConfigURL)
	}
	return fmt.Sprintf(c.GetURLTemplate(), c.v.GetString(dbConfigName))
}

// GetURLTemplate a template to use to generate the DB URL. The template is expected to have a single argument (DB name)
func (c DBConfig) GetURLTemplate() string {
	return c.v.GetString(dbConfigURLTemplate)
}
