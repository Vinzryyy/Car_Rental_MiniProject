package service

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/mail"

	"car_rental_miniproject/app/config"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var emailTemplates embed.FS

// EmailService handles sending emails via Gmail API
type EmailService struct {
	gmailService *gmail.Service
	fromEmail    string
	fromName     string
	isEnabled    bool
	templates    map[string]*template.Template
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

// NewEmailService creates a new Gmail API email service
func NewEmailService(cfg *config.Config) *EmailService {
	// Check if Gmail API is configured
	apiKey := cfg.Email.GmailAPIKey
	serviceAccount := cfg.Email.GmailServiceAccountJSON

	var gmailSvc *gmail.Service
	var err error
	isEnabled := false

	if serviceAccount != "" {
		// Use service account for sending emails
		gmailSvc, err = gmail.NewService(context.Background(), option.WithCredentialsJSON([]byte(serviceAccount)))
		if err != nil {
			isEnabled = false
		} else {
			isEnabled = true
		}
	} else if apiKey != "" {
		// Use API key (limited functionality)
		gmailSvc, err = gmail.NewService(context.Background(), option.WithAPIKey(apiKey))
		if err != nil {
			isEnabled = false
		} else {
			isEnabled = true
		}
	}

	fromEmail := cfg.Email.FromEmail
	if fromEmail == "" {
		fromEmail = "noreply@rentalcar.com"
	}

	fromName := cfg.Email.FromName
	if fromName == "" {
		fromName = "Rental Car Service"
	}

	// Load email templates
	templates := loadEmailTemplates()

	return &EmailService{
		gmailService: gmailSvc,
		fromEmail:    fromEmail,
		fromName:     fromName,
		isEnabled:    isEnabled,
		templates:    templates,
	}
}

// loadEmailTemplates loads all email templates from the embedded filesystem
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
			continue // Skip if template not found
		}

		tmpl, err := template.New(file).Parse(string(content))
		if err != nil {
			continue // Skip if template parsing fails
		}

		// Store template by name (without .html extension)
		name := file[:len(file)-5]
		templates[name] = tmpl
	}

	return templates
}

// renderTemplate renders a template with the given data
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

// SendEmail sends an email using Gmail API
func (s *EmailService) SendEmail(ctx context.Context, msg EmailMessage) error {
	if !s.isEnabled {
		// Log but don't fail - email is optional
		return nil
	}

	// Create MIME message
	var contentType string
	if msg.IsHTML {
		contentType = "text/html; charset=utf-8"
	} else {
		contentType = "text/plain; charset=utf-8"
	}

	mime := fmt.Sprintf("From: %s <%s>\r\n", s.fromName, s.fromEmail)
	mime += fmt.Sprintf("To: %s\r\n", msg.To)
	mime += fmt.Sprintf("Subject: %s\r\n", msg.Subject)
	mime += fmt.Sprintf("Content-Type: %s\r\n\r\n", contentType)
	mime += msg.Body

	// Encode message in base64
	encodedMessage := base64.URLEncoding.EncodeToString([]byte(mime))

	// Create Gmail message
	gmailMsg := &gmail.Message{
		Raw: encodedMessage,
	}

	// Send via Gmail API
	_, err := s.gmailService.Users.Messages.Send("me", gmailMsg).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendWelcomeEmail sends a welcome email after user registration
func (s *EmailService) SendWelcomeEmail(ctx context.Context, toEmail, userName string) error {
	subject := "Welcome to Rental Car Service!"
	
	body, err := s.renderTemplate("welcome", WelcomeEmailData{
		UserName: userName,
	})
	if err != nil {
		return err
	}

	return s.SendEmail(ctx, EmailMessage{
		To:      toEmail,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

// SendBookingConfirmationEmail sends booking confirmation email
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

	return s.SendEmail(ctx, EmailMessage{
		To:      toEmail,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

// SendPaymentConfirmationEmail sends payment confirmation email
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

	return s.SendEmail(ctx, EmailMessage{
		To:      toEmail,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

// SendTopUpConfirmationEmail sends deposit top-up confirmation email
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

	return s.SendEmail(ctx, EmailMessage{
		To:      toEmail,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

// SendRentalReminderEmail sends rental reminder email
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

	return s.SendEmail(ctx, EmailMessage{
		To:      toEmail,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

// SendPasswordResetEmail sends password reset email
func (s *EmailService) SendPasswordResetEmail(ctx context.Context, toEmail, userName, resetLink string) error {
	subject := "Password Reset Request"
	
	body, err := s.renderTemplate("password-reset", PasswordResetData{
		UserName:  userName,
		ResetLink: resetLink,
	})
	if err != nil {
		return err
	}

	return s.SendEmail(ctx, EmailMessage{
		To:      toEmail,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

// IsEnabled returns true if the email service is configured
func (s *EmailService) IsEnabled() bool {
	return s.isEnabled
}

// GetFromEmail returns the configured from email address
func (s *EmailService) GetFromEmail() string {
	return s.fromEmail
}

// parseEmailAddress validates and parses an email address
func parseEmailAddress(email string) (*mail.Address, error) {
	return mail.ParseAddress(email)
}