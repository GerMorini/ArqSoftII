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

// MongoActivitiesRepository implementa ActivitiesRepository usando DB
type MongoActivitiesRepository struct {
	col *mongo.Collection // Referencia a la colección "activities" en DB
}

// NewMongoActivitiesRepository crea una nueva instancia del repository
// Recibe una referencia a la base de datos DB
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
	// ⏰ Timeout para evitar que la operación se cuelgue
	// Esto es importante en producción para no bloquear indefinidamente
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 🔍 Find() sin filtros retorna todos los documentos de la colección
	// bson.M{} es un filtro vacío (equivale a {} en DB shell)
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx) // ⚠️ IMPORTANTE: Siempre cerrar el cursor para liberar recursos

	// 📦 Decodificar resultados en slice de DAO (modelo DB)
	// Usamos el modelo DAO porque maneja ObjectID y tags BSON
	var daoActivities []dao.ActivityDAO
	if err := cur.All(ctx, &daoActivities); err != nil {
		return nil, err
	}
	// Convertir de DAO a Domain
	dtoActivities := make([]dto.Activity, len(daoActivities))
	for i, daoAct := range daoActivities {
		dtoActivities[i] = daoAct.ToDomain()
	}

	return dtoActivities, nil
}

// Create inserta un nuevo activity en DB
func (r *MongoActivitiesRepository) Create(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
	activityDAO := dao.FromDomainDAO(activity) // Convertir a modelo DAO

	activityDAO.ID = primitive.NewObjectID()
	activityDAO.FechaCreacion = time.Now().UTC()

	// Insertar en DB
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
// Consigna 3: Update parcial + actualizar updatedAt
func (r *MongoActivitiesRepository) Update(ctx context.Context, id string, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
	// Validar que el ID es un ObjectID válido
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
	if len(set) == 0 {
		return dto.ActivityAdministration{}, errors.New("no fields to update")
	}
	set["fecha_creacion"] = time.Now().UTC()

	update := bson.M{"$set": set}

	// Ejecutar la actualización
	result, err := r.col.UpdateByID(ctx, objID, update)
	if err != nil {
		return dto.ActivityAdministration{}, err
	}
	if result.MatchedCount == 0 {
		return dto.ActivityAdministration{}, errors.New("activity not found")
	}

	// Retornar el activity actualizado
	return r.GetByID(ctx, id)
}

// Delete elimina un activity por ID
func (r *MongoActivitiesRepository) Delete(ctx context.Context, id string) error {
	// Validar que el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	// Ejecutar la eliminación
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
// Consigna 2: Validar que el ID sea un ObjectID válido
func (r *MongoActivitiesRepository) GetByID(ctx context.Context, id string) (dto.ActivityAdministration, error) {
	// Validar que el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return dto.ActivityAdministration{}, errors.New("invalid ID format")
	}

	// Buscar en DB
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
	// Validar que el ID es un ObjectID válido
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
	capMax := 0
	fmt.Sscanf(act.CapacidadMax, "%d", &capMax)
	if len(act.UsersInscribed) >= (capMax) {
		return "", errors.New("activity is full")
	}
	// check user not already inscribed
	for _, uid := range act.UsersInscribed {
		if uid == userID {
			return "", errors.New("user already inscribed")
		}
	}
	// Ejecutar la actualización
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
	// Validar que el ID es un ObjectID válido
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
	for _, uid := range act.UsersInscribed {
		if uid == userID {
			found = true
			break
		}
	}
	if !found {
		return "", errors.New("user not inscribed in activity")
	}
	// Ejecutar la actualización
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
	// Buscar en DB
	cur, err := r.col.Find(ctx, bson.M{"usuarios_inscritos": idint})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx) // ⚠️ IMPORTANTE: Siempre cerrar el cursor para liberar recursos

	// 📦 Decodificar resultados en slice de DAO (modelo DB)
	// Usamos el modelo DAO porque maneja ObjectID y tags BSON
	var daoActivities []dao.ActivityDAO
	if err := cur.All(ctx, &daoActivities); err != nil {
		return nil, err
	}
	// Convertir de DAO a Domain
	activityIDs := make([]string, len(daoActivities))
	for i, daoAct := range daoActivities {
		activityIDs[i] = daoAct.ID.Hex()
	}
	return activityIDs, nil
}
