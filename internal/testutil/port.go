// Package testutil contains utilities to support testing activities
package testutil

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"alielgamal.com/myservice/cmd"
	"alielgamal.com/myservice/internal/config"
	"alielgamal.com/myservice/internal/db"
	testDB "alielgamal.com/myservice/internal/db/test"
)

var httpPort atomic.Uint32

func init() {
	httpPort.Store(10000)
}

// SetUpIntegartionTest Starts a server and waits for it to be ready to accept connections. Returns: (BaseURL, TearDown, ApplicationConfig, DB)
func SetUpIntegartionTest(t *testing.T) (string, func(), config.Config, *db.SQLDB) {
	appConfig, v := config.ReadConfig()
	v.Set("SERVER.HTTP_ADDRESS", fmt.Sprintf(":%v", httpPort.Add(1)))
	db, dbTearDown, err := testDB.SetupTestDB(t.Name(), 0, appConfig)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	go cmd.ExecuteCommand(ctx, appConfig, db, "start")
	// REquire that server has started
	require.Eventually(t, func() bool {
		conn, err := net.Dial("tcp", appConfig.ServerConfig.GetHTTPAddress())
		if err == nil {
			defer conn.Close()
		}
		return err == nil
	}, 10*time.Second, 100*time.Millisecond, "Server didn't start")

	tearDown := func() {
		cancel()
		dbTearDown()
	}

	return fmt.Sprintf("http://%v", appConfig.ServerConfig.GetHTTPAddress()), tearDown, appConfig, db
}
