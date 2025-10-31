package controllers

import (
	"context"
	"net/http"
	"search/internal/dto"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ItemsService interface {
	List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error)
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
