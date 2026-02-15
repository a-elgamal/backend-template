package cmd

import (
	"context"
	"net/http"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"alielgamal.com/myservice/internal/app"
	"alielgamal.com/myservice/internal/auth"
	"alielgamal.com/myservice/internal/aws"
	"alielgamal.com/myservice/internal/config"
	internalDB "alielgamal.com/myservice/internal/db"
	"alielgamal.com/myservice/internal/google"
	"alielgamal.com/myservice/internal/health"
	"alielgamal.com/myservice/internal/telemetry"
)

func startCmd(logger logr.Logger, db *internalDB.SQLDB, appConfig config.Config) *cobra.Command {

	return &cobra.Command{
		Use:     "start",
		Aliases: []string{"up"},
		Short:   "Start the service",
		Run: func(cmd *cobra.Command, _ []string) {
			logger.Info("Running any pending DB migrations...")
			dbVersion, err := internalDB.UpgradeDB(db.DB)
			if err != nil {
				logger.Error(err, "Error while upgrading DB")
				panic(err)
			}
			logger.Info("Migrations done", "DBVersion", dbVersion)

			logger.Info("Initializing service...")
			telemetryShutdownFunc, err := telemetry.SetupMonitoring(cmd.Context(), appConfig.TelemetryConfig)
			if err != nil {
				logger.Error(err, "failed to set up OpenTelemetry monitoring")
			}
			if _, err = otelsql.RegisterDBStatsMetrics(db.DB, otelsql.WithAttributes(
				semconv.DBSystemPostgreSQL,
			)); err != nil {
				logger.Error(err, "failed to register otelsql metrics")
			}

			router := gin.New()
			router.Use(telemetry.Middleware(logger), otelgin.Middleware("myservice"), gin.Recovery())
			if origins := appConfig.ServerConfig.CORSAllowedOrigins(); len(origins) > 0 {
				corsConfig := cors.DefaultConfig()
				corsConfig.AllowOrigins = origins
				corsConfig.AllowCredentials = true
				router.Use(cors.New(corsConfig))
			}
			health.SetupRoutes(router, db, dbVersion)

			internalRouter := router.Group("/internal")
			var authProvider auth.Provider
			if appConfig.AWSConfig.ALBAuthEnabled() {
				authProvider = aws.NewAuthProvider(appConfig.AWSConfig)
			} else if appConfig.GCPConfig.IAPAuthEnabled() {
				authProvider = google.NewAuthProvider(appConfig.GCPConfig)
			}
			if authProvider != nil {
				internalRouter.Use(authProvider.Middleware(logger))
			}
			health.SetupRoutes(internalRouter, db, dbVersion)
			app.SetupRoutes(internalRouter, logger, db)
			internalRouter.Static("portal", appConfig.ServerConfig.PortalPath())

			externalRouter := router.Group("/external")
			health.SetupRoutes(externalRouter, db, dbVersion)

			server := &http.Server{
				Addr:    appConfig.ServerConfig.GetHTTPAddress(),
				Handler: router,
			}

			logger.Info("Starting to serve connections...", "port", appConfig.ServerConfig.GetHTTPAddress())

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error(err, "failed to set up http server")
					panic(err)
				}
			}()

			<-cmd.Context().Done()
			logger.Info("Shutting down...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(appConfig.ServerConfig.ShutdownTimeoutSeconds())*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			telemetryShutdownFunc(ctx)

			logger.Info("Shutting down completed!")
		},
	}
}
