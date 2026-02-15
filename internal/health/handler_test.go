package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	testDB "alielgamal.com/myservice/internal/db/test"
)

func TestHandler(t *testing.T) {
	t.Run("When DB is accessible, returns ok status", func(t *testing.T) {
		h := Health{
			DBVersion:      1,
			ServiceVersion: "v1",
			GitTag:         "v1 Tag",
			GitCommit:      "abcew3r4",
			BuildDate:      "2024-04-16",
			Status:         StatusOk,
		}
		mockDB := testDB.NewDB(t)
		mockDB.On("PingContext", mock.Anything).Return(nil)
		r := gin.Default()
		setupRoutes(r, mockDB, h.DBVersion, h.ServiceVersion, h.GitTag, h.GitCommit, h.BuildDate)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+RouteRelativePath, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expectedJSON, err := json.Marshal(h)
		assert.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), w.Body.String())
	})

	t.Run("When DB is not accessible, returns error status", func(t *testing.T) {
		mockDB := testDB.NewDB(t)
		mockDB.On("PingContext", mock.Anything).Return(errors.New("mocked db error")).Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			_, deadlineSet := ctx.Deadline()
			assert.True(t, deadlineSet)
		})

		h := Health{
			Status:     StatusError,
			StatusText: dbErrorStatusText,
		}

		r := gin.Default()
		SetupRoutes(r, mockDB, 0)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+RouteRelativePath, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expectedJSON, err := json.Marshal(h)
		assert.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), w.Body.String())
	})
}
