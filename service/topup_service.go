package service

import (
	"context"
	"fmt"
	"time"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/model"
	"car_rental_miniproject/repository"

	"github.com/google/uuid"
)

type TopUpService interface {
	CreateTopUp(ctx context.Context, userID uuid.UUID, req dto.TopUpRequest) (*model.TopUpTransaction, error)
	GetTopUpByID(ctx context.Context, id uuid.UUID) (*model.TopUpTransaction, error)
	GetTopUpsByUserID(ctx context.Context, userID uuid.UUID) ([]model.TopUpTransaction, error)
	ConfirmTopUp(ctx context.Context, transactionID uuid.UUID) error
	CancelTopUp(ctx context.Context, transactionID uuid.UUID) error
}

type topUpService struct {
	topUpRepo        repository.TopUpRepository
	userRepo         repository.UserRepository
	paymentService   *XenditPaymentService
	emailService     *EmailService
}

func NewTopUpService(topUpRepo repository.TopUpRepository, userRepo repository.UserRepository, paymentService *XenditPaymentService, emailService *EmailService) TopUpService {
	return &topUpService{
		topUpRepo:      topUpRepo,
		userRepo:       userRepo,
		paymentService: paymentService,
		emailService:   emailService,
	}
}

func (s *topUpService) CreateTopUp(ctx context.Context, userID uuid.UUID, req dto.TopUpRequest) (*model.TopUpTransaction, error) {
	// Get user email for payment
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	transaction := &model.TopUpTransaction{
		ID:            uuid.New(),
		UserID:        userID,
		Amount:        req.Amount,
		Status:        "pending",
		PaymentMethod: "xendit",
		PaymentURL:    "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.topUpRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	// Generate payment invoice URL using Xendit
	orderID := fmt.Sprintf("TOPUP-%s", transaction.ID.String())
	description := fmt.Sprintf("Deposit top-up for account %s", user.Email)
	paymentURL, err := s.paymentService.CreateInvoice(ctx, orderID, req.Amount, user.Email, description)
	if err != nil {
		// Continue without payment URL if gateway fails
		paymentURL = ""
	}

	// Update transaction with payment URL
	transaction.PaymentURL = paymentURL
	if err := s.topUpRepo.Update(ctx, transaction); err != nil {
		return nil, err
	}

	// Send top-up confirmation email (non-blocking)
	if s.emailService != nil && s.emailService.IsEnabled() {
		go func() {
			_ = s.emailService.SendTopUpConfirmationEmail(
				context.Background(),
				user.Email,
				user.Email,
				req.Amount,
				transaction.ID.String()[:8],
			)
		}()
	}

	return transaction, nil
}

func (s *topUpService) GetTopUpByID(ctx context.Context, id uuid.UUID) (*model.TopUpTransaction, error) {
	return s.topUpRepo.GetByID(ctx, id)
}

func (s *topUpService) GetTopUpsByUserID(ctx context.Context, userID uuid.UUID) ([]model.TopUpTransaction, error) {
	return s.topUpRepo.GetByUserID(ctx, userID)
}

func (s *topUpService) ConfirmTopUp(ctx context.Context, transactionID uuid.UUID) error {
	transaction, err := s.topUpRepo.GetByID(ctx, transactionID)
	if err != nil {
		return err
	}

	// Idempotency check: if already completed or cancelled, return success
	if transaction.Status == "completed" {
		return nil
	}
	
	if transaction.Status == "cancelled" {
		return fmt.Errorf("cannot confirm a cancelled transaction")
	}

	// Update user deposit
	if err := s.userRepo.UpdateDeposit(ctx, transaction.UserID, transaction.Amount); err != nil {
		return err
	}

	// Update transaction status
	return s.topUpRepo.UpdateStatus(ctx, transactionID, "completed")
}

func (s *topUpService) CancelTopUp(ctx context.Context, transactionID uuid.UUID) error {
	return s.topUpRepo.UpdateStatus(ctx, transactionID, "cancelled")
}
