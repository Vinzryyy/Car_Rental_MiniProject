package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"car_rental_miniproject/app/config"
	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/model"
	"car_rental_miniproject/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) WithTx(tx pgx.Tx) repository.UserRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.UserRepository)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) UpdateDeposit(ctx context.Context, id uuid.UUID, amount float64) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	args := m.Called(ctx, id, password)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockSessionRepository is a mock implementation of SessionRepository
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *model.UserSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetByToken(ctx context.Context, token string) (*model.UserSession, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserSession), args.Error(1)
}

func (m *MockSessionRepository) DeleteByToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestAuthService_Register(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	cfg := &config.JWTConfig{Secret: "testsecret", Expiration: 1}
	service := NewAuthService(mockUserRepo, mockSessionRepo, cfg, nil)

	t.Run("successful registration", func(t *testing.T) {
		req := dto.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockUserRepo.On("GetByEmail", mock.Anything, req.Email).Return(nil, errors.New("not found"))
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

		user, err := service.Register(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Email, user.Email)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("registration fails if user exists", func(t *testing.T) {
		req := dto.RegisterRequest{
			Email:    "existing@example.com",
			Password: "password123",
		}

		existingUser := &model.User{Email: req.Email}
		mockUserRepo.On("GetByEmail", mock.Anything, req.Email).Return(existingUser, nil)

		user, err := service.Register(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrUserAlreadyExists, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	cfg := &config.JWTConfig{Secret: "testsecret", Expiration: 1}
	service := NewAuthService(mockUserRepo, mockSessionRepo, cfg, nil)

	t.Run("successful login", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := &model.User{
			ID:       uuid.New(),
			Email:    "test@example.com",
			Password: string(hashedPassword),
			Role:     "user",
		}

		req := dto.LoginRequest{
			Email:    user.Email,
			Password: password,
		}

		mockUserRepo.On("GetByEmail", mock.Anything, req.Email).Return(user, nil)
		mockSessionRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.UserSession")).Return(nil)

		token, resp, err := service.Login(context.Background(), req)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, resp)
		assert.Equal(t, user.Email, resp.User.Email)
		mockUserRepo.AssertExpectations(t)
		mockSessionRepo.AssertExpectations(t)
	})

	t.Run("login fails with invalid credentials", func(t *testing.T) {
		req := dto.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		mockUserRepo.On("GetByEmail", mock.Anything, req.Email).Return(nil, errors.New("not found"))

		token, resp, err := service.Login(context.Background(), req)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInvalidCredentials, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	cfg := &config.JWTConfig{Secret: "testsecret", Expiration: 1}
	service := NewAuthService(mockUserRepo, mockSessionRepo, cfg, nil)

	t.Run("successful validation", func(t *testing.T) {
		userID := uuid.New()
		email := "test@example.com"
		role := "user"
		
		// Generate a real token for testing
		s := service.(*authService)
		token, _ := s.generateToken(userID, email, role)

		session := &model.UserSession{
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour),
		}

		mockSessionRepo.On("GetByToken", mock.Anything, token).Return(session, nil)

		valUserID, valRole, err := service.ValidateToken(context.Background(), token)

		assert.NoError(t, err)
		assert.Equal(t, userID, valUserID)
		assert.Equal(t, role, valRole)
		mockSessionRepo.AssertExpectations(t)
	})

	t.Run("validation fails if session not found", func(t *testing.T) {
		token := "invalidtoken"
		mockSessionRepo.On("GetByToken", mock.Anything, token).Return(nil, errors.New("not found"))

		valUserID, valRole, err := service.ValidateToken(context.Background(), token)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, valUserID)
		assert.Empty(t, valRole)
		assert.Equal(t, ErrInvalidToken, err)
		mockSessionRepo.AssertExpectations(t)
	})
}
