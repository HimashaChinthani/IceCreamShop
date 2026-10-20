
package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FullName     string    `gorm:"size:100;not null" json:"full_name"`
	Email        string    `gorm:"size:150;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"size:20;not null" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
