package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/maxluvds/booking-platform/services/user-service/internal/auth"
	"github.com/maxluvds/booking-platform/services/user-service/internal/domain"
)

func TestUserUseCase_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("secret", 15*time.Minute, 24*time.Hour)
	uc := NewUserUseCase(mockRepo, jwtManager)

	req := &domain.RegisterRequest{
		Email:    "new@example.com",
		Password: "password123",
		Name:     "New User",
	}

	mockRepo.On("GetByEmail", mock.Anything, "new@example.com").Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(
		&domain.User{
			ID:        1,
			Email:     "new@example.com",
			Name:      "New User",
			Role:      domain.RoleUser,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		nil,
	)

	resp, err := uc.Register(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.User.ID)
	assert.Equal(t, "new@example.com", resp.User.Email)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_Register_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("secret", 15*time.Minute, 24*time.Hour)
	uc := NewUserUseCase(mockRepo, jwtManager)

	req := &domain.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Name:     "Existing User",
	}

	existingUser := &domain.User{
		ID:    1,
		Email: "existing@example.com",
		Name:  "Existing User",
		Role:  domain.RoleUser,
	}

	mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)

	resp, err := uc.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "already exists")

	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_Register_ShortPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("secret", 15*time.Minute, 24*time.Hour)
	uc := NewUserUseCase(mockRepo, jwtManager)

	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "short",
		Name:     "Test User",
	}

	resp, err := uc.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "password must be at least 8 characters")
}

func TestUserUseCase_Register_EmptyEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("secret", 15*time.Minute, 24*time.Hour)
	uc := NewUserUseCase(mockRepo, jwtManager)

	req := &domain.RegisterRequest{
		Email:    "",
		Password: "password123",
		Name:     "Test User",
	}

	resp, err := uc.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "email is required")
}

func TestUserUseCase_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("secret", 15*time.Minute, 24*time.Hour)
	uc := NewUserUseCase(mockRepo, jwtManager)

	password := "password123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: string(hash),
		Name:         "Test User",
		Role:         domain.RoleUser,
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)

	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	resp, err := uc.Login(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.User.ID)
	assert.NotEmpty(t, resp.AccessToken)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_Login_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("secret", 15*time.Minute, 24*time.Hour)
	uc := NewUserUseCase(mockRepo, jwtManager)

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)

	existingUser := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: string(hash),
		Name:         "Test User",
		Role:         domain.RoleUser,
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)

	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrong-password",
	}

	resp, err := uc.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid email or password")

	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("secret", 15*time.Minute, 24*time.Hour)
	uc := NewUserUseCase(mockRepo, jwtManager)

	mockRepo.On("GetByEmail", mock.Anything, "missing@example.com").Return(nil, nil)

	req := &domain.LoginRequest{
		Email:    "missing@example.com",
		Password: "password123",
	}

	resp, err := uc.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid email or password")

	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("secret", 15*time.Minute, 24*time.Hour)
	uc := NewUserUseCase(mockRepo, jwtManager)

	user := &domain.User{
		ID:    1,
		Email: "test@example.com",
		Name:  "Test User",
		Role:  domain.RoleUser,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(user, nil)

	result, err := uc.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "test@example.com", result.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetByID_InvalidID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("secret", 15*time.Minute, 24*time.Hour)
	uc := NewUserUseCase(mockRepo, jwtManager)

	result, err := uc.GetByID(context.Background(), 0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid user id")
}

func TestUserUseCase_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("secret", 15*time.Minute, 24*time.Hour)
	uc := NewUserUseCase(mockRepo, jwtManager)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, nil)

	result, err := uc.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")

	mockRepo.AssertExpectations(t)
}
