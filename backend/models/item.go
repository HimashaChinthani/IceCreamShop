package models

import "time"

type Item struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null;index" json:"name"`
	Category    string    `gorm:"size:50;not null" json:"category"`
	Description string    `gorm:"size:500" json:"description"`
	Quantity    int       `gorm:"not null;default:0" json:"quantity"`
	Price       float64   `gorm:"not null" json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
