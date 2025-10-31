package controllers

import (
	"context"
	"net/http"
	"search/internal/dto"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ItemsService interface {
	List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error)
	Create(ctx context.Context, activity dto.Activity) (dto.Activity, error)
	GetByID(ctx context.Context, id string) (dto.Activity, error)
	Update(ctx context.Context, id string, activity dto.Activity) (dto.Activity, error)
	Delete(ctx context.Context, id string) error
}

type ItemsController struct {
	service ItemsService
}

const (
	listDefaultPage  = 1
	listDefaultCount = 10
)

func NewActivitiesController(activitysService ItemsService) *ItemsController {
	return &ItemsController{
		service: activitysService,
	}
}

func (c *ItemsController) List(ctx *gin.Context) {
	filters := dto.SearchFilters{}

	if titulo := ctx.Query("titulo"); titulo != "" {
		filters.Titulo = titulo
	}

	if descripcion := ctx.Query("descripcion"); descripcion != "" {
		filters.Descripcion = descripcion
	}

	if diaSemana := ctx.Query("diaSemana"); diaSemana != "" {
		filters.DiaSemana = diaSemana
	}

	filters.SortBy = ctx.DefaultQuery("sortBy", "fecha_creacion asc")

	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filters.Page = page
		}
	} else {
		filters.Page = listDefaultPage
	}

	if countStr := ctx.Query("count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil {
			filters.Count = count
		}
	} else {
		filters.Count = listDefaultCount
	}

	resp, err := c.service.List(ctx.Request.Context(), filters)
	if err != nil {
		log.Errorf("error al realizar busqueda: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch activitys",
			"details": err.Error(),
		})
		return
	}

	log.Infof("exito al realizar busqueda")
	ctx.JSON(http.StatusOK, resp)
}

func (c *ItemsController) CreateActivity(ctx *gin.Context) {
	var activity dto.Activity
	if err := ctx.ShouldBindJSON(&activity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	created, err := c.service.Create(ctx.Request.Context(), activity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create activity",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"activity": created,
	})
}

func (c *ItemsController) GetActivityByID(ctx *gin.Context) {
	id := ctx.Param("id")

	activity, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get activity", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"activity": activity})
}

func (c *ItemsController) UpdateActivity(ctx *gin.Context) {
	id := ctx.Param("id")

	var activity dto.Activity

	if err := ctx.BindJSON(&activity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "consuta con formato incorrecto"})
		return
	}

	updated, err := c.service.Update(ctx, id, activity)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

func (c *ItemsController) DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx, id); err != nil {
		if strings.Contains(err.Error(), "invalid") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete activity", "details": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
