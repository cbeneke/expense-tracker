package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"expense-tracker/internal/auth"
	"expense-tracker/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the test database
	err = db.AutoMigrate(&models.User{}, &models.Budget{}, &models.Expense{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	SetupRoutes(router, db)
	return router
}

func setupTestUser(t *testing.T, db *gorm.DB) *models.User {
	hashedPassword, _ := auth.HashPassword("password123")
	user := &models.User{
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}
	result := db.Create(user)
	assert.NoError(t, result.Error)
	return user
}

// Auth Handler Tests
func TestSignUp(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)

	tests := []struct {
		name       string
		input      map[string]string
		wantStatus int
	}{
		{
			name: "valid signup",
			input: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid email",
			input: map[string]string{
				"email":    "invalid-email",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusCreated {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
			}
		})
	}
}

// Budget Handler Tests
func TestCreateBudget(t *testing.T) {
	db := setupTestDB(t)
	user := setupTestUser(t, db)

	// Create auth token
	token, err := auth.GenerateToken(user.ID)
	assert.NoError(t, err)

	router := setupTestRouter(db)

	tests := []struct {
		name       string
		input      map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid budget",
			input: map[string]interface{}{
				"name":   "Groceries",
				"amount": 500.00,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid amount",
			input: map[string]interface{}{
				"name":   "Groceries",
				"amount": "invalid",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/budgets", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

// Expense Handler Tests
func TestCreateExpense(t *testing.T) {
	db := setupTestDB(t)
	user := setupTestUser(t, db)

	// Create auth token
	token, err := auth.GenerateToken(user.ID)
	assert.NoError(t, err)

	router := setupTestRouter(db)

	// Create a test budget
	budget := &models.Budget{
		UserID: user.ID,
		Name:   "Groceries",
		Amount: 500.00,
		Month:  time.Now().Format("2006-01"),
	}
	db.Create(budget)

	tests := []struct {
		name       string
		input      map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid expense",
			input: map[string]interface{}{
				"amount":      50.00,
				"budget_id":   budget.ID,
				"description": "Weekly groceries",
				"date":        time.Now().Format("2006-01-02"),
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid amount",
			input: map[string]interface{}{
				"amount":      "invalid",
				"budget_id":   budget.ID,
				"description": "Weekly groceries",
				"date":        time.Now().Format("2006-01-02"),
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/expenses", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

// Helper function to create a test token
func createTestToken(userID uint) (string, error) {
	// You'll need to implement this based on your auth package
	// Return a valid JWT token for testing
	return "test-token", nil
}

// Add more test functions for other handlers...
