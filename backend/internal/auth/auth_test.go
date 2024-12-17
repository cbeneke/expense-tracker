package auth

import (
	"net/http/httptest"
	"testing"

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

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "valid user",
			email:    "test@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "duplicate email",
			email:    "test@example.com",
			password: "password123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := CreateUser(db, tt.email, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.email, user.Email)
		})
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	userID := uint(1)

	token, err := GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Create a mock gin context with the token
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	// Validate the token
	extractedUserID, err := ValidateToken(c)
	assert.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)
}
