package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/mail"

	"car_rental_miniproject/app/config"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// EmailService handles sending emails via Gmail API
type EmailService struct {
	gmailService *gmail.Service
	fromEmail    string
	fromName     string
	isEnabled    bool
}

// EmailMessage represents an email message to be sent
type EmailMessage struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
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

	return &EmailService{
		gmailService: gmailSvc,
		fromEmail:    fromEmail,
		fromName:     fromName,
		isEnabled:    isEnabled,
	}
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
	body := s.getWelcomeEmailTemplate(userName)

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
	body := s.getBookingConfirmationTemplate(userName, carName, rentalDate, totalCost, paymentURL)

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
	body := s.getPaymentConfirmationTemplate(userName, orderID, amount, status)

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
	body := s.getTopUpConfirmationTemplate(userName, amount, transactionID)

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
	body := s.getRentalReminderTemplate(userName, carName, returnDate)

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
	body := s.getPasswordResetTemplate(userName, resetLink)

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

// getWelcomeEmailTemplate returns the welcome email HTML template
func (s *EmailService) getWelcomeEmailTemplate(userName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background: #4CAF50; color: white; text-decoration: none; border-radius: 4px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Rental Car Service!</h1>
        </div>
        <div class="content">
            <p>Dear %s,</p>
            <p>Thank you for registering with Rental Car Service! We're excited to have you on board.</p>
            <p>You can now browse our collection of premium cars and book your perfect ride.</p>
            <p style="text-align: center; margin: 30px 0;">
                <a href="https://rentalcar.com/login" class="button">Start Browsing Cars</a>
            </p>
            <p>If you have any questions, feel free to contact our support team.</p>
            <p>Best regards,<br>The Rental Car Team</p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Rental Car Service. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, userName)
}

// getBookingConfirmationTemplate returns the booking confirmation HTML template
func (s *EmailService) getBookingConfirmationTemplate(userName, carName, rentalDate string, totalCost float64, paymentURL string) string {
	paymentButton := ""
	if paymentURL != "" {
		paymentButton = fmt.Sprintf(`
            <p style="text-align: center; margin: 30px 0;">
                <a href="%s" class="button">Complete Payment</a>
            </p>
            <p><strong>Note:</strong> Your booking will be confirmed once payment is completed.</p>`, paymentURL)
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #2196F3; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .booking-details { background: white; padding: 15px; margin: 20px 0; border-radius: 4px; }
        .button { display: inline-block; padding: 12px 24px; background: #2196F3; color: white; text-decoration: none; border-radius: 4px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Booking Confirmation</h1>
        </div>
        <div class="content">
            <p>Dear %s,</p>
            <p>Your car rental booking has been received!</p>
            
            <div class="booking-details">
                <h3>Booking Details:</h3>
                <p><strong>Car:</strong> %s</p>
                <p><strong>Rental Date:</strong> %s</p>
                <p><strong>Total Cost:</strong> $%.2f</p>
                <p><strong>Status:</strong> Pending Payment</p>
            </div>
            
            %s
            
            <p>Thank you for choosing Rental Car Service!</p>
            <p>Best regards,<br>The Rental Car Team</p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Rental Car Service. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, userName, carName, rentalDate, totalCost, paymentButton)
}

// getPaymentConfirmationTemplate returns the payment confirmation HTML template
func (s *EmailService) getPaymentConfirmationTemplate(userName, orderID string, amount float64, status string) string {
	statusColor := "#4CAF50"
	if status != "completed" && status != "settlement" {
		statusColor = "#FF9800"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: %s; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .payment-details { background: white; padding: 15px; margin: 20px 0; border-radius: 4px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Payment Confirmation</h1>
        </div>
        <div class="content">
            <p>Dear %s,</p>
            <p>Your payment has been processed.</p>
            
            <div class="payment-details">
                <h3>Payment Details:</h3>
                <p><strong>Order ID:</strong> %s</p>
                <p><strong>Amount:</strong> $%.2f</p>
                <p><strong>Status:</strong> %s</p>
            </div>
            
            <p>Thank you for your payment!</p>
            <p>Best regards,<br>The Rental Car Team</p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Rental Car Service. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, statusColor, userName, orderID, amount, status)
}

// getTopUpConfirmationTemplate returns the top-up confirmation HTML template
func (s *EmailService) getTopUpConfirmationTemplate(userName string, amount float64, transactionID string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #9C27B0; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .topup-details { background: white; padding: 15px; margin: 20px 0; border-radius: 4px; }
        .button { display: inline-block; padding: 12px 24px; background: #9C27B0; color: white; text-decoration: none; border-radius: 4px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Deposit Top-Up Confirmation</h1>
        </div>
        <div class="content">
            <p>Dear %s,</p>
            <p>Your deposit top-up request has been received.</p>
            
            <div class="topup-details">
                <h3>Top-Up Details:</h3>
                <p><strong>Transaction ID:</strong> %s</p>
                <p><strong>Amount:</strong> $%.2f</p>
                <p><strong>Status:</strong> Pending Payment</p>
            </div>
            
            <p>Your deposit balance will be updated once the payment is confirmed.</p>
            <p>Best regards,<br>The Rental Car Team</p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Rental Car Service. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, userName, transactionID, amount)
}

// getRentalReminderTemplate returns the rental reminder HTML template
func (s *EmailService) getRentalReminderTemplate(userName, carName, returnDate string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #FF9800; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .reminder-details { background: white; padding: 15px; margin: 20px 0; border-radius: 4px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Rental Reminder</h1>
        </div>
        <div class="content">
            <p>Dear %s,</p>
            <p>This is a friendly reminder about your car rental.</p>
            
            <div class="reminder-details">
                <h3>Rental Details:</h3>
                <p><strong>Car:</strong> %s</p>
                <p><strong>Expected Return Date:</strong> %s</p>
            </div>
            
            <p>Please ensure the car is returned on time to avoid any additional charges.</p>
            <p>Best regards,<br>The Rental Car Team</p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Rental Car Service. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, userName, carName, returnDate)
}

// getPasswordResetTemplate returns the password reset HTML template
func (s *EmailService) getPasswordResetTemplate(userName, resetLink string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #FF5722; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background: #FF5722; color: white; text-decoration: none; border-radius: 4px; }
        .warning { background: #fff3cd; padding: 15px; border-left: 4px solid #FF5722; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <p>Dear %s,</p>
            <p>We received a request to reset your password. Click the button below to reset it:</p>

            <p style="text-align: center; margin: 30px 0;">
                <a href="%s" class="button">Reset Password</a>
            </p>

            <p>Or copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #2196F3;">%s</p>

            <div class="warning">
                <strong>Important:</strong> This link will expire in 15 minutes.
                If you didn't request this password reset, you can safely ignore this email.
                Your password will remain unchanged.
            </div>

            <p>Best regards,<br>The Rental Car Team</p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Rental Car Service. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, userName, resetLink, resetLink)
}
