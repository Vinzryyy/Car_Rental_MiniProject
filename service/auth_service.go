package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"car_rental_miniproject/app/config"
	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/model"
	"car_rental_miniproject/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidOldPassword = errors.New("invalid old password")
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req dto.LoginRequest) (string, *dto.LoginResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, *dto.LoginResponse, error)
	Logout(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	UpdateProfile(ctx context.Context, userID uuid.UUID, email string) (*model.User, error)
}

type authService struct {
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	cfg          *config.JWTConfig
	emailService *EmailService
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, cfg *config.JWTConfig, emailService *EmailService) AuthService {
	return &authService{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		cfg:          cfg,
		emailService: emailService,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:            uuid.New(),
		Email:         req.Email,
		Password:      string(hashedPassword),
		DepositAmount: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Send welcome email (non-blocking)
	if s.emailService != nil && s.emailService.IsEnabled() {
		go func(email string) {
			_ = s.emailService.SendWelcomeEmail(context.Background(), email, email)
		}(user.Email)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (string, *dto.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return "", nil, err
	}

	// Create session in database
	session := &model.UserSession{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(s.cfg.Expiration)),
		CreatedAt: time.Now(),
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", nil, fmt.Errorf("failed to create session: %w", err)
	}

	response := &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			DepositAmount: user.DepositAmount,
		},
		ExpiresIn: s.cfg.Expiration,
	}

	return token, response, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (uuid.UUID, error) {
	// Check if token exists in database (login state)
	session, err := s.sessionRepo.GetByToken(ctx, tokenString)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		_ = s.sessionRepo.DeleteByToken(ctx, tokenString)
		return uuid.Nil, ErrTokenExpired
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.Secret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token claims")
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *authService) generateToken(userID uuid.UUID, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.cfg.Expiration)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

// RefreshToken generates a new access token using a valid refresh token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, *dto.LoginResponse, error) {
	// Validate the refresh token
	userID, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return "", nil, ErrInvalidToken
	}

	// Get user from database
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", nil, ErrUserNotFound
	}

	// Generate new access token
	newToken, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return "", nil, err
	}

	// Create new session
	session := &model.UserSession{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     newToken,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(s.cfg.Expiration)),
		CreatedAt: time.Now(),
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", nil, err
	}

	// Delete old session
	_ = s.sessionRepo.DeleteByToken(ctx, refreshToken)

	response := &dto.LoginResponse{
		Token: newToken,
		User: dto.UserResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			DepositAmount: user.DepositAmount,
		},
		ExpiresIn: s.cfg.Expiration,
	}

	return newToken, response, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.DeleteByToken(ctx, token)
}

// ForgotPassword initiates password reset by generating a reset token and sending email
func (s *authService) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not for security reasons
		// Return success anyway to prevent email enumeration
		return "password reset token (if email exists)", nil
	}

	// Generate reset token (JWT with short expiration)
	resetToken, err := s.generateResetToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	// Send password reset email (non-blocking)
	if s.emailService != nil && s.emailService.IsEnabled() {
		go func(email, token string) {
			resetLink := fmt.Sprintf("https://yourapp.com/reset-password?token=%s", token)
			_ = s.emailService.SendPasswordResetEmail(context.Background(), email, email, resetLink)
		}(user.Email, resetToken)
	}

	return "password reset token sent", nil
}

// ResetPassword resets the user password using a valid reset token
func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate reset token
	claims, err := s.validateResetToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return ErrInvalidToken
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password in database
	return s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}

// ChangePassword changes the password for an authenticated user
func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidOldPassword
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password in database
	return s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}

// UpdateProfile updates user profile information
func (s *authService) UpdateProfile(ctx context.Context, userID uuid.UUID, email string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Update email if provided
	if email != "" && email != user.Email {
		// Check if email already exists
		_, err := s.userRepo.GetByEmail(ctx, email)
		if err == nil {
			return nil, errors.New("email already in use")
		}
		user.Email = email
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// generateResetToken generates a short-lived JWT for password reset
func (s *authService) generateResetToken(userID uuid.UUID, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID.String(),
		"email":     email,
		"purpose":   "password_reset",
		"exp":       time.Now().Add(15 * time.Minute).Unix(), // Short expiration
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

// validateResetToken validates a password reset token
func (s *authService) validateResetToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, ErrTokenExpired
		}
	}

	// Check purpose
	if purpose, ok := claims["purpose"].(string); !ok || purpose != "password_reset" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
