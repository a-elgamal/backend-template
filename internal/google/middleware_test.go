package google

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	testDB "alielgamal.com/myservice/internal/db/test"
	"alielgamal.com/myservice/internal/health"
	"alielgamal.com/myservice/internal/response"
)

func TestAuthMiddleware(t *testing.T) {
	zl := zerolog.New(zerolog.NewConsoleWriter())
	logger := zerologr.New(&zl)

	t.Run("Returns 401 with header not present", func(t *testing.T) {
		r := gin.Default()
		r.Use(authMiddleware(logger, 1234, "", 456))
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

	t.Run("Returns 401 with invalid header", func(t *testing.T) {
		r := gin.Default()
		r.Use(authMiddleware(logger, 1234, "", 456))
		health.SetupRoutes(r, testDB.NewDB(t), 0)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+health.RouteRelativePath, nil)
		req.Header.Set("x-goog-iap-jwt-assertion", "invalid")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		result := response.ErrorResponse{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, result.Err.Code)
	})
}
