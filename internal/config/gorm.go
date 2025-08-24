package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"shipping-gateway/internal/entity"
	"time"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	// Load database configuration from viper
	username := viper.GetString("db.user")
	password := viper.GetString("db.password")
	host := viper.GetString("db.host")
	port := viper.GetInt("db.port")
	database := viper.GetString("db.name")
	idleConnection := viper.GetInt("db.pool.idle")
	maxConnection := viper.GetInt("db.pool.max")
	maxLifeTimeConnection := viper.GetInt("db.pool.lifetime")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	logLevel := gormLogger.Info
	if viper.GetBool("db.debug") {
		logLevel = gormLogger.Info
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.New(&logrusWriter{Logger: log}, gormLogger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logLevel,
		}),
	})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))
	log.Info("Database connection established successfully")

	AutoMigrate(db,
		entity.Area{},
		entity.Courier{},
		entity.CourierService{},
		entity.ShipmentTrackingLog{},
	)
	return db
}

func AutoMigrate(db *gorm.DB, models ...interface{}) {
	if err := db.AutoMigrate(models...); err != nil {
		panic(fmt.Sprintf("failed to auto migrate: %v", err))
	}
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}
