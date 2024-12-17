package main

import (
	"expense-tracker/internal/api"
	"expense-tracker/internal/config"
	"expense-tracker/internal/database"
	"expense-tracker/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// Initialize API routes
	api.SetupRoutes(router, db)

	// Start server
	if err := router.Run(":" + config.Load().Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
