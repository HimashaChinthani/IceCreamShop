package models

import "time"

type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;index" json:"name"`
	Email     string    `gorm:"size:150;uniqueIndex" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Address   string    `gorm:"size:300" json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy *uint     `json:"created_by" gorm:"index"`
	Creator   *User     `json:"-" gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
