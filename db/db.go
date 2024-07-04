package db

import (
	"log101/konulu-konum-backend/models"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func InitDB() {
	dbPath := os.Getenv("DB_PATH")
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to create database")
	}

	db.AutoMigrate(&models.KonuluKonum{})
}

func GetDB() *gorm.DB {
	return db
}

func SetDB(database *gorm.DB) {
	db = database
}
