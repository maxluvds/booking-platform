package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateAndValidateAccessToken(t *testing.T) {
	secret := "test-secret-key"
	accessTTL := 15 * time.Minute
	refreshTTL := 24 * time.Hour
	manager := NewJWTManager(secret, accessTTL, refreshTTL)

	userID := int64(42)
	email := "test@example.com"
	role := "user"

	token, err := manager.GenerateAccessToken(userID, email, role)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := manager.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
}

func TestJWTManager_GenerateAndValidateRefreshToken(t *testing.T) {
	secret := "test-secret-key"
	manager := NewJWTManager(secret, 15*time.Minute, 24*time.Hour)

	userID := int64(1)
	email := "refresh@example.com"
	role := "admin"

	token, err := manager.GenerateRefreshToken(userID, email, role)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := manager.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
}

func TestJWTManager_ValidateInvalidToken(t *testing.T) {
	manager := NewJWTManager("secret", 15*time.Minute, 24*time.Hour)

	_, err := manager.ValidateToken("invalid.token.here")
	assert.Error(t, err)
}

func TestJWTManager_ValidateTokenWithWrongSecret(t *testing.T) {
	manager1 := NewJWTManager("secret-1", 15*time.Minute, 24*time.Hour)
	manager2 := NewJWTManager("secret-2", 15*time.Minute, 24*time.Hour)

	token, err := manager1.GenerateAccessToken(1, "test@example.com", "user")
	require.NoError(t, err)

	_, err = manager2.ValidateToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestJWTManager_ExpiredToken(t *testing.T) {
	manager := NewJWTManager("secret", 1*time.Millisecond, 1*time.Millisecond)

	token, err := manager.GenerateAccessToken(1, "expired@example.com", "user")
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	_, err = manager.ValidateToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}
