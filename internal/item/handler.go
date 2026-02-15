package item

import (
	"database/sql"
	"fmt"
	"net/http"

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
const RouteRelativePath = "items"

const idParamName = "id"

// SetupRoutes adds item routes handling
func SetupRoutes(routes gin.IRoutes, logger logr.Logger, db db.DB) {
	setupRoutes(routes, logger, stored.NewStore[Item](db, itemTableName))
}

func setupRoutes(routes gin.IRoutes, logger logr.Logger, db stored.Store[Item]) {
	h := handler{
		logger: logger.WithName("item.handler"),
		db:     db,
	}

	routes.POST(RouteRelativePath, h.addItem)
	routes.GET(RouteRelativePath+"/:"+idParamName, h.getItem)
	routes.PATCH(RouteRelativePath+"/:"+idParamName, h.patchItem)
	routes.GET(RouteRelativePath, h.listItems)
}

type handler struct {
	logger logr.Logger
	db     stored.Store[Item]
}

func (h *handler) addItem(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler.addItem")
	defer span.End()

	p := stored.Stored[Item]{}
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

	result, err := h.db.Add(ctx, internal.UserFromGinContext(c), p.ID, p.Content)
	if err != nil {
		h.logger.Error(err, "failed to add item to store", "id", p.ID)
		code := http.StatusInternalServerError
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code.Name() == "unique_violation" {
			h.logger.Error(pgErr, "PG Error while adding item", "code", pgErr.Code, "name", pgErr.Code.Name(), "column", pgErr.Column, "constraint", pgErr.Constraint, "table", pgErr.Table)
			code = http.StatusConflict
			err = fmt.Errorf("an item with the id '%v' already exists", p.ID)
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

func (h *handler) getItem(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler.getItem")
	defer span.End()

	p := stored.Stored[Item]{}
	if err := c.BindUri(&p); err != nil {
		h.logger.Error(err, "attempting to get item with invalid id", "id", p.ID)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			}})
		return
	}

	result, err := h.db.Get(ctx, p.ID)

	if err == sql.ErrNoRows {
		h.logger.Error(err, "cannot find item by id", "id", p.ID)
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusNotFound,
				Msg:  fmt.Sprintf("cannot find item with id: %v", p.ID),
			}})
		return
	}

	if err != nil {
		h.logger.Error(err, "failed to get item from store", "id", p.ID)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *handler) patchItem(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler.patchItem")
	defer span.End()

	p := stored.Stored[any]{}
	if err := c.BindUri(&p); err != nil {
		h.logger.Error(err, "attempting to patch item with invalid id", "id", p.ID)
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
		h.logger.Error(err, "attempt to patch a non-existing item", "id", p.ID)
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusNotFound,
				Msg:  err.Error(),
			}})
		return
	} else if err != nil {
		h.logger.Error(err, "failed to patch item in db", "id", p.ID)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *handler) listItems(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler.listItems")
	defer span.End()

	listCondition := []stored.Condition{}
	if v, hasName := c.GetQuery(nameJSONKey); hasName {
		listCondition = append(listCondition, stored.Condition{Attribute: nameJSONKey, Op: stored.EqualOperator, Value: v})
	}

	result, err := h.db.List(ctx, listCondition...)
	if err != nil {
		h.logger.Error(err, "failed to list items from store")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Err: response.ErrorDetail{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}})
		return
	}

	c.JSON(http.StatusOK, result)
}
