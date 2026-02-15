package aws

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/zerologr"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"alielgamal.com/myservice/internal"
	testDB "alielgamal.com/myservice/internal/db/test"
	"alielgamal.com/myservice/internal/health"
	"alielgamal.com/myservice/internal/response"
)

type mockKeyFetcher struct {
	key *ecdsa.PublicKey
	err error
}

func (m *mockKeyFetcher) FetchPublicKey(_, _ string) (*ecdsa.PublicKey, error) {
	return m.key, m.err
}

func signedALBToken(t *testing.T, key *ecdsa.PrivateKey, kid string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = kid
	signed, err := token.SignedString(key)
	assert.NoError(t, err)
	return signed
}

func TestALBAuthMiddleware(t *testing.T) {
	zl := zerolog.New(zerolog.NewConsoleWriter())
	logger := zerologr.New(&zl)

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	fetcher := &mockKeyFetcher{key: &privateKey.PublicKey}

	t.Run("Returns 401 with header not present", func(t *testing.T) {
		r := gin.Default()
		r.Use(albAuthMiddleware(logger, "us-east-1", fetcher))
		health.SetupRoutes(r, testDB.NewDB(t), 0)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+health.RouteRelativePath, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		result := response.ErrorResponse{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, result.Err.Code)
	})

	t.Run("Returns 401 with invalid JWT", func(t *testing.T) {
		r := gin.Default()
		r.Use(albAuthMiddleware(logger, "us-east-1", fetcher))
		health.SetupRoutes(r, testDB.NewDB(t), 0)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+health.RouteRelativePath, nil)
		req.Header.Set("x-amzn-oidc-data", "not-a-jwt")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Returns 401 when key fetch fails", func(t *testing.T) {
		failFetcher := &mockKeyFetcher{err: fmt.Errorf("network error")}
		r := gin.Default()
		r.Use(albAuthMiddleware(logger, "us-east-1", failFetcher))
		health.SetupRoutes(r, testDB.NewDB(t), 0)

		token := signedALBToken(t, privateKey, "test-kid", jwt.MapClaims{
			"sub":   "user123",
			"email": "test@example.com",
			"exp":   time.Now().Add(time.Hour).Unix(),
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+health.RouteRelativePath, nil)
		req.Header.Set("x-amzn-oidc-data", token)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Returns 401 with expired JWT", func(t *testing.T) {
		r := gin.Default()
		r.Use(albAuthMiddleware(logger, "us-east-1", fetcher))
		health.SetupRoutes(r, testDB.NewDB(t), 0)

		token := signedALBToken(t, privateKey, "test-kid", jwt.MapClaims{
			"sub":   "user123",
			"email": "test@example.com",
			"exp":   time.Now().Add(-time.Hour).Unix(),
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+health.RouteRelativePath, nil)
		req.Header.Set("x-amzn-oidc-data", token)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Succeeds with valid JWT and sets context", func(t *testing.T) {
		r := gin.Default()
		r.Use(albAuthMiddleware(logger, "us-east-1", fetcher))

		var capturedUserID, capturedEmail string
		r.GET("/test", func(c *gin.Context) {
			if v, ok := c.Get(internal.UserIDContextKey); ok {
				capturedUserID = v.(string)
			}
			if v, ok := c.Get(internal.UserEmailContextKey); ok {
				capturedEmail = v.(string)
			}
			c.Status(http.StatusOK)
		})

		token := signedALBToken(t, privateKey, "test-kid", jwt.MapClaims{
			"sub":   "user123",
			"email": "test@example.com",
			"exp":   time.Now().Add(time.Hour).Unix(),
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("x-amzn-oidc-data", token)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "user123", capturedUserID)
		assert.Equal(t, "test@example.com", capturedEmail)
	})

	t.Run("Returns 401 with wrong signing key", func(t *testing.T) {
		wrongKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		r := gin.Default()
		r.Use(albAuthMiddleware(logger, "us-east-1", fetcher))
		health.SetupRoutes(r, testDB.NewDB(t), 0)

		token := signedALBToken(t, wrongKey, "test-kid", jwt.MapClaims{
			"sub":   "user123",
			"email": "test@example.com",
			"exp":   time.Now().Add(time.Hour).Unix(),
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+health.RouteRelativePath, nil)
		req.Header.Set("x-amzn-oidc-data", token)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
