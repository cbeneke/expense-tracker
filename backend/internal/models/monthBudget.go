package models

import (
	"time"

	"gorm.io/gorm"
)

// MonthBudget represents the monthly instance of a budget
type MonthBudget struct {
	ID             int64     `json:"id" db:"id"`
	BudgetID       int64     `json:"budget_id" db:"budget_id"`
	Month          string    `json:"month" db:"month"`
	Amount         float64   `json:"amount" db:"amount"`
	UsedAmount     float64   `json:"used_amount" db:"used_amount"`
	RollOverAmount *float64  `json:"roll_over_amount,omitempty" db:"roll_over_amount"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	Budget         Budget    `gorm:"foreignKey:BudgetID" json:"-"`
}

func (mb *MonthBudget) ReconcileIsActive(db *gorm.DB) error {
	// If rollover amount is not null, then the month budget is active
	if mb.RollOverAmount != nil {
		return nil
	}

	lastMonthBudget := MonthBudget{}
	if err := db.Where("budget_id = ? AND month < ?", mb.BudgetID, mb.Month).Order("month desc").First(&lastMonthBudget).Error; err != nil {
		// Initialize rollover amount if it's the first month
		if err == gorm.ErrRecordNotFound {
			newRollOverAmount := 0.0
			mb.RollOverAmount = &newRollOverAmount
			return db.Save(&mb).Error
		}
		return err
	}

	tx := db.Begin()
	// Calculate the rollover amount for the current month
	newRollOverAmount := lastMonthBudget.Amount - lastMonthBudget.UsedAmount + *lastMonthBudget.RollOverAmount

	// Get all months between last month budget and current month budget
	lastMonthBudgetDate, err := time.Parse("2006-01", lastMonthBudget.Month)
	if err != nil {
		return err
	}
	monthBudgetDate, err := time.Parse("2006-01", mb.Month)
	if err != nil {
		return err
	}
	monthDiff := int(monthBudgetDate.Sub(lastMonthBudgetDate).Hours() / 24 / 30)
	if monthDiff > 1 {
		// Create month budgets for missing months
		for i := 1; i < monthDiff; i++ {
			intermediateMonth := lastMonthBudgetDate.AddDate(0, i, 0)
			intermediateBudget := MonthBudget{
				BudgetID:       mb.BudgetID,
				Month:          intermediateMonth.Format("2006-01"),
				Amount:         mb.Amount, // The current saved month is the latest version of the budget
				UsedAmount:     0,
				RollOverAmount: &newRollOverAmount,
			}
			if err := tx.Create(&intermediateBudget).Error; err != nil {
				tx.Rollback()
				return err
			}
			// For every month without a monthBudget there was no expense, add the full amount to the rollover amount
			newRollOverAmount += lastMonthBudget.Amount
		}
	}

	mb.RollOverAmount = &newRollOverAmount
	tx.Save(&mb)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
