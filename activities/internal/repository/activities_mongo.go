package repository

import (
	"activities/internal/dao"
	"activities/internal/dto"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoActivitiesRepository struct {
	col *mongo.Collection
}

func NewMongoActivitiesRepository(ctx context.Context, uri, dbName, collectionName string) *MongoActivitiesRepository {
	opt := options.Client().ApplyURI(uri)
	opt.SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
		return nil
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		log.Fatalf("Error pinging DB: %v", err)
		return nil
	}

	return &MongoActivitiesRepository{
		col: client.Database(dbName).Collection(collectionName), // Conecta con la colección "activities"
	}
}

// List obtiene todos los activities de DB
func (r *MongoActivitiesRepository) List(ctx context.Context) ([]dto.Activity, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var daoActivities []dao.ActivityDAO
	if err := cur.All(ctx, &daoActivities); err != nil {
		return nil, err
	}

	dtoActivities := make([]dto.Activity, len(daoActivities))
	for i, daoAct := range daoActivities {
		dtoActivities[i] = daoAct.ToDomain()
	}

	return dtoActivities, nil
}

// GetMany obtiene multiples activities por IDs (ignora IDs no encontrados)
func (r *MongoActivitiesRepository) GetMany(ctx context.Context, ids []string) ([]dto.Activity, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Convert string IDs to ObjectIDs, skip invalid ones
	objectIDs := []primitive.ObjectID{}
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objectIDs = append(objectIDs, objID)
		}
	}

	if len(objectIDs) == 0 {
		return []dto.Activity{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var daoActivities []dao.ActivityDAO
	if err := cur.All(ctx, &daoActivities); err != nil {
		return nil, err
	}

	dtoActivities := make([]dto.Activity, len(daoActivities))
	for i, daoAct := range daoActivities {
		dtoActivities[i] = daoAct.ToDomain()
	}

	return dtoActivities, nil
}

// Create inserta un nuevo activity en DB
func (r *MongoActivitiesRepository) Create(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
	activityDAO := dao.FromDomainDAO(activity)

	activityDAO.ID = primitive.NewObjectID()
	activityDAO.FechaCreacion = time.Now().UTC()

	_, err := r.col.InsertOne(ctx, activityDAO)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return dto.ActivityAdministration{}, errors.New("activity with the same ID already exists")
		}
		return dto.ActivityAdministration{}, err
	}

	return dao.ToDomainAdministration(activityDAO), nil
}

// Update actualiza un activity existente
func (r *MongoActivitiesRepository) Update(ctx context.Context, id string, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return dto.ActivityAdministration{}, errors.New("invalid ID format")
	}

	// Preparar los campos a actualizar
	// Construir update dinámico a partir de campos no vacíos
	set := bson.M{}
	if activity.Nombre != "" {
		set["nombre"] = activity.Nombre
	}
	if activity.Descripcion != "" {
		set["descripcion"] = activity.Descripcion
	}
	if activity.Profesor != "" {
		set["profesor_id"] = activity.Profesor
	}
	if activity.DiaSemana != "" {
		set["dia_semana"] = activity.DiaSemana
	}
	if activity.HoraInicio != "" {
		set["hora_inicio"] = activity.HoraInicio
	}
	if activity.HoraFin != "" {
		set["hora_fin"] = activity.HoraFin
	}
	if activity.FotoUrl != "" {
		set["foto_url"] = activity.FotoUrl
	}
	if activity.CapacidadMax > 0 {
		set["capacidad_max"] = activity.CapacidadMax
	}
	// If admin provided explicit users list, convert to integer slice and set it
	if activity.UsersInscribed != nil {
		set["usuarios_inscritos"] = activity.UsersInscribed
	}
	if len(set) == 0 {
		return dto.ActivityAdministration{}, errors.New("no fields to update")
	}
	set["fecha_creacion"] = time.Now().UTC()

	update := bson.M{"$set": set}

	result, err := r.col.UpdateByID(ctx, objID, update)
	if err != nil {
		return dto.ActivityAdministration{}, err
	}
	if result.MatchedCount == 0 {
		return dto.ActivityAdministration{}, errors.New("activity not found")
	}

	return r.GetByID(ctx, id)
}

// Delete elimina un activity por ID
func (r *MongoActivitiesRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	result, err := r.col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("activity not found")
	}

	return nil
}

