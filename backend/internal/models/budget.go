package models

import (
	"time"

	"gorm.io/gorm"
)

type Budget struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"not null" json:"user_id"`
	Name           string         `gorm:"not null" json:"name"`
	Amount         float64        `gorm:"not null" json:"amount"`
	Month          string         `gorm:"not null" json:"month"` // Format: "2024-01"
	RollOverAmount float64        `json:"roll_over_amount"`
	OverrunAmount  float64        `json:"overrun_amount"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	User           User           `gorm:"foreignKey:UserID" json:"-"`
}
