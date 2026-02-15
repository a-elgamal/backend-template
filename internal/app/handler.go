package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"alielgamal.com/myservice/internal"
	"alielgamal.com/myservice/internal/db"
	"alielgamal.com/myservice/internal/response"
	"alielgamal.com/myservice/internal/stored"
)

// RouteRelativePath the relative path that handlers will be registered under
const RouteRelativePath = "apps"

const idParamName = "id"

// SetupRoutes adds app routes handling
func SetupRoutes(routes gin.IRoutes, logger logr.Logger, db db.DB) {
	setupRoutes(routes, logger, stored.NewStore[App](db, appTableName))
}

func setupRoutes(routes gin.IRoutes, logger logr.Logger, db stored.Store[App]) {
	h := handler{
		logger: logger.WithName("app.handler"),
		db:     db,
	}

	routes.POST(RouteRelativePath, h.addApp)
	routes.GET(RouteRelativePath+"/:"+idParamName, h.getApp)
	routes.PATCH(RouteRelativePath+"/:"+idParamName, h.patchApp)
	routes.GET(RouteRelativePath, h.listApps)
	routes.POST(RouteRelativePath+"/:"+idParamName+"/api-key", h.resetAPIKey)
}

type handler struct {
	logger logr.Logger
	db     stored.Store[App]
}

func (h *handler) addApp(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler.addApp")
	defer span.End()

	p := stored.Stored[App]{}
	if err := c.BindJSON(&p); err != nil {
		h.logger.Error(err, "failed to read request JSON")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			}})
		return
	}

	if p.ID == "" {
		p.ID = uuid.NewString()
	}

	p.Content.APIKey = uuid.NewString()
	p.Content.Disabled = false

	result, err := h.db.Add(ctx, internal.UserFromGinContext(c), p.ID, p.Content)
	if err != nil {
		h.logger.Error(err, "failed to add app to store", "id", p.ID)
		code := http.StatusInternalServerError
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code.Name() == "unique_violation" {
			h.logger.Error(pgErr, "PG Error while adding app", "code", pgErr.Code, "name", pgErr.Code.Name(), "column", pgErr.Column, "constraint", pgErr.Constraint, "table", pgErr.Table)
			code = http.StatusConflict
			err = fmt.Errorf("an app with the id '%v' already exists", p.ID)
		}
		c.JSON(code, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: code,
				Msg:  err.Error(),
			}})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *handler) getApp(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler.getApp")
	defer span.End()

	p := stored.Stored[App]{}
	if err := c.BindUri(&p); err != nil {
		h.logger.Error(err, "attempting to get app with invalid id", "id", p.ID)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			}})
		return
	}

	result, err := h.db.Get(ctx, p.ID)

	if err == sql.ErrNoRows {
		h.logger.Error(err, "cannot find app by id", "id", p.ID)
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusNotFound,
				Msg:  fmt.Sprintf("cannot find app with id: %v", p.ID),
			}})
		return
	}

	if err != nil {
		h.logger.Error(err, "failed to get app from store", "id", p.ID)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *handler) patchApp(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler.patchApp")
	defer span.End()

	p := stored.Stored[any]{}
	if err := c.BindUri(&p); err != nil {
		h.logger.Error(err, "attempting to patch app with invalid id", "id", p.ID)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			}})
		return
	}

	delta := map[string]any{}
	if err := c.BindJSON(&delta); err != nil {
		h.logger.Error(err, "unable to parse patch content", "id", p.ID)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			}})
		return
	}

	_, err := h.db.Patch(ctx, internal.UserFromGinContext(c), p.ID, delta)
	if err == sql.ErrNoRows {
		h.logger.Error(err, "attempt to patch a non-existing app", "id", p.ID)
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusNotFound,
				Msg:  err.Error(),
			}})
		return
	} else if err != nil {
		h.logger.Error(err, "failed to patch app in db", "id", p.ID)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *handler) listApps(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler.listApps")
	defer span.End()

	listCondition := []stored.Condition{}
	if v, hasDisabled := c.GetQuery(disabledJSONKey); hasDisabled {
		disabled, err := strconv.ParseBool(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{
				Err: response.ErrorDetail{
					Code: http.StatusBadRequest,
					Msg:  fmt.Sprintf("invalid value for disabled: %v", v),
				}})
			return
		}
		listCondition = append(listCondition, stored.Condition{Attribute: disabledJSONKey, Op: stored.EqualOperator, Value: disabled})
	}

	result, err := h.db.List(ctx, listCondition...)
	if err != nil {
		h.logger.Error(err, "failed to list apps from store")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *handler) resetAPIKey(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler.resetAPIKey")
	defer span.End()

	id := c.Param(idParamName)

	newKey := uuid.NewString()
	_, err := h.db.Patch(ctx, internal.UserFromGinContext(c), id, map[string]any{"apiKey": newKey})
	if err == sql.ErrNoRows {
		h.logger.Error(err, "attempt to reset API key for a non-existing app", "id", id)
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusNotFound,
				Msg:  fmt.Sprintf("cannot find app with id: %v", id),
			}})
		return
	} else if err != nil {
		h.logger.Error(err, "failed to reset API key", "id", id)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}})
		return
	}

	c.String(http.StatusOK, newKey)
}
