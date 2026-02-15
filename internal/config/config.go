/*
Package config reads application configurations by default read from application.yaml and can be overriden using environment variables.
To override with an environment variable, replace any '.' in the configuration name with '_'
*/
package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// ReadConfig by default searches for application.yaml in current working directory or parent directories [up to 2 levels], then overrides the configuration read with relevant environment variables
func ReadConfig() (Config, *viper.Viper) {
	v := viper.New()
	v.SetConfigName("application")
	v.SetConfigType("yaml")
	v.AddConfigPath("./")
	v.AddConfigPath("../")
	v.AddConfigPath("../../")
	v.AddConfigPath("../../../")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error encountered while parsing server config: %v", err)
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	log.Printf("Using config file: %v\n", v.ConfigFileUsed())

	return NewConfigFromViper(v), v
}

// Config is the main struct that holds all the application configurations
type Config struct {
	ServerConfig    ServerConfig
	DBConfig        DBConfig
	TelemetryConfig TelemetryConfig
	GCPConfig       GCPConfig
	AWSConfig       AWSConfig
}

// NewConfigFromViper Creates a new Config struct from a Viper object
func NewConfigFromViper(v *viper.Viper) Config {
	return Config{
		ServerConfig:    ServerConfig{v},
		DBConfig:        DBConfig{v},
		TelemetryConfig: TelemetryConfig{v},
		GCPConfig:       GCPConfig{v},
		AWSConfig:       AWSConfig{v},
	}
}
