package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/model"
	"car_rental_miniproject/repository"

	"github.com/google/uuid"
)

var (
	ErrRentalNotFound       = errors.New("rental not found")
	ErrInsufficientDeposit  = errors.New("insufficient deposit")
	ErrRentalAlreadyExists  = errors.New("rental already exists")
	ErrUnauthorizedAccess   = errors.New("unauthorized access to this rental")
)

type RentalService interface {
	RentCar(ctx context.Context, userID uuid.UUID, req dto.RentCarRequest) (*model.RentalHistory, error)
	GetRentalByID(ctx context.Context, id uuid.UUID) (*model.RentalHistory, error)
	GetRentalsByUserID(ctx context.Context, userID uuid.UUID) ([]dto.RentalHistoryResponse, error)
	GetBookingReport(ctx context.Context, userID uuid.UUID) (*dto.BookingReportResponse, error)
	ConfirmPayment(ctx context.Context, rentalID uuid.UUID, userID uuid.UUID) error
	ConfirmExternalPayment(ctx context.Context, rentalID uuid.UUID) error
	CancelRental(ctx context.Context, rentalID uuid.UUID, userID uuid.UUID) error
	ReturnCar(ctx context.Context, rentalID uuid.UUID) error
	ProcessOverdueRentals(ctx context.Context) error
	GetAdminDashboardStats(ctx context.Context) (*model.AdminStats, []model.PopularCar, error)
}

type rentalService struct {
	pool           repository.DBPool
	rentalRepo     repository.RentalRepository
	carRepo        repository.CarRepository
	userRepo       repository.UserRepository
	paymentService PaymentService
	emailService   *EmailService
}

func NewRentalService(pool repository.DBPool, rentalRepo repository.RentalRepository, carRepo repository.CarRepository, userRepo repository.UserRepository, paymentService PaymentService, emailService *EmailService) RentalService {
	return &rentalService{
		pool:           pool,
		rentalRepo:     rentalRepo,
		carRepo:        carRepo,
		userRepo:       userRepo,
		paymentService: paymentService,
		emailService:   emailService,
	}
}

func (s *rentalService) RentCar(ctx context.Context, userID uuid.UUID, req dto.RentCarRequest) (*model.RentalHistory, error) {
	// Start transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Create repository instances with transaction
	rentalRepoTx := s.rentalRepo.WithTx(tx)
	carRepoTx := s.carRepo.WithTx(tx)
	userRepoTx := s.userRepo.WithTx(tx)

	// Get car details with lock
	car, err := carRepoTx.GetByIDForUpdate(ctx, req.CarID)
	if err != nil {
		return nil, ErrCarNotFound
	}

	// Check availability
	if !car.Availability || car.StockAvailability <= 0 {
		return nil, ErrCarNotAvailable
	}

	// Calculate total cost
	totalCost := car.RentalCosts * float64(req.RentalDays)

	// Get user details
	user, err := userRepoTx.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Note: We no longer fail here if deposit is insufficient.
	// Instead, we allow the rental to be created as 'pending' 
	// and the user can pay via the generated Xendit payment link.
	// If they want to use deposit, they can call ConfirmPayment later.

	// Decrease car stock
	if err := carRepoTx.DecreaseStock(ctx, req.CarID); err != nil {
		return nil, ErrCarNotAvailable
	}

	// Calculate return date
	returnDate := time.Now().Add(time.Duration(req.RentalDays) * 24 * time.Hour)

	// Create rental record
	rental := &model.RentalHistory{
		ID:            uuid.New(),
		UserID:        userID,
		CarID:         req.CarID,
		RentalDate:    time.Now(),
		ReturnDate:    &returnDate,
		TotalCost:     totalCost,
		Status:        "pending",
		PaymentStatus: "pending",
		PaymentURL:    "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := rentalRepoTx.Create(ctx, rental); err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Post-transaction steps (external APIs)
	
	// Generate payment invoice URL using Xendit
	orderID := fmt.Sprintf("RENTAL-%s", rental.ID.String())
	description := fmt.Sprintf("Car rental payment for %s", car.Name)
	paymentURL, err := s.paymentService.CreateInvoice(context.Background(), orderID, totalCost, user.Email, description)
	if err == nil && paymentURL != "" {
		// Update rental with payment URL (outside transaction is fine here)
		rental.PaymentURL = paymentURL
		_ = s.rentalRepo.UpdatePaymentURL(context.Background(), rental.ID, paymentURL)
	}

	// Send booking confirmation email (non-blocking)
	if s.emailService != nil && s.emailService.IsEnabled() {
		go func() {
			_ = s.emailService.SendBookingConfirmationEmail(
				context.Background(),
				user.Email,
				user.Email,
				car.Name,
				time.Now().Format("2006-01-02"),
				totalCost,
				paymentURL,
			)
		}()
	}

	return rental, nil
}

func (s *rentalService) GetRentalByID(ctx context.Context, id uuid.UUID) (*model.RentalHistory, error) {
	rental, err := s.rentalRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrRentalNotFound
	}
	return rental, nil
}

