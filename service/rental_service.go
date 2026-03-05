package service

import (
	"context"
	"errors"
	"fmt"
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
	CancelRental(ctx context.Context, rentalID uuid.UUID, userID uuid.UUID) error
	ReturnCar(ctx context.Context, rentalID uuid.UUID) error
}

type rentalService struct {
	rentalRepo     repository.RentalRepository
	carRepo        repository.CarRepository
	userRepo       repository.UserRepository
	paymentService *XenditPaymentService
	emailService   *EmailService
}

func NewRentalService(rentalRepo repository.RentalRepository, carRepo repository.CarRepository, userRepo repository.UserRepository, paymentService *XenditPaymentService, emailService *EmailService) RentalService {
	return &rentalService{
		rentalRepo:     rentalRepo,
		carRepo:        carRepo,
		userRepo:       userRepo,
		paymentService: paymentService,
		emailService:   emailService,
	}
}

func (s *rentalService) RentCar(ctx context.Context, userID uuid.UUID, req dto.RentCarRequest) (*model.RentalHistory, error) {
	// Get car details
	car, err := s.carRepo.GetByID(ctx, req.CarID)
	if err != nil {
		return nil, ErrCarNotFound
	}

	// Check availability
	if !car.Availability || car.StockAvailability <= 0 {
		return nil, ErrCarNotAvailable
	}

	// Calculate total cost
	totalCost := car.RentalCosts * float64(req.RentalDays)

	// Check user deposit
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.DepositAmount < totalCost {
		return nil, ErrInsufficientDeposit
	}

	// Create rental record
	rental := &model.RentalHistory{
		ID:            uuid.New(),
		UserID:        userID,
		CarID:         req.CarID,
		RentalDate:    time.Now(),
		ReturnDate:    nil,
		TotalCost:     totalCost,
		Status:        "pending",
		PaymentStatus: "pending",
		PaymentURL:    "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.rentalRepo.Create(ctx, rental); err != nil {
		return nil, err
	}

	// Decrease car stock
	if err := s.carRepo.DecreaseStock(ctx, req.CarID); err != nil {
		return nil, err
	}

	// Generate payment invoice URL using Xendit
	orderID := fmt.Sprintf("RENTAL-%s-%s", rental.ID.String()[:8], time.Now().Format("20060102"))
	description := fmt.Sprintf("Car rental payment for %s", car.Name)
	paymentURL, err := s.paymentService.CreateInvoice(ctx, orderID, totalCost, user.Email, description)
	if err != nil {
		// Continue without payment URL if gateway fails
		paymentURL = ""
	}

	// Update rental with payment URL
	rental.PaymentURL = paymentURL
	if err := s.rentalRepo.UpdatePaymentURL(ctx, rental.ID, paymentURL); err != nil {
		// Log error but don't fail the request
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
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
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
	if err := s.userRepo.UpdateDeposit(ctx, rental.UserID, -rental.TotalCost); err != nil {
		return err
	}

	// Update rental status
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
	if err := s.rentalRepo.UpdateStatus(ctx, rentalID, "completed", "paid"); err != nil {
		return err
	}

	// Increase car stock back
	return s.carRepo.IncreaseStock(ctx, rental.CarID)
}
