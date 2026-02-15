/*
Package cmd sets up the commands available for execution. The default command is 'start' or 'up' which starts the server
*/
package cmd

import (
	"context"
	"database/sql/driver"

	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/XSAM/otelsql"
	// Postgres driver import
	_ "github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"

	"alielgamal.com/myservice/internal/config"
	internalDB "alielgamal.com/myservice/internal/db"
	"alielgamal.com/myservice/internal/telemetry"
)

// ExecuteCommand executes the command line command. if db is nil, a new DB will be creatd based on the configuration in the passed AppConfig.
func ExecuteCommand(ctx context.Context, appConfig config.Config, db *internalDB.SQLDB, args ...string) error {
	var err error
	logger, err := telemetry.SetupLogger(appConfig.TelemetryConfig)

	if err != nil {
		logger.Error(err, "Error while setting logger")
	}

	if db == nil {
		// Get a database driver.Connector for a fixed configuration.
		logger.Info("Connecting to db...", "URL", appConfig.DBConfig.GetURL())
		var connector driver.Connector
		connector, err = pq.NewConnector(appConfig.DBConfig.GetURL())
		if err != nil {
			logger.Error(err, "Error while connecting to DB")
			return err
		}

		// Wrap the driver.Connector with ocsql.
		connector = ocsql.WrapConnector(connector, ocsql.WithAllTraceOptions())

		// Use the wrapped driver.Connector.
		db = &internalDB.SQLDB{
			DB: otelsql.OpenDB(connector, otelsql.WithAttributes(
				semconv.DBSystemPostgreSQL)),
		}
		defer db.DB.Close()
	}

	// RootCmd is the top-level command that all other commands are added to
	rootCmd := &cobra.Command{
		Use:   "myservice",
		Short: "My Service",
	}

	rootCmd.AddCommand(startCmd(logger, db, appConfig))
	rootCmd.AddCommand(migrateCmd(db))
	rootCmd.AddCommand(versionCmd(db))

	if len(args) > 0 {
		rootCmd.SetArgs(args)
	}
	return rootCmd.ExecuteContext(ctx)
}
