package api

import (
	"net/http"
	"time"

	"expense-tracker/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetExpenses(c *gin.Context) {
	userID := c.GetUint("user_id")
	month := c.Query("month") // Get month query parameter
	var expenses []models.Expense

	query := h.db.Where("user_id = ?", userID)
	if month != "" {
		// If month is provided, filter expenses for that month
		query = query.Where("DATE_TRUNC('month', date)::date = ?", month+"-01")
	}

	if err := query.
		Preload("Budget").
		Order("date DESC").
		Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

func (h *Handler) CreateExpense(c *gin.Context) {
	var input struct {
		Amount      float64 `json:"amount" binding:"required"`
		BudgetID    uint    `json:"budget_id" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Date        string  `json:"date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Get the user ID from the context
	userID := c.GetUint("user_id")

	budget := &models.Budget{}
	if err := h.db.Where("id = ? AND user_id = ?", input.BudgetID, userID).First(budget).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Budget not found"})
		return
	}

	monthBudget, err := budget.GetOrCreateMonthBudget(h.db, date.Format("2006-01"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get or create month budget"})
		return
	}

	// Verify the expense date matches the budget month
	expenseMonth := date.Format("2006-01")
	if expenseMonth != monthBudget.Month {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expense date must be in the same month as the budget"})
		return
	}

	// Create the expense
	expense := models.Expense{
		UserID:        userID,
		MonthBudgetID: input.BudgetID,
		Amount:        input.Amount,
		Description:   input.Description,
		Date:          date,
	}

	if err := h.db.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	// Update the month budget's used amount
	monthBudget.UsedAmount += input.Amount
	if err := h.db.Save(monthBudget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update month budget"})
		return
	}

	// Reconcile rollover for historic budgets
	if monthBudget.Month != time.Now().Format("2006-01") {
		if err := budget.ReconcileRollover(h.db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reconcile rollover"})
			return
		}
	}

	c.JSON(http.StatusCreated, expense)
}

func (h *Handler) UpdateExpense(c *gin.Context) {
	userID := c.GetUint("user_id")
	expenseID := c.Param("id")

	var input struct {
		Amount      float64 `json:"amount"`
		BudgetID    *uint   `json:"budget_id"`
		Description string  `json:"description"`
		Date        string  `json:"date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expense models.Expense
	if err := h.db.Where("id = ? AND user_id = ?", expenseID, userID).First(&expense).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	// If budget is being changed, update both old and new budgets
	if input.BudgetID != nil {
		tx := h.db.Begin()

		newBudget := &models.Budget{}
		if err := tx.Where("id = ? AND user_id = ?", input.BudgetID, userID).First(newBudget).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Budget not found"})
			tx.Rollback()
			return
		}

		expense.MonthBudget.UsedAmount -= expense.Amount
		if err := tx.Save(expense.MonthBudget).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update old month budget"})
			tx.Rollback()
			return
		}

		// Date changes are handled later in the code
		newMonthBudget, err := newBudget.GetOrCreateMonthBudget(h.db, expense.Date.Format("2006-01"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get or create month budget"})
			tx.Rollback()
			return
		}
		newMonthBudget.UsedAmount += expense.Amount
		if err := tx.Save(newMonthBudget).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update old month budget"})
			tx.Rollback()
			return
		}

		// Reconcile rollover for historic budgets
		if expense.MonthBudget.Month != time.Now().Format("2006-01") {
			if err := expense.MonthBudget.Budget.ReconcileRollover(h.db); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reconcile rollover"})
				tx.Rollback()
				return
			}
			if err := newBudget.ReconcileRollover(h.db); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reconcile rollover"})
				tx.Rollback()
				return
			}
		}

		expense.MonthBudget = *newMonthBudget
		tx.Save(&expense)

		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			tx.Rollback()
			return
		}
	}

	if input.Amount != 0 {
		// Update budget used amount
		expense.MonthBudget.UsedAmount += input.Amount - expense.Amount
		if err := h.db.Save(&expense.MonthBudget).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
			return
		}
		expense.Amount = input.Amount
		if err := h.db.Save(&expense).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
			return
		}
	}

	if input.Description != "" {
		expense.Description = input.Description
		if err := h.db.Save(&expense).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
			return
		}
	}

	if input.Date != "" {
		date, err := time.Parse("2006-01-02", input.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}

		if expense.MonthBudget.Month != date.Format("2006-01") {
			tx := h.db.Begin()
			newMonthBudget, err := expense.MonthBudget.Budget.GetOrCreateMonthBudget(h.db, date.Format("2006-01"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get or create month budget"})
				tx.Rollback()
				return
			}
			expense.MonthBudget.UsedAmount -= expense.Amount
			if err := tx.Save(&expense.MonthBudget).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update old month budget"})
				tx.Rollback()
				return
			}
			newMonthBudget.UsedAmount += expense.Amount
			if err := tx.Save(&newMonthBudget).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update new month budget"})
				tx.Rollback()
				return
			}
			expense.MonthBudget = *newMonthBudget

			// Reconcile rollover for the budget
			if err := expense.MonthBudget.Budget.ReconcileRollover(h.db); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reconcile rollover"})
				tx.Rollback()
				return
			}

			tx.Save(&expense)
			if err := tx.Commit().Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
				tx.Rollback()
				return
			}
		}
		expense.Date = date
		if err := h.db.Save(&expense).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
			return
		}
	}

	// Load the budget relationship for the response
	if err := h.db.Model(&expense).Association("MonthBudget.Budget").Find(&expense.MonthBudget.Budget); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load budget association"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

func (h *Handler) DeleteExpense(c *gin.Context) {
	userID := c.GetUint("user_id")
	expenseID := c.Param("id")

	var expense models.Expense
	if err := h.db.Where("id = ? AND user_id = ?", expenseID, userID).First(&expense).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	expense.MonthBudget.UsedAmount -= expense.Amount
	if err := h.db.Save(&expense.MonthBudget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update month budget"})
		return
	}

	// Reconcile rollover for historic budgets
	if expense.MonthBudget.Month != time.Now().Format("2006-01") {
		if err := expense.MonthBudget.Budget.ReconcileRollover(h.db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reconcile rollover"})
			return
		}
	}

	if err := h.db.Delete(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
