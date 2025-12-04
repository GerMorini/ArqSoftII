package controllers

import (
	"activities/internal/dto"
	"activities/internal/repository"
	"activities/internal/services"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

type ActivitiesService interface {
	List(ctx context.Context) ([]dto.Activity, error)
	GetMany(ctx context.Context, ids []string) ([]dto.Activity, error)
	Create(ctx context.Context, actividad dto.ActivityAdministration) (dto.ActivityAdministration, error)
	GetByID(ctx context.Context, id string) (dto.ActivityAdministration, error)
	Update(ctx context.Context, id string, actividad dto.ActivityAdministration) (dto.ActivityAdministration, error)
	Delete(ctx context.Context, id string) error
	Inscribir(ctx context.Context, id string, userID string) (string, error)
	Desinscribir(ctx context.Context, id string, userID string) (string, error)
	GetInscripcionesByUserID(ctx context.Context, userID string) ([]string, error)
	GetActivitiesByUserID(ctx context.Context, userID string) (dto.Activities, error)
	GetStatistics(ctx context.Context) (dto.ActivityStatistics, error)
}

type ActivitiesController struct {
	service ActivitiesService
}

func NewActivitiesController(s ActivitiesService) *ActivitiesController {
	return &ActivitiesController{service: s}
}

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

func splitAndTrim(s string, sep string) []string {
	var result []string
	parts := strings.Split(s, sep)

	for _, item := range parts {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// Authentication is handled by middleware; handlers can read claims from context if needed

// GetActivities maneja GET /activities
func (c *ActivitiesController) GetActivities(ctx *gin.Context) {
	activities, err := c.service.List(ctx.Request.Context())
	if err != nil {
		log.Errorf("error al obtener todas las actividades: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities", "details": err.Error()})
		return
	}

	log.Infof("actividades obtenidas exitosamente")
	ctx.JSON(http.StatusOK, gin.H{"activities": activities, "count": len(activities)})
}

// GetManyActivities maneja GET /activities/many?ids=id1,id2,id3
func (c *ActivitiesController) GetManyActivities(ctx *gin.Context) {
	idsParam := ctx.Query("ids")
	if idsParam == "" {
		log.Warnf("peticion sin parametro ids")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ids query parameter is required"})
		return
	}

	// Parse comma-separated IDs
	ids := []string{}
	for _, id := range splitAndTrim(idsParam, ",") {
		if id != "" {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		log.Warnf("peticion con parametro ids vacio")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "at least one id is required"})
		return
	}

	activities, err := c.service.GetMany(ctx.Request.Context(), ids)
	if err != nil {
		log.Errorf("error al obtener actividades: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities", "details": err.Error()})
		return
	}

	log.Infof("actividades obtenidas exitosamente, solicitadas: %d, encontradas: %d", len(ids), len(activities))
	ctx.JSON(http.StatusOK, gin.H{"activities": activities, "count": len(activities)})
}

// CreateActivity maneja POST /activities
func (c *ActivitiesController) CreateActivity(ctx *gin.Context) {
	var newAct dto.ActivityAdministration
	if err := ctx.ShouldBindJSON(&newAct); err != nil {
		log.Warnf("error al parsear body JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	//admin only
	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		log.Warnf("token sin claims")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}
	if !isAdminFromClaims(claims) {
		log.Warnf("operacion sin privilegios para el usuario: %s@%s", claims["username"], ctx.RemoteIP())
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only admin users can create activities"})
		return
	}
	created, err := c.service.Create(ctx.Request.Context(), newAct)
	if err != nil {
		if errors.Is(err, services.ErrValidation) {
			log.Warnf("error de validación al crear actividad: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
			return
		}
		log.Errorf("fallo al crear actividad: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity", "details": err.Error()})
		return
	}

	log.Infof("actividad creada exitosamente por usuario: %s", claims["username"])
	ctx.JSON(http.StatusCreated, gin.H{"activity": created})
}

// GetActivityByID maneja GET /activities/:id
func (c *ActivitiesController) GetActivityByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		log.Warnf("peticion sin ID de actividad")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		log.Warnf("token sin claims")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}

	actAdmin, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrActivityNotFound) {
			log.Warnf("actividad no encontrada: %s", id)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		log.Errorf("error al obtener actividad %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity", "details": err.Error()})
		return
	}

	// If requester is admin, return the administration DTO (with users_inscritos, fecha)
	if isAdminFromClaims(claims) {
		log.Infof("actividad %s (admin view) obtenida exitosamente por usuario: %s", id, claims["username"])
		ctx.JSON(http.StatusOK, gin.H{"activity": actAdmin})
		return
	}

	// Non-admin users: return the public DTO (sin datos sensibles como usuarios inscritos)
	public := dto.Activity{
		ID:                 actAdmin.ID,
		Nombre:             actAdmin.Nombre,
		Descripcion:        actAdmin.Descripcion,
		Profesor:           actAdmin.Profesor,
		DiaSemana:          actAdmin.DiaSemana,
		HoraInicio:         actAdmin.HoraInicio,
		HoraFin:            actAdmin.HoraFin,
		CapacidadMax:       actAdmin.CapacidadMax,
		LugaresDisponibles: actAdmin.LugaresDisponibles,
		FotoUrl:            actAdmin.FotoUrl,
	}

	log.Infof("actividad %s (public view) obtenida exitosamente por usuario: %s", id, claims["username"])
	ctx.JSON(http.StatusOK, gin.H{"activity": public})
}

