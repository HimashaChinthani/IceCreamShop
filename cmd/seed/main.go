package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"icecream-shop-api/database"
	"icecream-shop-api/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SeedUser struct {
	FullName    string
	Email       string
	Role        string
	Password    string
	PasswordEnv string
}

func main() {
	// Passwords must be supplied through environment variables.
	users := []SeedUser{
		{
			FullName:    "Shop Owner",
			Email:       "owner@icecream.local",
			Role:        "OWNER",
			Password:    os.Getenv("SEED_OWNER_PASSWORD"),
			PasswordEnv: "SEED_OWNER_PASSWORD",
		},
		{
			FullName:    "Shop Manager",
			Email:       "manager@icecream.local",
			Role:        "MANAGER",
			Password:    os.Getenv("SEED_MANAGER_PASSWORD"),
			PasswordEnv: "SEED_MANAGER_PASSWORD",
		},
		{
			FullName:    "Shop Cashier",
			Email:       "cashier@icecream.local",
			Role:        "CASHIER",
			Password:    os.Getenv("SEED_CASHIER_PASSWORD"),
			PasswordEnv: "SEED_CASHIER_PASSWORD",
		},
	}

	// Validate all passwords before changing the database.
	for _, user := range users {
		if len(user.Password) < 12 {
			log.Fatalf(
				"Password for role %s must contain at least 12 bytes; set %s before running the seed command",
				user.Role,
				user.PasswordEnv,
			)
		}
	}

	db := database.Connect()

	for _, user := range users {
		if err := seedUser(db, user); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("User seeding completed successfully.")
}

func seedUser(db *gorm.DB, user SeedUser) error {
	var existing models.User

	err := db.Where("email = ?", user.Email).First(&existing).Error

	if err == nil {
		log.Printf(
			"User already exists: %s (%s). Skipping.",
			user.Email,
			user.Role,
		)
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf(
			"failed to check user %s: %w",
			user.Email,
			err,
		)
	}

	// Hash the password before storing it.
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to hash password for %s: %w",
			user.Email,
			err,
		)
	}

	newUser := models.User{
		FullName:     user.FullName,
		Email:        user.Email,
		PasswordHash: string(hash),
		Role:         user.Role,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return fmt.Errorf(
			"failed to create user %s: %w",
			user.Email,
			err,
		)
	}

	log.Printf(
		"Created user: %s (%s)",
		user.Email,
		user.Role,
	)

	return nil
}
