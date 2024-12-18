package models

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	MonthBudgetID uint           `gorm:"not null" json:"month_budget_id"`
	Amount        float64        `gorm:"not null" json:"amount"`
	Description   string         `gorm:"not null" json:"description"`
	Date          time.Time      `gorm:"not null" json:"date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	User          User           `gorm:"foreignKey:UserID" json:"-"`
	MonthBudget   MonthBudget    `gorm:"foreignKey:MonthBudgetID" json:"month_budget"`
}
