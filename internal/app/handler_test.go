package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"alielgamal.com/myservice/internal/response"
	"alielgamal.com/myservice/internal/stored"
	storedTest "alielgamal.com/myservice/internal/stored/test"
)

func newLogger() zerologr.Logger {
	zl := zerolog.New(zerolog.NewConsoleWriter())
	return zerologr.New(&zl)
}

func TestAddApp(t *testing.T) {
	t.Run("Successfully adds an app", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		result := &stored.Stored[App]{ID: "test-id", Content: App{APIKey: "generated-key", Disabled: false}}
		mockStore.On("Add", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(result, nil)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		body, _ := json.Marshal(stored.Stored[App]{ID: "test-id", Content: App{}})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/"+RouteRelativePath, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Returns error when store fails", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		mockStore.On("Add", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return((*stored.Stored[App])(nil), errors.New("db error"))

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		body, _ := json.Marshal(stored.Stored[App]{ID: "test-id", Content: App{}})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/"+RouteRelativePath, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAddAppAutoGeneratesAPIKey(t *testing.T) {
	mockStore := &storedTest.Store[App]{}
	result := &stored.Stored[App]{ID: "test-id", Content: App{APIKey: "key", Disabled: false}}
	mockStore.On("Add", mock.Anything, mock.Anything, "test-id", mock.MatchedBy(func(a App) bool {
		return a.APIKey != "" && !a.Disabled
	})).Return(result, nil)

	r := gin.Default()
	setupRoutes(r, newLogger(), mockStore)

	body, _ := json.Marshal(stored.Stored[App]{ID: "test-id", Content: App{}})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/"+RouteRelativePath, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetApp(t *testing.T) {
	t.Run("Successfully gets an app", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		result := &stored.Stored[App]{ID: "test-id", Content: App{APIKey: "some-key"}}
		mockStore.On("Get", mock.Anything, "test-id").Return(result, nil)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+RouteRelativePath+"/test-id", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Returns 404 when app not found", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		mockStore.On("Get", mock.Anything, "missing").Return((*stored.Stored[App])(nil), sql.ErrNoRows)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+RouteRelativePath+"/missing", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		result := response.ErrorResponse{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, result.Err.Code)
	})
}

func TestGetAppInternalError(t *testing.T) {
	mockStore := &storedTest.Store[App]{}
	mockStore.On("Get", mock.Anything, "err-id").Return((*stored.Stored[App])(nil), errors.New("db error"))

	r := gin.Default()
	setupRoutes(r, newLogger(), mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/"+RouteRelativePath+"/err-id", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPatchApp(t *testing.T) {
	t.Run("Successfully patches an app", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		result := &stored.Stored[App]{ID: "test-id", Content: App{APIKey: "key", Disabled: true}}
		mockStore.On("Patch", mock.Anything, mock.Anything, "test-id", mock.Anything).Return(result, nil)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		body, _ := json.Marshal(map[string]any{"disabled": true})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, "/"+RouteRelativePath+"/test-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Returns 404 when app not found", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		mockStore.On("Patch", mock.Anything, mock.Anything, "missing", mock.Anything).Return((*stored.Stored[App])(nil), sql.ErrNoRows)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		body, _ := json.Marshal(map[string]any{"disabled": true})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, "/"+RouteRelativePath+"/missing", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Returns 500 on internal error", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		mockStore.On("Patch", mock.Anything, mock.Anything, "err-id", mock.Anything).Return((*stored.Stored[App])(nil), errors.New("db error"))

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		body, _ := json.Marshal(map[string]any{"disabled": true})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, "/"+RouteRelativePath+"/err-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestListApps(t *testing.T) {
	t.Run("Successfully lists apps", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		apps := []stored.Stored[App]{
			{ID: "1", Content: App{APIKey: "key-1"}},
			{ID: "2", Content: App{APIKey: "key-2"}},
		}
		mockStore.On("List", mock.Anything).Return(apps, nil)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+RouteRelativePath, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Filters by disabled query param", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		apps := []stored.Stored[App]{
			{ID: "1", Content: App{APIKey: "key-1", Disabled: true}},
		}
		mockStore.On("List", mock.Anything, stored.Condition{Attribute: disabledJSONKey, Op: stored.EqualOperator, Value: true}).Return(apps, nil)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+RouteRelativePath+"?disabled=true", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestListAppsError(t *testing.T) {
	mockStore := &storedTest.Store[App]{}
	mockStore.On("List", mock.Anything).Return([]stored.Stored[App](nil), errors.New("db error"))

	r := gin.Default()
	setupRoutes(r, newLogger(), mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/"+RouteRelativePath, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestListAppsInvalidDisabled(t *testing.T) {
	mockStore := &storedTest.Store[App]{}

	r := gin.Default()
	setupRoutes(r, newLogger(), mockStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/"+RouteRelativePath+"?disabled=notabool", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetAPIKey(t *testing.T) {
	t.Run("Successfully resets API key", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		result := &stored.Stored[App]{ID: "test-id", Content: App{APIKey: "new-key"}}
		mockStore.On("Patch", mock.Anything, mock.Anything, "test-id", mock.Anything).Return(result, nil)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/"+RouteRelativePath+"/test-id/api-key", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Body.String())
	})

	t.Run("Returns 404 when app not found", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		mockStore.On("Patch", mock.Anything, mock.Anything, "missing", mock.Anything).Return((*stored.Stored[App])(nil), sql.ErrNoRows)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/"+RouteRelativePath+"/missing/api-key", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Returns 500 on internal error", func(t *testing.T) {
		mockStore := &storedTest.Store[App]{}
		mockStore.On("Patch", mock.Anything, mock.Anything, "err-id", mock.Anything).Return((*stored.Stored[App])(nil), errors.New("db error"))

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/"+RouteRelativePath+"/err-id/api-key", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
