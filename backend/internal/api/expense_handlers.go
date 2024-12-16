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
		BudgetID    *uint   `json:"budget_id"`
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

	var budget *models.Budget
	if input.BudgetID != nil {
		budget = &models.Budget{}
		if err := h.db.Where("id = ? AND user_id = ?", input.BudgetID, userID).First(budget).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Budget not found"})
			return
		}

		// Verify the expense date matches the budget month
		expenseMonth := date.Format("2006-01")
		if expenseMonth != budget.Month {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Expense date must be in the same month as the budget"})
			return
		}
	}

	// Create the expense
	expense := models.Expense{
		UserID:      userID,
		BudgetID:    input.BudgetID,
		Amount:      input.Amount,
		Category:    budget.Category,
		Description: input.Description,
		Date:        date,
	}

	if err := h.db.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	if budget != nil {
		// Update the budget's roll-over amount
		budget.RollOverAmount += input.Amount
		if err := h.db.Save(budget).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
			return
		}
	}

	// Load the budget relationship for the response
	h.db.Model(&expense).Association("Budget").Find(&expense.Budget)

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
	if input.BudgetID != nil && (expense.BudgetID == nil || *input.BudgetID != *expense.BudgetID) {
		// Update old budget
		if expense.BudgetID != nil {
			var oldBudget models.Budget
			if err := h.db.First(&oldBudget, expense.BudgetID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find old budget"})
				return
			}
			oldBudget.RollOverAmount -= expense.Amount
			if err := h.db.Save(&oldBudget).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update old budget"})
				return
			}
		}

		// Update new budget
		var newBudget models.Budget
		if err := h.db.First(&newBudget, input.BudgetID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find new budget"})
			return
		}
		newBudget.RollOverAmount += input.Amount
		if err := h.db.Save(&newBudget).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update new budget"})
			return
		}

		expense.BudgetID = input.BudgetID
	}

	if input.Amount != 0 {
		// Update budget amount
		var budget models.Budget
		if err := h.db.First(&budget, expense.BudgetID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find budget"})
			return
		}
		budget.RollOverAmount = budget.RollOverAmount - expense.Amount + input.Amount
		if err := h.db.Save(&budget).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
			return
		}
		expense.Amount = input.Amount
	}

	if input.Description != "" {
		expense.Description = input.Description
	}

	if input.Date != "" {
		date, err := time.Parse("2006-01-02", input.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		expense.Date = date
	}

	if err := h.db.Save(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}

	// Load the budget relationship for the response
	h.db.Model(&expense).Association("Budget").Find(&expense.Budget)

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

	// Update the budget's roll-over amount
	var budget models.Budget
	if err := h.db.First(&budget, expense.BudgetID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find budget"})
		return
	}

	budget.RollOverAmount -= expense.Amount
	if err := h.db.Save(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
		return
	}

	if err := h.db.Delete(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
