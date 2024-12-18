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
		Name     string  `json:"name" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
		RollOver *bool   `json:"roll_over" binding:"required"` // https://github.com/gin-gonic/gin/issues/685
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	budget := models.Budget{
		UserID:   int64(c.GetUint("user_id")),
		Name:     input.Name,
		Amount:   input.Amount,
		RollOver: *input.RollOver,
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
		Name   string  `json:"name"`
		Amount float64 `json:"amount"`
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

	// Start a transaction
	tx := h.db.Begin()

	// Update the main budget
	budget.Name = input.Name
	budget.Amount = input.Amount

	if err := tx.Save(&budget).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
		return
	}

	// Update month budgets
	currentTime := time.Now()
	query := tx.Model(&models.MonthBudget{}).Where("budget_id = ? AND month >= ?", budget.ID, currentTime.Format("2006-01"))

	// Update the amount for all relevant month budgets
	if err := query.Updates(map[string]interface{}{
		"amount": budget.Amount,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update month budgets"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
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
