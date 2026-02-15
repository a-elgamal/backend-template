// Package health allows introspecting the service version and health status
package health

import (
	"encoding/json"
)

// Health contains a summary of the service health
type Health struct {
	DBVersion      uint
	ServiceVersion string
	GitTag         string
	GitCommit      string
	BuildDate      string
	Status         Status
	StatusText     string
}

// Status a simple enum for the possible values of status
type Status uint8

const (
	// StatusOk means that the service or the component is functioning normally
	StatusOk Status = iota
	// StatusError means that the service or the component are not functioning
	StatusError = iota
)

func (s Status) String() string {
	names := [...]string{"ok", "error"}
	if int(s) > len(names) {
		return "invalid"
	}
	return names[s]
}

// MarshalJSON ensures that the Status is printed as a string not an int
func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
