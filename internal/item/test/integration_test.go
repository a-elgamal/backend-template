package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"alielgamal.com/myservice/internal/item"
	"alielgamal.com/myservice/internal/stored"
	"alielgamal.com/myservice/internal/testutil"
)

func TestIntegration(t *testing.T) {
	baseURL, tearDown, _, _ := testutil.SetUpIntegartionTest(t)
	defer tearDown()

	t.Run("Add and Get Item", func(t *testing.T) {
		newItem := stored.Stored[item.Item]{
			ID:      "test-item",
			Content: item.Item{Name: "Test Item", Description: "A test item"},
		}

		body, err := json.Marshal(newItem)
		require.NoError(t, err)

		resp, err := http.Post(baseURL+"/internal/items", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = http.Get(baseURL + "/internal/items/test-item")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result stored.Stored[item.Item]
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "test-item", result.ID)
		assert.Equal(t, "Test Item", result.Content.Name)
	})

	t.Run("List Items", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/internal/items")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
