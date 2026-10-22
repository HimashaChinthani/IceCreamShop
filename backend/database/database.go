package database

import (
	"log"

	"icecream-shop-api/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("icecream.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Inventory{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("SQLite database connected successfully")
	log.Println("Database tables created successfully")

	DB = db
	return db
}