// Inscribir maneja POST /activities/:id/inscribir
func (c *ActivitiesController) Inscribir(ctx *gin.Context) {
	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		log.Warnf("token sin claims")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}

	uid, ok := getUserIDFromClaims(claims)
	if !ok {
		log.Warnf("id de usuario invalido en claims del token")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id in token claims"})
		return
	}

	if isAdminFromClaims(claims) {
		log.Warnf("intento de inscripcion por usuario admin: %s", claims["username"])
		ctx.JSON(http.StatusForbidden, gin.H{"error": "admin users cannot inscribe"})
		return
	}

	activityID := ctx.Param("id")
	if activityID == "" {
		log.Warnf("peticion de inscripcion sin ID de actividad")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "activity id required"})
		return
	}

	_, err := c.service.Inscribir(ctx.Request.Context(), activityID, uid)
	if err != nil {
		if errors.Is(err, repository.ErrActivityNotFound) {
			log.Warnf("actividad no encontrada para inscribir: %s", activityID)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}

		if errors.Is(err, repository.ErrActivityFull) {
			log.Warnf("actividad llena: %s", activityID)
			ctx.JSON(http.StatusConflict, gin.H{"error": "Activity is full"})
			return
		}

		if errors.Is(err, repository.ErrUserAlreadyInscribed) {
			log.Warnf("usuario %s ya inscrito en actividad %s", uid, activityID)
			ctx.JSON(http.StatusConflict, gin.H{"error": "User already inscribed in this activity"})
			return
		}

		log.Errorf("fallo al inscribir usuario %s en actividad %s: %v", uid, activityID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to inscribe", "details": err.Error()})
		return
	}

	log.Infof("usuario %s inscrito exitosamente en actividad %s", uid, activityID)
	ctx.JSON(http.StatusOK, gin.H{"status": "inscribed", "activity_id": activityID, "user_id": uid})
}

// Desinscribir maneja POST /activities/:id/desinscribir
func (c *ActivitiesController) Desinscribir(ctx *gin.Context) {
	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		log.Warnf("token sin claims")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}

	uid, ok := getUserIDFromClaims(claims)
	if !ok {
		log.Warnf("id de usuario invalido en claims del token")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id in token claims"})
		return
	}

	if isAdminFromClaims(claims) {
		log.Warnf("intento de desinscripcion por usuario admin: %s", claims["username"])
		ctx.JSON(http.StatusForbidden, gin.H{"error": "admin users cannot desinscribe"})
		return
	}

	activityID := ctx.Param("id")
	if activityID == "" {
		log.Warnf("peticion de desinscripcion sin ID de actividad")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "activity id required"})
		return
	}

	_, err := c.service.Desinscribir(ctx.Request.Context(), activityID, uid)
	if err != nil {
		if errors.Is(err, repository.ErrActivityNotFound) {
			log.Warnf("actividad no encontrada para desinscribir: %s", activityID)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}

		if errors.Is(err, repository.ErrUserNotInscribed) {
			log.Warnf("usuario %s no inscrito en actividad %s", uid, activityID)
			ctx.JSON(http.StatusConflict, gin.H{"error": "User not inscribed in this activity"})
			return
		}

		log.Errorf("fallo al desinscribir usuario %s de actividad %s: %v", uid, activityID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to desinscribe", "details": err.Error()})
		return
	}

	log.Infof("usuario %s desinscrito exitosamente de actividad %s", uid, activityID)
	ctx.JSON(http.StatusOK, gin.H{"status": "unsubscribed", "activity_id": activityID, "user_id": uid})
}

