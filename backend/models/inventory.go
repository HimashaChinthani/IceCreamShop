package models

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	SKU       string         `json:"sku" gorm:"uniqueIndex;not null"`
	Quantity  int            `json:"quantity" gorm:"not null;default:0"`
	Price     float64        `json:"price" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	CreatedBy *uint          `json:"created_by" gorm:"index"`
	Creator   *User          `json:"-" gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
