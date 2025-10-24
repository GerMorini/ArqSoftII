package controllers

import (
	"activities/internal/dto"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// ActivitiesService define la lógica de negocio para Activities
type ActivitiesService interface {
	List(ctx context.Context) ([]dto.Activity, error)
	Create(ctx context.Context, actividad dto.ActivityAdministration) (dto.ActivityAdministration, error)
	GetByID(ctx context.Context, id string) (dto.ActivityAdministration, error)
	Update(ctx context.Context, id string, actividad dto.ActivityAdministration) (dto.ActivityAdministration, error)
	Delete(ctx context.Context, id string) error
	Inscribir(ctx context.Context, id string, userID string) (string, error)
	Desinscribir(ctx context.Context, id string, userID string) (string, error)
	GetInscripcionesByUserID(ctx context.Context, userID string) ([]string, error)
}

// ActivitiesController maneja las peticiones HTTP para Activities
type ActivitiesController struct {
	service ActivitiesService
}

// NewActivitiesController crea una nueva instancia del controller
func NewActivitiesController(s ActivitiesService) *ActivitiesController {
	return &ActivitiesController{service: s}
}

// Helpers to read claims from context
func getClaimsFromContext(c *gin.Context) (jwt.MapClaims, bool) {
	v, ok := c.Get("claims")
	if !ok {
		return nil, false
	}
	claims, ok := v.(jwt.MapClaims)
	return claims, ok
}

func getUserIDFromClaims(claims jwt.MapClaims) (string, bool) {
	if claims == nil {
		return "", false
	}
	// claims store id_usuario as numeric or string depending on origin; handle both
	if idv, ok := claims["id_usuario"]; ok {
		switch t := idv.(type) {
		case float64:
			return fmt.Sprintf("%d", int(t)), true
		case int:
			return fmt.Sprintf("%d", t), true
		case string:
			return t, true
		}
	}
	return "", false
}

func isAdminFromClaims(claims jwt.MapClaims) bool {
	if claims == nil {
		return false
	}
	if adm, ok := claims["is_admin"]; ok {
		switch v := adm.(type) {
		case bool:
			return v
		case float64:
			return v != 0
		case int:
			return v != 0
		case string:
			return v == "true" || v == "1"
		}
	}
	return false
}

// Authentication is handled by middleware; handlers can read claims from context if needed

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
	var newAct dto.ActivityAdministration
	if err := ctx.ShouldBindJSON(&newAct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	//admin only
	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}
	if !isAdminFromClaims(claims) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only admin users can create activities"})
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
	//admin only
	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}
	if !isAdminFromClaims(claims) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only admin users can view activity by ID"})
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

// Inscribir maneja POST /activities/:id/inscribir
func (c *ActivitiesController) Inscribir(ctx *gin.Context) {
	// auth middleware ensures claims exist
	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}

	uid, ok := getUserIDFromClaims(claims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id in token claims"})
		return
	}

	// Non-admin users can inscribirse; admins may also inscribirse but typically shouldn't
	// Here we allow non-admin users explicitly. If admin, block.
	if isAdminFromClaims(claims) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "admin users cannot inscribe"})
		return
	}

	activityID := ctx.Param("id")
	if activityID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "activity id required"})
		return
	}

	// call service Inscribir with activityID and userID (current repository uses only activity id; user id handling done inside repo)
	_, err := c.service.Inscribir(ctx.Request.Context(), activityID, uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to inscribe", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "inscribed", "activity_id": activityID, "user_id": uid})
}

// Desinscribir maneja POST /activities/:id/desinscribir
func (c *ActivitiesController) Desinscribir(ctx *gin.Context) {
	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}

	uid, ok := getUserIDFromClaims(claims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id in token claims"})
		return
	}

	if isAdminFromClaims(claims) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "admin users cannot desinscribe"})
		return
	}

	activityID := ctx.Param("id")
	if activityID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "activity id required"})
		return
	}

	_, err := c.service.Desinscribir(ctx.Request.Context(), activityID, uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to desinscribe", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "unsubscribed", "activity_id": activityID, "user_id": uid})
}

// UpdateActivity maneja PUT /activities/:id
func (c *ActivitiesController) UpdateActivity(ctx *gin.Context) {
	var toUpdate dto.ActivityAdministration
	if err := ctx.ShouldBindJSON(&toUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	//admin only
	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}
	if !isAdminFromClaims(claims) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only admin users can update activities"})
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

	//admin only
	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}
	if !isAdminFromClaims(claims) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only admin users can delete activities"})
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

// GetInscripcionesByUserID maneja GET /inscripciones/:userId
func (c *ActivitiesController) GetInscripcionesByUserID(ctx *gin.Context) {
	userID := ctx.Param("userId")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userId parameter is required"})
		return
	}
	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}
	requesterID, ok := getUserIDFromClaims(claims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id in token claims"})
		return
	}
	// Only allow users to fetch their own inscripciones unless admin
	if requesterID != userID && !isAdminFromClaims(claims) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "cannot access other user's inscripciones"})
		return
	}
	inscripciones, err := c.service.GetInscripcionesByUserID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch inscripciones", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user_id": userID, "inscripciones": inscripciones, "count": len(inscripciones)})
}
