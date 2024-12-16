package api

import (
	"net/http"
	"time"

	"expense-tracker/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetBudgets(c *gin.Context) {
	userID := c.GetUint("user_id")
	month := c.Query("month") // Get month query parameter
	var budgets []models.Budget

	query := h.db.Where("user_id = ?", userID)
	if month != "" {
		query = query.Where("month = ?", month)
	}

	if err := query.Find(&budgets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch budgets"})
		return
	}

	c.JSON(http.StatusOK, budgets)
}

func (h *Handler) CreateBudget(c *gin.Context) {
	var input struct {
		Name   string  `json:"name" binding:"required"`
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentMonth := time.Now().Format("2006-01")

	budget := models.Budget{
		UserID:         c.GetUint("user_id"),
		Name:           input.Name,
		Amount:         input.Amount,
		Month:          currentMonth,
		RollOverAmount: 0,
		OverrunAmount:  0,
	}

	if err := h.db.Create(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create budget"})
		return
	}

	c.JSON(http.StatusCreated, budget)
}

func (h *Handler) UpdateBudget(c *gin.Context) {
	userID := c.GetUint("user_id")
	budgetID := c.Param("id")

	var input struct {
		Amount         float64 `json:"amount"`
		RollOverAmount float64 `json:"roll_over_amount"`
		OverrunAmount  float64 `json:"overrun_amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var budget models.Budget
	if err := h.db.Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	budget.Amount = input.Amount
	budget.RollOverAmount = input.RollOverAmount
	budget.OverrunAmount = input.OverrunAmount

	if err := h.db.Save(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
		return
	}

	c.JSON(http.StatusOK, budget)
}

func (h *Handler) DeleteBudget(c *gin.Context) {
	userID := c.GetUint("user_id")
	budgetID := c.Param("id")

	result := h.db.Where("id = ? AND user_id = ?", budgetID, userID).Delete(&models.Budget{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete budget"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted successfully"})
}

func (h *Handler) GetBudgetOverview(c *gin.Context) {
	userID := c.GetUint("user_id")
	currentMonth := time.Now().Format("2006-01")

	var budgets []models.Budget
	if err := h.db.Where("user_id = ? AND month = ?", userID, currentMonth).Find(&budgets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch budgets"})
		return
	}

	var expenses []models.Expense
	if err := h.db.Where("user_id = ? AND DATE_TRUNC('month', date) = ?", userID, currentMonth).Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	overview := calculateBudgetOverview(budgets, expenses)
	c.JSON(http.StatusOK, overview)
}

type BudgetOverview struct {
	TotalBudget     float64                `json:"total_budget"`
	TotalExpenses   float64                `json:"total_expenses"`
	RemainingBudget float64                `json:"remaining_budget"`
	Budgets         map[string]BudgetStats `json:"budgets"`
}

type BudgetStats struct {
	Budget    float64 `json:"budget"`
	Spent     float64 `json:"spent"`
	Remaining float64 `json:"remaining"`
	Overrun   float64 `json:"overrun"`
}

func calculateBudgetOverview(budgets []models.Budget, expenses []models.Expense) BudgetOverview {
	overview := BudgetOverview{
		Budgets: make(map[string]BudgetStats),
	}

	// Initialize budgets
	for _, budget := range budgets {
		overview.TotalBudget += budget.Amount
		overview.Budgets[budget.Name] = BudgetStats{
			Budget: budget.Amount,
		}
	}

	// Calculate expenses by budget
	for _, expense := range expenses {
		overview.TotalExpenses += expense.Amount
		stats := overview.Budgets[expense.Budget.Name]
		stats.Spent += expense.Amount
		stats.Remaining = stats.Budget - stats.Spent
		if stats.Remaining < 0 {
			stats.Overrun = -stats.Remaining
			stats.Remaining = 0
		}
		overview.Budgets[expense.Budget.Name] = stats
	}

	overview.RemainingBudget = overview.TotalBudget - overview.TotalExpenses
	if overview.RemainingBudget < 0 {
		overview.RemainingBudget = 0
	}

	return overview
}
