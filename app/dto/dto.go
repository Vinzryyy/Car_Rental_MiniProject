package dto

import "github.com/google/uuid"

// FieldError represents a single field validation error
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// RegisterRequest represents the registration request payload
// @Description User registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest represents the login request payload
// @Description User login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents the refresh token request
// @Description Request to refresh JWT token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ForgotPasswordRequest represents the forgot password request
// @Description Request to reset password
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents the reset password request
// @Description Request to set new password with reset token
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ChangePasswordRequest represents the change password request for authenticated users
// @Description Request to change password for authenticated user
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UpdateProfileRequest represents the profile update request
// @Description Request to update user profile
type UpdateProfileRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
}

// TopUpRequest represents the deposit top-up request
// @Description Request to add funds to user deposit
type TopUpRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// RentCarRequest represents the car rental request
// @Description Request to rent a car
type RentCarRequest struct {
	CarID      uuid.UUID `json:"car_id" validate:"required,uuid"`
	RentalDays int       `json:"rental_days" validate:"required,gt=0"`
}

// UpdateCarRequest represents the car update request
// @Description Request to update car information
type UpdateCarRequest struct {
	Name              string  `json:"name" validate:"omitempty"`
	StockAvailability int     `json:"stock_availability" validate:"omitempty,gte=0"`
	RentalCosts       float64 `json:"rental_costs" validate:"omitempty,gt=0"`
	Category          string  `json:"category" validate:"omitempty"`
	Description       string  `json:"description" validate:"omitempty"`
	ImageURL          string  `json:"image_url" validate:"omitempty,url"`
}

// CreateCarRequest represents the car creation request
// @Description Request to create a new car
type CreateCarRequest struct {
	Name              string  `json:"name" validate:"required"`
	StockAvailability int     `json:"stock_availability" validate:"required,gte=0"`
	RentalCosts       float64 `json:"rental_costs" validate:"required,gt=0"`
	Category          string  `json:"category" validate:"required"`
	Description       string  `json:"description" validate:"required"`
	ImageURL          string  `json:"image_url" validate:"omitempty,url"`
}

// APIResponse represents a standard API response
// @Description Standard API response wrapper
type APIResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
	Errors  []FieldError  `json:"errors,omitempty"`
}

// LoginResponse represents the login response with JWT token
// @Description Login response containing JWT token and user info
type LoginResponse struct {
	Token     string      `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresIn int         `json:"expires_in"`
}

// UserResponse represents the user data in responses
// @Description User information returned in API responses
type UserResponse struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	DepositAmount float64 `json:"deposit_amount"`
	Role          string  `json:"role"`
}

// CarResponse represents the car data in responses
// @Description Car information returned in API responses
type CarResponse struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Availability      bool    `json:"availability"`
	StockAvailability int     `json:"stock_availability"`
	RentalCosts       float64 `json:"rental_costs"`
	Category          string  `json:"category"`
	Description       string  `json:"description"`
	ImageURL          string  `json:"image_url"`
}

// RentalHistoryResponse represents rental history in responses
// @Description Rental history information returned in API responses
type RentalHistoryResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	CarID         string  `json:"car_id"`
	CarName       string  `json:"car_name"`
	RentalDate    string  `json:"rental_date"`
	ReturnDate    *string `json:"return_date"`
	TotalCost     float64 `json:"total_cost"`
	Status        string  `json:"status"`
	PaymentStatus string  `json:"payment_status"`
	PaymentURL    string  `json:"payment_url"`
}

// BookingReportResponse represents the booking report
// @Description Comprehensive booking report for user
type BookingReportResponse struct {
	UserID          string                 `json:"user_id"`
	Email           string                 `json:"email"`
	TotalRentals    int                    `json:"total_rentals"`
	ActiveRentals   int                    `json:"active_rentals"`
	TotalSpent      float64                `json:"total_spent"`
	CurrentDeposit  float64                `json:"current_deposit"`
	RentalHistories []RentalHistoryResponse `json:"rental_histories"`
}

// PaymentResponse represents payment gateway response
// @Description Payment information from payment gateway
type PaymentResponse struct {
	PaymentURL    string `json:"payment_url"`
	PaymentID     string `json:"payment_id"`
	Status        string `json:"status"`
	RedirectURL   string `json:"redirect_url"`
}