// GetByID busca un activity por su ID
func (r *MongoActivitiesRepository) GetByID(ctx context.Context, id string) (dto.ActivityAdministration, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return dto.ActivityAdministration{}, errors.New("invalid ID format")
	}

	var activityDAO dao.ActivityDAO
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&activityDAO)
	if err != nil {
		// Manejar caso de no encontrado
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dto.ActivityAdministration{}, errors.New("activity not found")
		}
		return dto.ActivityAdministration{}, err
	}

	return dao.ToDomainAdministration(activityDAO), nil
}

func (r *MongoActivitiesRepository) Inscribir(ctx context.Context, id string, userID string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", errors.New("invalid ID format")
	}

	idint, err := strconv.Atoi(userID)
	if err != nil {
		return "", fmt.Errorf("error converting userID to int: %w", err)
	}

	act, err := r.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("error getting activity from repository: %w", err)
	}

	if len(act.UsersInscribed) >= (act.CapacidadMax) {
		return "", errors.New("activity is full")
	}

	// check user not already inscribed
	var userID_int int
	fmt.Sscanf(userID, "%d", &userID_int)
	for _, uid := range act.UsersInscribed {
		if uid == userID_int {
			return "", errors.New("user already inscribed")
		}
	}

	update := bson.M{"$push": bson.M{"usuarios_inscritos": idint}}
	result, err := r.col.UpdateByID(ctx, objID, update)
	if err != nil {
		return "", err
	}
	if result.MatchedCount == 0 {
		return "", errors.New("activity not found")
	}
	return id, nil
}

func (r *MongoActivitiesRepository) Desinscribir(ctx context.Context, id string, userID string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", errors.New("invalid ID format")
	}

	idint, err := strconv.Atoi(userID)
	if err != nil {
		return "", fmt.Errorf("error converting userID to int: %w", err)
	}

	act, err := r.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("error getting activity from repository: %w", err)
	}

	found := false
	var userID_int int
	fmt.Sscanf(userID, "%d", &userID_int)
	for _, uid := range act.UsersInscribed {
		if uid == userID_int {
			found = true
			break
		}
	}
	if !found {
		return "", errors.New("user not inscribed in activity")
	}

	update := bson.M{"$pull": bson.M{"usuarios_inscritos": idint}}
	result, err := r.col.UpdateByID(ctx, objID, update)
	if err != nil {
		return "", err
	}
	if result.MatchedCount == 0 {
		return "", errors.New("activity not found")
	}
	return id, nil
}

func (r *MongoActivitiesRepository) GetInscripcionesByUserID(ctx context.Context, userID string) ([]string, error) {
	idint, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("error converting userID to int: %w", err)
	}

	cur, err := r.col.Find(ctx, bson.M{"usuarios_inscritos": idint})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var daoActivities []dao.ActivityDAO
	if err := cur.All(ctx, &daoActivities); err != nil {
		return nil, err
	}

	activityIDs := make([]string, len(daoActivities))
	for i, daoAct := range daoActivities {
		activityIDs[i] = daoAct.ID.Hex()
	}
	return activityIDs, nil
}

func (r *MongoActivitiesRepository) GetActivitiesByUserID(ctx context.Context, userID string) (dto.Activities, error) {
	idint, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("error converting userID to int: %w", err)
	}

	cur, err := r.col.Find(ctx, bson.M{"usuarios_inscritos": idint})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var daoActivities []dao.ActivityDAO
	if err := cur.All(ctx, &daoActivities); err != nil {
		return nil, err
	}

	var dtoActivities dto.Activities = make(dto.Activities, len(daoActivities))
	for i, dao := range daoActivities {
		dtoActivities[i] = dao.ToDomain()
	}

	return dtoActivities, nil
}

// ListAllForAdmin retrieves all activities with full administration data (including UsersInscribed)
func (r *MongoActivitiesRepository) ListAllForAdmin(ctx context.Context) ([]dto.ActivityAdministration, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var daoActivities []dao.ActivityDAO
	if err := cur.All(ctx, &daoActivities); err != nil {
		return nil, err
	}

	dtoActivities := make([]dto.ActivityAdministration, len(daoActivities))
	for i, daoAct := range daoActivities {
		dtoActivities[i] = dao.ToDomainAdministration(daoAct)
	}

	return dtoActivities, nil
}
