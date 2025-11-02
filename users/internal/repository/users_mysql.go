package repository

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"sync"
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
	GetAll() ([]dao.User, error)
	Update(id int, user dao.User) (dao.User, error)
	Delete(id int) error
}

type MySQLUsersRepository struct {
	db             *gorm.DB
	cfg            config.MySQLConfig
	mu             sync.RWMutex // protege la variable isReconnecting (para cuando se debe escribir)
	isReconnecting bool
}

const MAX_RECONNECTION_TRIES = 5

func makeDSN(cfg config.MySQLConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DB_USER, cfg.DB_PASS, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_SCHEMA)
}

func NewMySQLUsersRepository(cfg config.MySQLConfig) *MySQLUsersRepository {
	var conn *gorm.DB

	dsn := makeDSN(cfg)
	log.Info("conectando a la base de datos con dsn: ", dsn)

	// reintentamos conectarnos a la BDD varias veces
	for tryN := range MAX_RECONNECTION_TRIES {
		time.Sleep(3 * time.Second)
		log.Warnf("Intentando conectar (%d/%d)\n", tryN+1, MAX_RECONNECTION_TRIES)

		var err error
		conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})

		if err != nil {
			log.Errorf("error al conectar a la base de datos: %v\n", err)
			log.Error("no se pudo establecer conexion a la BDD")
			continue
		}

		break
	}

	log.Info("conexion a base de datos establecida")
	conn.AutoMigrate(&dao.User{})

	repo := &MySQLUsersRepository{
		db:  conn,
		cfg: cfg,
	}

	return repo
}

func (r *MySQLUsersRepository) reconnect() {
	r.mu.Lock()
	if r.isReconnecting {
		r.mu.Unlock()
		return
	}

	r.isReconnecting = true
	r.mu.Unlock()

	defer func() {
		r.mu.Lock()
		r.isReconnecting = false
		r.mu.Unlock()
	}()

	dsn := makeDSN(r.cfg)
	log.Warn("Intentando reconectar a la base de datos...")
	for tryN := range MAX_RECONNECTION_TRIES {
		time.Sleep((1 << tryN) * time.Second)
		log.Warnf("Intento de reconexion (%d/%d)", tryN+1, MAX_RECONNECTION_TRIES)

		conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})

		if err != nil {
			log.Errorf("Error en intento de reconexion: %v", err)
			continue
		}

		r.mu.Lock()
		r.db = conn
		r.mu.Unlock()
		log.Info("Reconexion exitosa a la base de datos")
		return
	}

	log.Errorf("No se pudo reconectar a la base de datos despues de %d intentos", MAX_RECONNECTION_TRIES)
}

func (r *MySQLUsersRepository) isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, driver.ErrBadConn) {
		return true
	}

	errStr := strings.ToLower(err.Error())
	connectionErrors := []string{
		"connection refused",
		"broken pipe",
		"invalid connection",
		"connection reset by peer",
		"server has gone away",
		"error 2006",
		"eof",
		"no such host",
	}

	for _, connErr := range connectionErrors {
		if strings.Contains(errStr, connErr) {
			return true
		}
	}

	return false
}

func (r *MySQLUsersRepository) Create(user dao.User) (dao.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		if r.isConnectionError(err) {
			log.Errorf("error al conectar a la BDD: %s", err.Error())
			go r.reconnect()
		}
		return dao.User{}, err
	}

	return user, nil
}

func (r *MySQLUsersRepository) GetUserByID(id int) (dao.User, error) {
	var userData dao.User

	err := r.db.Where("id_usuario = ?", id).First(&userData).Error
	if err != nil {
		if r.isConnectionError(err) {
			log.Errorf("error al conectar a la BDD: %s", err.Error())
			go r.reconnect()
		}
		return dao.User{}, err
	}

	return userData, nil
}

func (r *MySQLUsersRepository) GetUserByUsername(username string) (dao.User, error) {
	var usuario dao.User

	err := r.db.Where("username = ?", username).First(&usuario).Error
	if err != nil {
		if r.isConnectionError(err) {
			log.Errorf("error al conectar a la BDD: %s", err.Error())
			go r.reconnect()
		}
		return dao.User{}, err
	}

	return usuario, nil
}

func (r *MySQLUsersRepository) GetUserByEmail(email string) (dao.User, error) {
	var usuario dao.User

	err := r.db.Where("email = ?", email).First(&usuario).Error
	if err != nil {
		if r.isConnectionError(err) {
			log.Errorf("error al conectar a la BDD: %s", err.Error())
			go r.reconnect()
		}
		return dao.User{}, err
	}

	return usuario, nil
}

func (r *MySQLUsersRepository) GetAll() ([]dao.User, error) {
	var usuarios []dao.User

	err := r.db.Find(&usuarios).Error
	if err != nil {
		if r.isConnectionError(err) {
			log.Errorf("error al conectar a la BDD: %s", err.Error())
			go r.reconnect()
		}
		return nil, err
	}

	return usuarios, nil
}

func (r *MySQLUsersRepository) Update(id int, user dao.User) (dao.User, error) {
	// Usar Save para asegurar que todos los campos se actualicen, incluyendo booleanos
	user.Id = id
	err := r.db.Save(&user).Error
	if err != nil {
		if r.isConnectionError(err) {
			log.Errorf("error al conectar a la BDD: %s", err.Error())
			go r.reconnect()
		}
		return dao.User{}, err
	}

	return user, nil
}

func (r *MySQLUsersRepository) Delete(id int) error {
	err := r.db.Where("id_usuario = ?", id).Delete(&dao.User{}).Error
	if r.isConnectionError(err) {
		log.Errorf("error al conectar a la BDD: %s", err.Error())
		go r.reconnect()
	}
	return err
}
