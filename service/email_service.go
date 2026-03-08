package service

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"car_rental_miniproject/app/config"
)

//go:embed templates/emails/*
var emailTemplates embed.FS

// EmailService handles sending emails via Resend API
type EmailService struct {
	apiKey    string
	fromEmail string
	fromName  string
	isEnabled bool
	templates map[string]*template.Template
	client    *http.Client
}

// ResendRequest represents the request body for Resend API
type ResendRequest struct {
	From    string `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
}

// EmailMessage represents an email message to be sent
type EmailMessage struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

// Template data structures
type WelcomeEmailData struct {
	UserName string
}

type BookingConfirmationData struct {
	UserName    string
	CarName     string
	RentalDate  string
	TotalCost   float64
	PaymentURL  string
}

type PaymentConfirmationData struct {
	UserName    string
	OrderID     string
	Amount      float64
	Status      string
	StatusColor string
}

type TopUpConfirmationData struct {
	UserName      string
	Amount        float64
	TransactionID string
}

type RentalReminderData struct {
	UserName   string
	CarName    string
	ReturnDate string
}

type PasswordResetData struct {
	UserName  string
	ResetLink string
}

// NewEmailService creates a new Resend API email service
func NewEmailService(cfg *config.Config) *EmailService {
	// Check if Email is enabled in config
	if !cfg.Email.IsEnabled {
		log.Println("Email service is explicitly disabled in config")
		return &EmailService{
			isEnabled: false,
			fromEmail: cfg.Email.FromEmail,
			fromName:  cfg.Email.FromName,
			templates: loadEmailTemplates(),
		}
	}

	apiKey := cfg.Email.ResendAPIKey
	if apiKey == "" {
		log.Println("Email enabled but RESEND_API_KEY not found. Service will be disabled.")
		return &EmailService{
			isEnabled: false,
			fromEmail: cfg.Email.FromEmail,
			fromName:  cfg.Email.FromName,
			templates: loadEmailTemplates(),
		}
	}

	fromEmail := cfg.Email.FromEmail
	if fromEmail == "" {
		fromEmail = "onboarding@resend.dev" // Default Resend test email
	}

	fromName := cfg.Email.FromName
	if fromName == "" {
		fromName = "Rental Car Service"
	}

	log.Println("Email Service (Resend) initialized successfully")

	return &EmailService{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
		isEnabled: true,
		templates: loadEmailTemplates(),
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

// loadEmailTemplates remains the same...
func loadEmailTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template)
	
	templateFiles := []string{
		"welcome.html",
		"booking-confirmation.html",
		"payment-confirmation.html",
		"topup-confirmation.html",
		"rental-reminder.html",
		"password-reset.html",
	}

	for _, file := range templateFiles {
		content, err := emailTemplates.ReadFile("templates/emails/" + file)
		if err != nil {
			continue
		}

		tmpl, err := template.New(file).Parse(string(content))
		if err != nil {
			continue
		}

		name := file[:len(file)-5]
		templates[name] = tmpl
	}

	return templates
}

func (s *EmailService) renderTemplate(name string, data interface{}) (string, error) {
	tmpl, ok := s.templates[name]
	if !ok {
		return "", fmt.Errorf("template %s not found", name)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SendEmail sends an email using Resend API
func (s *EmailService) SendEmail(ctx context.Context, msg EmailMessage) error {
	if !s.isEnabled {
		log.Printf("Email service disabled. Skip sending to %s", msg.To)
		return nil
	}

	resendReq := ResendRequest{
		From:    fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail),
		To:      []string{msg.To},
		Subject: msg.Subject,
		Html:    msg.Body,
	}

	body, err := json.Marshal(resendReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.resend.com/emails", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("Failed to call Resend API: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		log.Printf("Resend API error (Status %d): %v", resp.StatusCode, errResp)
		return fmt.Errorf("resend api error: status %d", resp.StatusCode)
	}

	log.Printf("Email sent successfully to %s via Resend", msg.To)
	return nil
}

// SendWelcomeEmail remains the same...
func (s *EmailService) SendWelcomeEmail(ctx context.Context, toEmail, userName string) error {
	subject := "Welcome to Rental Car Service!"
	body, err := s.renderTemplate("welcome", WelcomeEmailData{UserName: userName})
	if err != nil {
		return err
	}
	return s.SendEmail(ctx, EmailMessage{To: toEmail, Subject: subject, Body: body, IsHTML: true})
}

// SendBookingConfirmationEmail remains the same...
func (s *EmailService) SendBookingConfirmationEmail(ctx context.Context, toEmail, userName, carName, rentalDate string, totalCost float64, paymentURL string) error {
	subject := "Booking Confirmation - " + carName
	body, err := s.renderTemplate("booking-confirmation", BookingConfirmationData{
		UserName:   userName,
		CarName:    carName,
		RentalDate: rentalDate,
		TotalCost:  totalCost,
		PaymentURL: paymentURL,
	})
	if err != nil {
		return err
	}
	return s.SendEmail(ctx, EmailMessage{To: toEmail, Subject: subject, Body: body, IsHTML: true})
}

// SendPaymentConfirmationEmail remains the same...
func (s *EmailService) SendPaymentConfirmationEmail(ctx context.Context, toEmail, userName, orderID string, amount float64, status string) error {
	subject := "Payment Confirmation - " + orderID
	statusColor := "#4CAF50"
	if status != "completed" && status != "settlement" {
		statusColor = "#FF9800"
	}
	body, err := s.renderTemplate("payment-confirmation", PaymentConfirmationData{
		UserName:    userName,
		OrderID:     orderID,
		Amount:      amount,
		Status:      status,
		StatusColor: statusColor,
	})
	if err != nil {
		return err
	}
	return s.SendEmail(ctx, EmailMessage{To: toEmail, Subject: subject, Body: body, IsHTML: true})
}

func (s *EmailService) SendTopUpConfirmationEmail(ctx context.Context, toEmail, userName string, amount float64, transactionID string) error {
	subject := "Deposit Top-Up Confirmation"
	body, err := s.renderTemplate("topup-confirmation", TopUpConfirmationData{
		UserName:      userName,
		Amount:        amount,
		TransactionID: transactionID,
	})
	if err != nil {
		return err
	}
	return s.SendEmail(ctx, EmailMessage{To: toEmail, Subject: subject, Body: body, IsHTML: true})
}

func (s *EmailService) SendRentalReminderEmail(ctx context.Context, toEmail, userName, carName, returnDate string) error {
	subject := "Rental Reminder - Return Your " + carName
	body, err := s.renderTemplate("rental-reminder", RentalReminderData{
		UserName:   userName,
		CarName:    carName,
		ReturnDate: returnDate,
	})
	if err != nil {
		return err
	}
	return s.SendEmail(ctx, EmailMessage{To: toEmail, Subject: subject, Body: body, IsHTML: true})
}

func (s *EmailService) SendPasswordResetEmail(ctx context.Context, toEmail, userName, resetLink string) error {
	subject := "Password Reset Request"
	body, err := s.renderTemplate("password-reset", PasswordResetData{
		UserName:  userName,
		ResetLink: resetLink,
	})
	if err != nil {
		return err
	}
	return s.SendEmail(ctx, EmailMessage{To: toEmail, Subject: subject, Body: body, IsHTML: true})
}

func (s *EmailService) IsEnabled() bool {
	return s.isEnabled
}

func (s *EmailService) GetFromEmail() string {
	return s.fromEmail
}
