// Package google contains GCP-specific code
package google

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"google.golang.org/api/idtoken"

	"alielgamal.com/myservice/internal"
	"alielgamal.com/myservice/internal/config"
	"alielgamal.com/myservice/internal/response"
)

var tracer = otel.Tracer("myservice/google")

// AuthProvider implements auth.Provider for GCP IAP authentication
type AuthProvider struct {
	conf config.GCPConfig
}

// NewAuthProvider creates a new GCP IAP auth provider
func NewAuthProvider(conf config.GCPConfig) *AuthProvider {
	return &AuthProvider{conf: conf}
}

// Middleware returns a Gin middleware that validates GCP IAP signed headers
func (p *AuthProvider) Middleware(logger logr.Logger) gin.HandlerFunc {
	return authMiddleware(logger, p.conf.ProjectNumber(), p.conf.Region(), p.conf.InternalBackendServiceID())
}

// IsEnabled reports whether IAP authentication is configured
func (p *AuthProvider) IsEnabled() bool {
	return p.conf.IAPAuthEnabled()
}

// AuthMiddleware Returns an Gin Authentication middleware that authenticates requests by checking IAP signed authentication headers. Upon successful authentication, the Gin context is set with the user's id & email
func AuthMiddleware(logger logr.Logger, conf config.GCPConfig) gin.HandlerFunc {
	return authMiddleware(logger, conf.ProjectNumber(), conf.Region(), conf.InternalBackendServiceID())
}

func authMiddleware(logger logr.Logger, projectNumber int64, region string, backendServiceID int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "google.middleware")
		defer span.End()

		aud := fmt.Sprintf("/projects/%v/%v/backendServices/%v", projectNumber, region, backendServiceID)

		iapJWTs, ok := c.Request.Header[http.CanonicalHeaderKey("x-goog-iap-jwt-assertion")]
		if !ok {
			// Reject the request as the IAP JWT token is missing.
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Err: response.ErrorDetail{
					Code: http.StatusUnauthorized,
					Msg:  "Missing Google IAP JWT Header",
				},
			})
			return
		}
		// Only use the first value of the header
		iapJWT := iapJWTs[0]

		payload, err := idtoken.Validate(ctx, iapJWT, aud)
		if err != nil {
			logger.Error(err, "Unable to validate JWT Token", "token", iapJWT, "aud", aud)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Err: response.ErrorDetail{
					Code: http.StatusUnauthorized,
					Msg:  fmt.Sprintf("Invalid IAP JWT Header: %v", err.Error()),
				},
			})
			return
		}

		c.Set(internal.UserIDContextKey, payload.Subject[strings.LastIndex(payload.Subject, ":"):])

		switch emailClaim := payload.Claims["email"].(type) {
		case string:
			c.Set(internal.UserEmailContextKey, emailClaim)
		default:
			logger.Info("Cannot parse email claim from JWT ID Token", "emailClaim", emailClaim, "payload", payload)
		}

		c.Next()
	}
}
