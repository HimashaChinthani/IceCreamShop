
package database

import (
	"log"

	"icecream-shop-api/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("icecream.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Customer{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("SQLite database connected successfully")
	log.Println("Database tables created successfully")

	return db
}
