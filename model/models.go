package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents the users entity
// @Description User entity with authentication and deposit information
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	Password     string     `json:"-" db:"password"`
	DepositAmount float64   `json:"deposit_amount" db:"deposit_amount"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// Car represents the cars entity (main rental entity)
// @Description Car entity for rental with availability and pricing
type Car struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	Name               string     `json:"name" db:"name"`
	Availability       bool       `json:"availability" db:"availability"`
	StockAvailability  int        `json:"stock_availability" db:"stock_availability"`
	RentalCosts        float64    `json:"rental_costs" db:"rental_costs"`
	Category           string     `json:"category" db:"category"`
	Description        string     `json:"description" db:"description"`
	ImageURL           string     `json:"image_url" db:"image_url"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// RentalHistory represents the rental history table
// @Description Rental history tracking user car rentals
type RentalHistory struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	CarID           uuid.UUID  `json:"car_id" db:"car_id"`
	RentalDate      time.Time  `json:"rental_date" db:"rental_date"`
	ReturnDate      *time.Time `json:"return_date" db:"return_date"`
	TotalCost       float64    `json:"total_cost" db:"total_cost"`
	Status          string     `json:"status" db:"status"` // pending, active, completed, cancelled
	PaymentStatus   string     `json:"payment_status" db:"payment_status"` // pending, paid, failed
	PaymentURL      string     `json:"payment_url" db:"payment_url"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// TopUpTransaction represents deposit top-up transactions
// @Description Transaction record for user deposit top-ups
type TopUpTransaction struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	Amount        float64    `json:"amount" db:"amount"`
	Status        string     `json:"status" db:"status"` // pending, completed, failed
	PaymentMethod string     `json:"payment_method" db:"payment_method"`
	PaymentURL    string     `json:"payment_url" db:"payment_url"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// UserSession represents a login session
type UserSession struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
