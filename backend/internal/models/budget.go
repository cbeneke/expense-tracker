package models

import (
	"time"

	"gorm.io/gorm"
)

// Budget represents the high-level budget configuration
type Budget struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Amount    float64   `json:"amount" db:"amount"`
	UserID    int64     `json:"user_id" db:"user_id"`
	RollOver  bool      `json:"roll_over" db:"roll_over"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GetOrCreateMonthBudget creates a new monthly budget if it doesn't exist
func (b *Budget) GetOrCreateMonthBudget(db *gorm.DB, month string) (*MonthBudget, error) {
	// First try to get existing month budget
	mb := &MonthBudget{}
	result := db.Where("budget_id = ? AND month = ?", b.ID, month).First(mb)

	if result.Error == nil {
		return mb, nil
	}

	if result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	// Initialize new month budget
	mb = &MonthBudget{
		BudgetID:       b.ID,
		Month:          month,
		Amount:         b.Amount,
		UsedAmount:     0,
		RollOverAmount: nil,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := db.Create(mb).Error; err != nil {
		return nil, err
	}

	return mb, nil
}

func (b *Budget) ReconcileRollover(db *gorm.DB) error {
	if !b.RollOver {
		return nil
	}

	var monthBudgets []MonthBudget
	if err := db.Where("budget_id = ? AND roll_over_amount IS NOT NULL", b.ID).Order("month asc").Find(&monthBudgets).Error; err != nil {
		return err
	}
	tx := db.Begin()

	var totalAmount, totalUsed float64 = 0, 0
	for _, mb := range monthBudgets {
		// Update the month budget with new rollover amount
		if err := tx.Model(&mb).Update("roll_over_amount", totalAmount-totalUsed).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Calculate the rolling amounts
		totalAmount += mb.Amount
		totalUsed += mb.UsedAmount
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
