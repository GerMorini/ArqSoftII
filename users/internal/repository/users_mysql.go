package repository

import (
	"fmt"
	"time"
	"users/internal/config"
	"users/internal/dao"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UsersRepository interface {
	Create(user dao.User) (dao.User, error)
	GetUserByID(id int) (dao.User, error)
	GetUserByUsername(username string) (dao.User, error)
	GetUserByEmail(email string) (dao.User, error)
}

type MySQLUsersRepository struct {
	db *gorm.DB
}

func NewMySQLUsersRepository(cfg config.MySQLConfig) *MySQLUsersRepository {
	var conn *gorm.DB

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DB_USER, cfg.DB_PASS, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_SCHEMA)
	log.Info("Conectando a la base de datos con dsn: ", dsn)

	// reintentamos conectarnos a la BDD varias veces
	for i := range 10 {
		time.Sleep(3 * time.Second)
		log.Debugf("Intentando conectar (%d/%d)\n", i+1, 10)

		var err error
		conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})

		if err != nil {
			log.Errorf("Error al conectar a la base de datos: %v\n", err)
			log.Error("No se pudo establecer conexion a la BDD")
			continue
		}

		break
	}

	log.Info("conexion a base de datos establecida")

	conn.AutoMigrate(&dao.User{})

	return &MySQLUsersRepository{db: conn}
}

func (r *MySQLUsersRepository) Create(user dao.User) (dao.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return dao.User{}, err
	}

	return user, nil
}

func (r *MySQLUsersRepository) GetUserByID(id int) (dao.User, error) {
	var userData dao.User

	err := r.db.Where("id_usuario = ?", id).First(&userData).Error
	if err != nil {
		return dao.User{}, err
	}

	return userData, nil
}

func (r *MySQLUsersRepository) GetUserByUsername(username string) (dao.User, error) {
	var usuario dao.User

	err := r.db.Where("username = ?", username).First(&usuario).Error
	if err != nil {
		log.Errorf("error al buscar un usuario por su nombre\nusername: %s\nerror: %v\n", username, err)
		return dao.User{}, err
	}

	return usuario, nil
}

func (r *MySQLUsersRepository) GetUserByEmail(email string) (dao.User, error) {
	var usuario dao.User

	err := r.db.Where("email = ?", email).First(&usuario).Error
	if err != nil {
		log.Errorf("error al buscar un usuario por su email\nemail: %s\nerror: %v\n", email, err)
		return dao.User{}, err
	}

	return usuario, nil
}
