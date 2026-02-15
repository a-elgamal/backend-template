// Package item provides a sample domain entity for the service template
package item

import (
	"go.opentelemetry.io/otel"
)

const itemTableName = "item"

var tracer = otel.Tracer("myservice.item")

const nameJSONKey = "name"

// Item represents a sample domain entity
type Item struct {
	// Name The name of the item
	Name string `json:"name"`

	// Description A description of the item
	Description string `json:"description"`
}
