package starter

import (
	"fmt"
	"github.com/feature-vector/harbor/base/conf"
	"github.com/feature-vector/harbor/base/db"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func InitPostgresql() {
	host := conf.Get("database.host")
	port := conf.Get("database.port")
	username := conf.Get("database.username")
	password := conf.Get("database.password")
	dbName := conf.Get("database.db_name")
	applicationName := conf.Get("database.application_name")
	dsn := fmt.Sprintf(
		"host=%s port=%s, user=%s password=%s dbname=%s application_name=%s",
		host, port, username, password, dbName, applicationName,
	)
	configuredLogger := logger.New(
		log.New(&lumberjack.Logger{Filename: "logs/gorm.log"}, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	pgDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: configuredLogger,
	})
	if err != nil {
		panic(err)
	}
	db.SetGlobalDb(pgDb)
}
