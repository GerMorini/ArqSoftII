package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"users/internal/config"
	"users/internal/dao"
	"users/internal/domain"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MySQLUsersRepository struct {
	db *gorm.DB
}

func NewMySQLUsersRepository(ctx context.Context, cfg config.MySQLConfig) *MySQLUsersRepository {
	var conn *gorm.DB

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DB_USER, cfg.DB_PASS, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_SCHEMA)
	log.Info("Conectando a la base de datos con dsn: ", dsn)

	// reintentamos conectarnos a la BDD varias veces
	for i := range 10 {
		time.Sleep(3 * time.Second)
		log.Debugf("Intentando conectar (%d/%d)\n", i+1, 10)

		var err error
		conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err != nil {
			log.Errorf("Error al conectar a la base de datos: %v\n", err)
			log.Error("No se pudo establecer conexion a la BDD")
			continue
		}

		break
	}

	log.Info("Conexion a base de datos establecida")

	return &MySQLUsersRepository{db: conn}
}

func (r *MySQLUsersRepository) List(ctx context.Context) ([]domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var daoUsers dao.Users
	if err := cur.All(ctx, &daoItems); err != nil {
		return nil, err
	}

	domainItems := make([]domain.User, len(daoItems))
	for i, daoItem := range daoItems {
		domainItems[i] = daoItem.ToDomain()
	}

	return domainItems, nil
}

func (r *MySQLUsersRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now

	daoItem := dao.FromDomain(item)

	res, err := r.col.InsertOne(ctx, daoItem)
	if err != nil {
		return domain.Item{}, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		daoItem.ID = oid
	}

	return daoItem.ToDomain(), nil
}

func (r *MySQLUsersRepository) GetByID(ctx context.Context, id string) (domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Item{}, errors.New("invalid id format")
	}

	var d dao.Item
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&d)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Item{}, errors.New("item not found")
		}
		return domain.Item{}, err
	}

	return d.ToDomain(), nil
}

func (r *MySQLUsersRepository) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Item{}, errors.New("invalid id format")
	}

	update := bson.M{
		"$set": bson.M{
			"name":       item.Name,
			"price":      item.Price,
			"updated_at": time.Now().UTC(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated dao.Item
	err = r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Item{}, errors.New("item not found")
		}
		return domain.Item{}, err
	}

	return updated.ToDomain(), nil
}

func (r *MySQLUsersRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("item not found")
	}
	return nil
}
