package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"alielgamal.com/myservice/internal"
	"alielgamal.com/myservice/internal/db"
)

// RouteRelativePath the relative path that the route will be configured at
const RouteRelativePath = "health"
const pingTimeout = time.Second
const dbErrorStatusText = "Error connecting to DB"

// SetupRoutes adds health routes handling
func SetupRoutes(routes gin.IRoutes, db db.DB, dbVersion uint) {
	setupRoutes(routes, db, dbVersion, internal.Version, internal.GitTag, internal.GitCommit, internal.BuildDate)
}

func setupRoutes(routes gin.IRoutes, db db.DB, dbVersion uint, version string, gitTag string, gitCommit string, buildDate string) {
	routes.GET(RouteRelativePath, (&handler{db, dbVersion, version, gitTag, gitCommit, buildDate}).healthHandler)
}

type handler struct {
	db             db.DB
	DBVersion      uint
	ServiceVersion string
	gitTag         string
	gitCommit      string
	buildDate      string
}

func (h *handler) healthHandler(ctx *gin.Context) {
	limitedCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	health := Health{
		DBVersion:      h.DBVersion,
		ServiceVersion: h.ServiceVersion,
		GitTag:         h.gitTag,
		GitCommit:      h.gitCommit,
		BuildDate:      h.buildDate,
		Status:         StatusOk,
	}

	err := h.db.PingContext(limitedCtx)
	if err != nil {
		health.Status = StatusError
		health.StatusText = dbErrorStatusText
		ctx.JSON(http.StatusInternalServerError, health)
	} else {
		ctx.JSON(http.StatusOK, health)
	}
}