// UpdateActivity maneja PUT /activities/:id
func (c *ActivitiesController) UpdateActivity(ctx *gin.Context) {
	var toUpdate dto.ActivityAdministration
	if err := ctx.ShouldBindJSON(&toUpdate); err != nil {
		log.Warnf("error al parsear body JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		log.Warnf("peticion de actualizacion sin ID de actividad")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		log.Warnf("token sin claims")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}
	if !isAdminFromClaims(claims) {
		log.Warnf("operacion sin privilegios para el usuario: %s@%s", claims["username"], ctx.RemoteIP())
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only admin users can update activities"})
		return
	}

	updated, err := c.service.Update(ctx.Request.Context(), id, toUpdate)
	if err != nil {
		if errors.Is(err, repository.ErrActivityNotFound) {
			log.Warnf("actividad no encontrada para actualizar: %s", id)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		if errors.Is(err, services.ErrValidation) {
			log.Warnf("error de validación al actualizar actividad %s: %v", id, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
			return
		}
		log.Errorf("error al actualizar actividad %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity", "details": err.Error()})
		return
	}

	log.Infof("actividad %s actualizada exitosamente por usuario: %s", id, claims["username"])
	ctx.JSON(http.StatusOK, gin.H{"activity": updated})
}

// DeleteActivity maneja DELETE /activities/:id
func (c *ActivitiesController) DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		log.Warnf("peticion de eliminacion sin ID de actividad")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		log.Warnf("token sin claims")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}
	if !isAdminFromClaims(claims) {
		log.Warnf("operacion sin privilegios para el usuario: %s@%s", claims["username"], ctx.RemoteIP())
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only admin users can delete activities"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrActivityNotFound) {
			log.Warnf("actividad no encontrada para eliminar: %s", id)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		log.Errorf("error al eliminar actividad %s: %v", id, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity", "details": err.Error()})
		return
	}

	log.Infof("actividad %s eliminada exitosamente por usuario: %s", id, claims["username"])
	ctx.JSON(http.StatusNoContent, nil)
}

// GetInscripcionesByUserID maneja GET /inscriptions/:userId
func (c *ActivitiesController) GetInscripcionesByUserID(ctx *gin.Context) {
	userID := ctx.Param("userId")
	if userID == "" {
		log.Warnf("peticion de inscripciones sin userId")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userId parameter is required"})
		return
	}

	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		log.Warnf("token sin claims")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}

	requesterID, ok := getUserIDFromClaims(claims)
	if !ok {
		log.Warnf("id de usuario invalido en claims del token")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id in token claims"})
		return
	}

	// Only allow users to fetch their own inscripciones unless admin
	if requesterID != userID && !isAdminFromClaims(claims) {
		log.Warnf("usuario %s intento acceder a inscripciones de usuario %s sin permisos", requesterID, userID)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "cannot access other user's inscripciones"})
		return
	}

	inscripciones, err := c.service.GetInscripcionesByUserID(ctx.Request.Context(), userID)
	if err != nil {
		log.Errorf("error al obtener inscripciones para usuario %s: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch inscripciones", "details": err.Error()})
		return
	}

	log.Infof("inscripciones obtenidas exitosamente para usuario %s", userID)
	ctx.JSON(http.StatusOK, gin.H{"user_id": userID, "inscripciones": inscripciones, "count": len(inscripciones)})
}

// GetInscribedActivities maneja GET /inscriptions/data/:userId
func (c *ActivitiesController) GetInscribedActivities(ctx *gin.Context) {
	userID := ctx.Param("userId")
	if userID == "" {
		log.Warnf("peticion de actividades inscritas sin userId")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userId parameter is required"})
		return
	}

	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		log.Warnf("token sin claims")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token claims"})
		return
	}

	requesterID, ok := getUserIDFromClaims(claims)
	if !ok {
		log.Warnf("id de usuario invalido en claims del token")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id in token claims"})
		return
	}

	// Only allow users to fetch their own activities unless admin
	if requesterID != userID && !isAdminFromClaims(claims) {
		log.Warnf("usuario %s intento acceder a actividades inscritas de usuario %s sin permisos", requesterID, userID)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "cannot access other user's inscribed activities"})
		return
	}

	activities, err := c.service.GetActivitiesByUserID(ctx.Request.Context(), userID)
	if err != nil {
		log.Errorf("error al obtener inscripciones para usuario %s: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch inscripciones", "details": err.Error()})
		return
	}

	log.Infof("actividades inscritas obtenidas exitosamente para usuario %s: %d actividades", userID, len(activities))
	ctx.JSON(http.StatusOK, gin.H{"activities": activities, "count": len(activities)})
}

// GetStatistics obtiene estadísticas de actividades (solo admin)
func (c *ActivitiesController) GetStatistics(ctx *gin.Context) {
	claims, ok := getClaimsFromContext(ctx)
	if !ok {
		log.Warn("no se pudieron obtener claims del contexto")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if !isAdminFromClaims(claims) {
		log.Warn("usuario no es admin")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin access required"})
		return
	}

	stats, err := c.service.GetStatistics(ctx.Request.Context())
	if err != nil {
		log.WithError(err).Error("error al obtener estadísticas")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch statistics", "details": err.Error()})
		return
	}

	log.Info("estadísticas obtenidas exitosamente")
	ctx.JSON(http.StatusOK, stats)
}
