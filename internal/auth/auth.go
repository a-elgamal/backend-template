// Package auth defines the authentication provider interface used by the server
package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
)

// Provider is implemented by cloud-specific authentication middleware (e.g. GCP IAP, AWS ALB OIDC)
type Provider interface {
	// Middleware returns a Gin middleware that authenticates requests and sets user identity on the context
	Middleware(logger logr.Logger) gin.HandlerFunc
	// IsEnabled reports whether this auth provider is configured and should be used
	IsEnabled() bool
}
