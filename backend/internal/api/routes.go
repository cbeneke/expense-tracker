package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	handler := NewHandler(db)

	// Auth routes (no middleware)
	router.POST("/auth/login", handler.Login)
	router.POST("/auth/signup", handler.SignUp)
	router.GET("/auth/validate", handler.AuthMiddleware(), handler.ValidateToken)

	// Protected routes
	api := router.Group("/api")
	api.Use(handler.AuthMiddleware())
	{
		// Budget routes
		api.GET("/budgets", handler.GetBudgets)
		api.POST("/budgets", handler.CreateBudget)
		api.PUT("/budgets/:id", handler.UpdateBudget)
		api.DELETE("/budgets/:id", handler.DeleteBudget)
		api.GET("/budgets/overview", handler.GetBudgetOverview)

		// Expense routes
		api.GET("/expenses", handler.GetExpenses)
		api.POST("/expenses", handler.CreateExpense)
		api.PUT("/expenses/:id", handler.UpdateExpense)
		api.DELETE("/expenses/:id", handler.DeleteExpense)
	}
}
