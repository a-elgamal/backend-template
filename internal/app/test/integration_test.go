package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"alielgamal.com/myservice/internal/app"
	"alielgamal.com/myservice/internal/stored"
	"alielgamal.com/myservice/internal/testutil"
)

func TestIntegration(t *testing.T) {
	baseURL, tearDown, _, _ := testutil.SetUpIntegartionTest(t)
	defer tearDown()

	t.Run("Add and Get App", func(t *testing.T) {
		newApp := stored.Stored[app.App]{
			ID:      "test-app",
			Content: app.App{},
		}

		body, err := json.Marshal(newApp)
		require.NoError(t, err)

		resp, err := http.Post(baseURL+"/internal/apps", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = http.Get(baseURL + "/internal/apps/test-app")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result stored.Stored[app.App]
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "test-app", result.ID)
		assert.NotEmpty(t, result.Content.APIKey)
		assert.False(t, result.Content.Disabled)
	})

	t.Run("List Apps", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/internal/apps")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
