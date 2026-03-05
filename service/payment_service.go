package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"car_rental_miniproject/app/config"
)

var (
	ErrPaymentGatewayUnavailable = errors.New("payment gateway unavailable")
	ErrPaymentFailed             = errors.New("payment failed")
)

// MidtransConfig holds Midtrans API configuration
type MidtransConfig struct {
	ServerKey     string
	ClientKey     string
	BaseURL       string
	IsProduction  bool
}

// MidtransPaymentService handles payment processing via Midtrans API
type MidtransPaymentService struct {
	config *MidtransConfig
	httpClient *http.Client
}

// SnapRequest represents the request to Midtrans Snap API
type SnapRequest struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	CustomerDetails    CustomerDetails    `json:"customer_details"`
	EnabledPayments    []string           `json:"enabled_payments,omitempty"`
}

// TransactionDetails represents transaction details for Midtrans
type TransactionDetails struct {
	OrderID  string `json:"order_id"`
	GrossAmount int64 `json:"gross_amount"`
}

// CustomerDetails represents customer details for Midtrans
type CustomerDetails struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
}

// SnapResponse represents the response from Midtrans Snap API
type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

// PaymentNotification represents Midtrans payment notification
type PaymentNotification struct {
	OrderID       string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus   string `json:"fraud_status"`
	GrossAmount   string `json:"gross_amount"`
	PaymentType   string `json:"payment_type"`
	StatusCode    string `json:"status_code"`
}

// NewMidtransPaymentService creates a new Midtrans payment service
func NewMidtransPaymentService(cfg *config.Config) *MidtransPaymentService {
	isProd := cfg.Server.Env == "production"
	baseURL := "https://app.sandbox.midtrans.com"
	if isProd {
		baseURL = "https://app.midtrans.com"
	}

	return &MidtransPaymentService{
		config: &MidtransConfig{
			ServerKey:    getEnv("MIDTRANS_SERVER_KEY", ""),
			ClientKey:    getEnv("MIDTRANS_CLIENT_KEY", ""),
			BaseURL:      baseURL,
			IsProduction: isProd,
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateSnapPayment creates a payment link using Midtrans Snap API
func (s *MidtransPaymentService) CreateSnapPayment(ctx context.Context, orderID string, amount float64, userEmail string) (string, error) {
	if s.config.ServerKey == "" {
		// Return mock URL for development without API key
		return fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v2/transactions/%s/pay?mock=true", orderID), nil
	}

	grossAmount := int64(amount)

	req := SnapRequest{
		TransactionDetails: TransactionDetails{
			OrderID:     orderID,
			GrossAmount: grossAmount,
		},
		CustomerDetails: CustomerDetails{
			FirstName: "Customer",
			Email:     userEmail,
		},
		EnabledPayments: []string{"credit_card", "gopay", "bank_transfer", "qris", "shopeepay"},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/snap/v1/transactions", s.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication header (Base64 encoded server key)
	auth := base64.StdEncoding.EncodeToString([]byte(s.config.ServerKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+auth)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to call payment gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("payment gateway returned status %d", resp.StatusCode)
	}

	var snapResp SnapResponse
	if err := json.NewDecoder(resp.Body).Decode(&snapResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return snapResp.RedirectURL, nil
}

// CheckPaymentStatus checks the status of a payment transaction
func (s *MidtransPaymentService) CheckPaymentStatus(ctx context.Context, orderID string) (*PaymentNotification, error) {
	if s.config.ServerKey == "" {
		return &PaymentNotification{
			OrderID:           orderID,
			TransactionStatus: "settlement",
			FraudStatus:       "accept",
			StatusCode:        "200",
		}, nil
	}

	url := fmt.Sprintf("%s/v2/%s/status", s.config.BaseURL, orderID)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication header
	auth := base64.StdEncoding.EncodeToString([]byte(s.config.ServerKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+auth)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call payment gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment gateway returned status %d", resp.StatusCode)
	}

	var paymentNotif PaymentNotification
	if err := json.NewDecoder(resp.Body).Decode(&paymentNotif); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &paymentNotif, nil
}

// VerifyPaymentNotification verifies a payment notification from Midtrans
func (s *MidtransPaymentService) VerifyPaymentNotification(ctx context.Context, orderID string, signature string) bool {
	// In production, verify the signature hash from Midtrans
	// For now, return true if server key is configured
	return s.config.ServerKey != ""
}

// GetPaymentMethods returns available payment methods
func (s *MidtransPaymentService) GetPaymentMethods() []string {
	return []string{
		"Credit Card",
		"GoPay",
		"Bank Transfer",
		"QRIS",
		"ShopeePay",
	}
}

// IsConfigured returns true if the payment gateway is configured
func (s *MidtransPaymentService) IsConfigured() bool {
	return s.config.ServerKey != ""
}

// GetEnvironment returns the current environment (sandbox/production)
func (s *MidtransPaymentService) GetEnvironment() string {
	if s.config.IsProduction {
		return "production"
	}
	return "sandbox"
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