func (s *rentalService) GetRentalsByUserID(ctx context.Context, userID uuid.UUID) ([]dto.RentalHistoryResponse, error) {
	rentals, err := s.rentalRepo.GetByUserIDWithCarDetails(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.RentalHistoryResponse
	for _, r := range rentals {
		response := dto.RentalHistoryResponse{
			ID:            r.RentalHistory.ID.String(),
			UserID:        r.RentalHistory.UserID.String(),
			CarID:         r.RentalHistory.CarID.String(),
			CarName:       r.CarName,
			RentalDate:    r.RentalHistory.RentalDate.Format(time.RFC3339),
			TotalCost:     r.RentalHistory.TotalCost,
			Status:        r.RentalHistory.Status,
			PaymentStatus: r.RentalHistory.PaymentStatus,
			PaymentURL:    r.RentalHistory.PaymentURL,
		}
		if r.RentalHistory.ReturnDate != nil {
			returnDate := r.RentalHistory.ReturnDate.Format(time.RFC3339)
			response.ReturnDate = &returnDate
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *rentalService) GetBookingReport(ctx context.Context, userID uuid.UUID) (*dto.BookingReportResponse, error) {
	report, err := s.rentalRepo.GetBookingReport(ctx, userID)
	if err != nil {
		return nil, err
	}

	rentals, err := s.GetRentalsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.BookingReportResponse{
		UserID:          userID.String(),
		Email:           report.Email,
		TotalRentals:    report.TotalRentals,
		ActiveRentals:   report.ActiveRentals,
		TotalSpent:      report.TotalSpent,
		CurrentDeposit:  report.CurrentDeposit,
		RentalHistories: rentals,
	}, nil
}

func (s *rentalService) ConfirmPayment(ctx context.Context, rentalID uuid.UUID, userID uuid.UUID) error {
	// Start transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Create repository instances with transaction
	rentalRepoTx := s.rentalRepo.WithTx(tx)
	userRepoTx := s.userRepo.WithTx(tx)

	rental, err := rentalRepoTx.GetByID(ctx, rentalID)
	if err != nil {
		return ErrRentalNotFound
	}

	// Ownership check
	if rental.UserID != userID {
		return ErrUnauthorizedAccess
	}

	// Idempotency check: if already paid, return success
	if rental.PaymentStatus == "paid" {
		return nil
	}

	// Deduct deposit from user
	if err := userRepoTx.UpdateDeposit(ctx, rental.UserID, -rental.TotalCost); err != nil {
		return err
	}

	// Update rental status
	if err := rentalRepoTx.UpdateStatus(ctx, rentalID, "active", "paid"); err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit(ctx)
}

func (s *rentalService) ConfirmExternalPayment(ctx context.Context, rentalID uuid.UUID) error {
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return ErrRentalNotFound
	}

	// Idempotency check: if already paid, return success
	if rental.PaymentStatus == "paid" {
		return nil
	}

	// Update rental status WITHOUT deducting internal deposit
	// This is for external payments (Xendit) where the money is already collected
	return s.rentalRepo.UpdateStatus(ctx, rentalID, "active", "paid")
}

func (s *rentalService) CancelRental(ctx context.Context, rentalID uuid.UUID, userID uuid.UUID) error {
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return ErrRentalNotFound
	}

	// Ownership check
	if rental.UserID != userID {
		return ErrUnauthorizedAccess
	}

	// Increase car stock back
	if err := s.carRepo.IncreaseStock(ctx, rental.CarID); err != nil {
		return err
	}

	return s.rentalRepo.UpdateStatus(ctx, rentalID, "cancelled", "refunded")
}

func (s *rentalService) ReturnCar(ctx context.Context, rentalID uuid.UUID) error {
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return ErrRentalNotFound
	}

	// Update rental with return date and status
	// We update return_date to now to reflect actual return time
	// But our repository UpdateStatus only updates status and payment_status.
	// Let's assume we want to keep the original ReturnDate (due date) or update it.
	// For simplicity, let's keep the existing UpdateStatus behavior and focus on the worker.
	if err := s.rentalRepo.UpdateStatus(ctx, rentalID, "completed", "paid"); err != nil {
		return err
	}

	// Increase car stock back
	return s.carRepo.IncreaseStock(ctx, rental.CarID)
}

func (s *rentalService) ProcessOverdueRentals(ctx context.Context) error {
	overdueRentals, err := s.rentalRepo.GetOverdueRentals(ctx)
	if err != nil {
		return err
	}

	for _, r := range overdueRentals {
		// Update status to overdue
		err := s.rentalRepo.UpdateStatus(ctx, r.RentalHistory.ID, "overdue", r.RentalHistory.PaymentStatus)
		if err != nil {
			log.Printf("Failed to update rental %s to overdue: %v", r.RentalHistory.ID, err)
			continue
		}

		// Send reminder email
		user, err := s.userRepo.GetByID(ctx, r.RentalHistory.UserID)
		if err == nil && s.emailService != nil && s.emailService.IsEnabled() {
			_ = s.emailService.SendRentalReminderEmail(
				ctx,
				user.Email,
				user.Email, // using email as username for now
				r.CarName,
				r.RentalHistory.ReturnDate.Format("2006-01-02 15:04:05"),
			)
		}
	}

	return nil
}

func (s *rentalService) GetAdminDashboardStats(ctx context.Context) (*model.AdminStats, []model.PopularCar, error) {
	stats, err := s.rentalRepo.GetAdminStats(ctx)
	if err != nil {
		return nil, nil, err
	}

	popularCars, err := s.rentalRepo.GetPopularCars(ctx, 5) // Top 5 popular cars
	if err != nil {
		return nil, nil, err
	}

	return stats, popularCars, nil
}
