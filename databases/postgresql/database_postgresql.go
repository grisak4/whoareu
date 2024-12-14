package database_postgresql

import (
	"fmt"
	"log"
	"whoareu/config/confget/dbconf"
	postgresqlmodels "whoareu/models/postgresql_models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase() {
	dbconf := dbconf.GetDBConf()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		dbconf.Host, dbconf.Port, dbconf.User, dbconf.Password, dbconf.DBName)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[ERROR] %s", err.Error())
	}

	if err := db.AutoMigrate(&postgresqlmodels.User{}); err != nil {
		log.Fatalf("[ERROR] %s", err.Error())
	}

	log.Println("Connected to database!")
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	if sqlDB, err := db.DB(); err != nil {
		log.Fatalf("[ERROR] %s", err.Error())
	} else {
		sqlDB.Close()
	}
}
