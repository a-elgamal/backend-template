// Package aws contains AWS-specific code
package aws

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel"

	"alielgamal.com/myservice/internal"
	"alielgamal.com/myservice/internal/config"
	"alielgamal.com/myservice/internal/response"
)

var tracer = otel.Tracer("myservice/aws")

// AuthProvider implements auth.Provider for AWS ALB OIDC authentication
type AuthProvider struct {
	conf config.AWSConfig
}

// NewAuthProvider creates a new AWS ALB auth provider
func NewAuthProvider(conf config.AWSConfig) *AuthProvider {
	return &AuthProvider{conf: conf}
}

// Middleware returns a Gin middleware that validates ALB-signed OIDC JWTs
func (p *AuthProvider) Middleware(logger logr.Logger) gin.HandlerFunc {
	return albAuthMiddleware(logger, p.conf.ALBRegion(), &httpKeyFetcher{})
}

// IsEnabled reports whether ALB OIDC authentication is configured
func (p *AuthProvider) IsEnabled() bool {
	return p.conf.ALBAuthEnabled()
}

// keyFetcher abstracts fetching ALB public keys for testing
type keyFetcher interface {
	FetchPublicKey(region, keyID string) (*ecdsa.PublicKey, error)
}

// httpKeyFetcher fetches ALB public keys over HTTPS
type httpKeyFetcher struct {
	mu    sync.RWMutex
	cache map[string]*ecdsa.PublicKey
}

func (f *httpKeyFetcher) FetchPublicKey(region, keyID string) (*ecdsa.PublicKey, error) {
	f.mu.RLock()
	if f.cache != nil {
		if key, ok := f.cache[keyID]; ok {
			f.mu.RUnlock()
			return key, nil
		}
	}
	f.mu.RUnlock()

	url := fmt.Sprintf("https://public-keys.auth.elb.%s.amazonaws.com/%s", region, keyID)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ALB public key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ALB public key endpoint returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read ALB public key response: %w", err)
	}

	block, _ := pem.Decode(body)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from ALB public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ALB public key: %w", err)
	}

	ecKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("ALB public key is not ECDSA")
	}

	f.mu.Lock()
	if f.cache == nil {
		f.cache = make(map[string]*ecdsa.PublicKey)
	}
	f.cache[keyID] = ecKey
	f.mu.Unlock()

	return ecKey, nil
}

func albAuthMiddleware(logger logr.Logger, region string, fetcher keyFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracer.Start(c.Request.Context(), "aws.middleware")
		defer span.End()

		oidcData := c.GetHeader("x-amzn-oidc-data")
		if oidcData == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Err: response.ErrorDetail{
					Code: http.StatusUnauthorized,
					Msg:  "Missing ALB OIDC data header",
				},
			})
			return
		}

		// Parse without verification first to get the key ID from the header
		unverified, _, err := jwt.NewParser().ParseUnverified(oidcData, jwt.MapClaims{})
		if err != nil {
			logger.Error(err, "Unable to parse ALB JWT header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Err: response.ErrorDetail{
					Code: http.StatusUnauthorized,
					Msg:  fmt.Sprintf("Invalid ALB OIDC JWT: %v", err.Error()),
				},
			})
			return
		}

		kid, ok := unverified.Header["kid"].(string)
		if !ok || kid == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Err: response.ErrorDetail{
					Code: http.StatusUnauthorized,
					Msg:  "ALB JWT missing kid header",
				},
			})
			return
		}

		pubKey, err := fetcher.FetchPublicKey(region, kid)
		if err != nil {
			logger.Error(err, "Unable to fetch ALB public key", "kid", kid, "region", region)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Err: response.ErrorDetail{
					Code: http.StatusUnauthorized,
					Msg:  fmt.Sprintf("Failed to fetch ALB public key: %v", err.Error()),
				},
			})
			return
		}

		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(oidcData, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return pubKey, nil
		})
		if err != nil {
			logger.Error(err, "Unable to validate ALB JWT", "kid", kid)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Err: response.ErrorDetail{
					Code: http.StatusUnauthorized,
					Msg:  fmt.Sprintf("Invalid ALB OIDC JWT: %v", err.Error()),
				},
			})
			return
		}

		if sub, ok := claims["sub"].(string); ok {
			c.Set(internal.UserIDContextKey, sub)
		}

		if email, ok := claims["email"].(string); ok {
			c.Set(internal.UserEmailContextKey, email)
		} else {
			logger.Info("Cannot parse email claim from ALB JWT", "claims", claims)
		}

		c.Next()
	}
}
