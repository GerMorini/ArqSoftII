package controllers

import (
	"activities/internal/domain"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ActivitiesService define la lógica de negocio para Activities
type ActivitiesService interface {
	List(ctx context.Context) ([]domain.Activity, error)
	Create(ctx context.Context, actividad domain.Activity) (domain.Activity, error)
	GetByID(ctx context.Context, id string) (domain.Activity, error)
	Update(ctx context.Context, id string, actividad domain.Activity) (domain.Activity, error)
	Delete(ctx context.Context, id string) error
}

// ActivitiesController maneja las peticiones HTTP para Activities
type ActivitiesController struct {
	service ActivitiesService
}

// NewActivitiesController crea una nueva instancia del controller
func NewActivitiesController(s ActivitiesService) *ActivitiesController {
	return &ActivitiesController{service: s}
}

// GetActivities maneja GET /activities
func (c *ActivitiesController) GetActivities(ctx *gin.Context) {
	activities, err := c.service.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"activities": activities, "count": len(activities)})
}

// CreateActivity maneja POST /activities
func (c *ActivitiesController) CreateActivity(ctx *gin.Context) {
	var newAct domain.Activity
	if err := ctx.ShouldBindJSON(&newAct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	created, err := c.service.Create(ctx.Request.Context(), newAct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"activity": created})
}

// GetActivityByID maneja GET /activities/:id
func (c *ActivitiesController) GetActivityByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	act, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "activity not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"activity": act})
}

// UpdateActivity maneja PUT /activities/:id
func (c *ActivitiesController) UpdateActivity(ctx *gin.Context) {
	var toUpdate domain.Activity
	if err := ctx.ShouldBindJSON(&toUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	updated, err := c.service.Update(ctx.Request.Context(), id, toUpdate)
	if err != nil {
		if err.Error() == "activity not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"activity": updated})
}

// DeleteActivity maneja DELETE /activities/:id
func (c *ActivitiesController) DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		if err.Error() == "activity not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
