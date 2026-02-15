package item

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

func TestAddItem(t *testing.T) {
	t.Run("Successfully adds an item", func(t *testing.T) {
		mockStore := &storedTest.Store[Item]{}
		result := &stored.Stored[Item]{ID: "test-id", Content: Item{Name: "Test", Description: "A test item"}}
		mockStore.On("Add", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(result, nil)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		body, _ := json.Marshal(stored.Stored[Item]{ID: "test-id", Content: Item{Name: "Test", Description: "A test item"}})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/"+RouteRelativePath, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Returns error when store fails", func(t *testing.T) {
		mockStore := &storedTest.Store[Item]{}
		mockStore.On("Add", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return((*stored.Stored[Item])(nil), errors.New("db error"))

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		body, _ := json.Marshal(stored.Stored[Item]{ID: "test-id", Content: Item{Name: "Test"}})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/"+RouteRelativePath, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetItem(t *testing.T) {
	t.Run("Successfully gets an item", func(t *testing.T) {
		mockStore := &storedTest.Store[Item]{}
		result := &stored.Stored[Item]{ID: "test-id", Content: Item{Name: "Test"}}
		mockStore.On("Get", mock.Anything, "test-id").Return(result, nil)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+RouteRelativePath+"/test-id", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Returns 404 when item not found", func(t *testing.T) {
		mockStore := &storedTest.Store[Item]{}
		mockStore.On("Get", mock.Anything, "missing").Return((*stored.Stored[Item])(nil), sql.ErrNoRows)

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

func TestListItems(t *testing.T) {
	t.Run("Successfully lists items", func(t *testing.T) {
		mockStore := &storedTest.Store[Item]{}
		items := []stored.Stored[Item]{
			{ID: "1", Content: Item{Name: "Item 1"}},
			{ID: "2", Content: Item{Name: "Item 2"}},
		}
		mockStore.On("List", mock.Anything).Return(items, nil)

		r := gin.Default()
		setupRoutes(r, newLogger(), mockStore)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+RouteRelativePath, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
