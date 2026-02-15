// Package app provides the app domain entity for the service
package app

import (
	"go.opentelemetry.io/otel"
)

const appTableName = "app"

var tracer = otel.Tracer("myservice.app")

const disabledJSONKey = "disabled"

// App represents an application entity
type App struct {
	// APIKey The API key for the app
	APIKey string `json:"apiKey"`

	// Disabled Whether the app is disabled
	Disabled bool `json:"disabled"`
}
