package jwt

import (
	"testing"
	"time"

	"os"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTManager(t *testing.T) {
	// Generate keys for testing if they don't exist
	// In this project they already exist in the backend folder
	privateKey, _ := os.ReadFile("../../private.pem")
	publicKey, _ := os.ReadFile("../../public.pem")

	manager, err := NewJWTManager(string(privateKey), string(publicKey), 15*time.Minute, 7*24*time.Hour)
	assert.NoError(t, err)

	userID := uuid.New()
	email := "test@example.com"

	t.Run("Generate and Verify Access Token", func(t *testing.T) {
		token, err := manager.GenerateAccessToken(userID, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := manager.VerifyToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, "access", claims.TokenType)
	})

	t.Run("Generate and Verify Refresh Token", func(t *testing.T) {
		token, err := manager.GenerateRefreshToken(userID, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := manager.VerifyToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, "refresh", claims.TokenType)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, err := manager.VerifyToken("invalid.token.here")
		assert.Error(t, err)
	})
}
