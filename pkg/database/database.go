package database

import (
	"fmt"
	"log"

	"forex/models"
	"forex/pkg/config"

	postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	conf := config.GetConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", conf.PostgresHost, conf.PostgresUser, conf.PostgresPassword, conf.PostgresDbname, conf.PostgresPort, conf.PostgresSslmode, conf.PostgresTimezone)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		log.Panicln(err)
	}

	migrateDB(db)

	return db
}

func migrateDB(db *gorm.DB) {
	db.AutoMigrate(&models.Data{})
}
