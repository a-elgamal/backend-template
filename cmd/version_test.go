package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"alielgamal.com/myservice/internal"
	"alielgamal.com/myservice/internal/config"
	testDB "alielgamal.com/myservice/internal/db/test"
)

func TestVersionCmd(t *testing.T) {
	t.Run("Prints the correct version information", func(t *testing.T) {
		appConfig, _ := config.ReadConfig()
		db, tearDown, err := testDB.SetupTestDB(t.Name(), 0, appConfig)
		require.NoError(t, err)
		defer tearDown()

		versionCmd := versionCmd(db)

		internal.Version = "testVersion"
		var buffer bytes.Buffer
		versionCmd.SetOutput(&buffer)
		err = versionCmd.ExecuteContext(context.Background())
		assert.NoError(t, err)
		output := buffer.String()
		assert.Contains(t, output, internal.Version)
	})
}
