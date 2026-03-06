package service

import (
	"context"
	"testing"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRentalRepository is a mock implementation of RentalRepository
type MockRentalRepository struct {
	mock.Mock
}

func (m *MockRentalRepository) Create(ctx context.Context, rental *model.RentalHistory) error {
	args := m.Called(ctx, rental)
	return args.Error(0)
}

func (m *MockRentalRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.RentalHistory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RentalHistory), args.Error(1)
}

func (m *MockRentalRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.RentalHistory, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RentalHistory), args.Error(1)
}

func (m *MockRentalRepository) GetByUserIDWithCarDetails(ctx context.Context, userID uuid.UUID) ([]model.RentalWithCarDetails, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RentalWithCarDetails), args.Error(1)
}

func (m *MockRentalRepository) GetBookingReport(ctx context.Context, userID uuid.UUID) (*model.BookingReport, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BookingReport), args.Error(1)
}

func (m *MockRentalRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status, paymentStatus string) error {
	args := m.Called(ctx, id, status, paymentStatus)
	return args.Error(0)
}

func (m *MockRentalRepository) UpdatePaymentURL(ctx context.Context, id uuid.UUID, paymentURL string) error {
	args := m.Called(ctx, id, paymentURL)
	return args.Error(0)
}

// MockPaymentService is a mock implementation of PaymentService
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) CreateInvoice(ctx context.Context, orderID string, amount float64, userEmail string, description string) (string, error) {
	args := m.Called(ctx, orderID, amount, userEmail, description)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentService) CheckPaymentStatus(ctx context.Context, orderID string) (*PaymentNotification, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PaymentNotification), args.Error(1)
}

func (m *MockPaymentService) VerifyPaymentNotification(ctx context.Context, orderID string, callbackToken string) bool {
	args := m.Called(ctx, orderID, callbackToken)
	return args.Bool(0)
}

func (m *MockPaymentService) GetPaymentMethods() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockPaymentService) IsConfigured() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockPaymentService) GetEnvironment() string {
	args := m.Called()
	return args.String(0)
}

func TestRentalService_RentCar(t *testing.T) {
	mockRentalRepo := new(MockRentalRepository)
	mockCarRepo := new(MockCarRepository)
	mockUserRepo := new(MockUserRepository)
	mockPaymentService := new(MockPaymentService)
	
	service := NewRentalService(mockRentalRepo, mockCarRepo, mockUserRepo, mockPaymentService, nil)

	t.Run("successful car rental", func(t *testing.T) {
		userID := uuid.New()
		carID := uuid.New()
		req := dto.RentCarRequest{
			CarID:      carID,
			RentalDays: 2,
		}

		car := &model.Car{
			ID:                carID,
			Name:              "Test Car",
			Availability:      true,
			StockAvailability: 5,
			RentalCosts:       100.00,
		}

		user := &model.User{
			ID:            userID,
			Email:         "test@example.com",
			DepositAmount: 500.00,
		}

		mockCarRepo.On("GetByID", mock.Anything, carID).Return(car, nil)
		mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
		mockCarRepo.On("DecreaseStock", mock.Anything, carID).Return(nil)
		mockRentalRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.RentalHistory")).Return(nil)
		mockPaymentService.On("CreateInvoice", mock.Anything, mock.Anything, 200.00, user.Email, mock.Anything).Return("https://payment.url", nil)
		mockRentalRepo.On("UpdatePaymentURL", mock.Anything, mock.Anything, "https://payment.url").Return(nil)

		rental, err := service.RentCar(context.Background(), userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, rental)
		assert.Equal(t, 200.00, rental.TotalCost)
		assert.Equal(t, "https://payment.url", rental.PaymentURL)
		mockCarRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockRentalRepo.AssertExpectations(t)
		mockPaymentService.AssertExpectations(t)
	})

	t.Run("rental fails if car not available", func(t *testing.T) {
		userID := uuid.New()
		carID := uuid.New()
		req := dto.RentCarRequest{CarID: carID, RentalDays: 1}

		car := &model.Car{ID: carID, Availability: false, StockAvailability: 0}

		mockCarRepo.On("GetByID", mock.Anything, carID).Return(car, nil)

		rental, err := service.RentCar(context.Background(), userID, req)

		assert.Error(t, err)
		assert.Nil(t, rental)
		assert.Equal(t, ErrCarNotAvailable, err)
		mockCarRepo.AssertExpectations(t)
	})

	t.Run("rental fails if insufficient deposit", func(t *testing.T) {
		userID := uuid.New()
		carID := uuid.New()
		req := dto.RentCarRequest{CarID: carID, RentalDays: 5}

		car := &model.Car{ID: carID, Availability: true, StockAvailability: 5, RentalCosts: 100.00}
		user := &model.User{ID: userID, DepositAmount: 50.00}

		mockCarRepo.On("GetByID", mock.Anything, carID).Return(car, nil)
		mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

		rental, err := service.RentCar(context.Background(), userID, req)

		assert.Error(t, err)
		assert.Nil(t, rental)
		assert.Equal(t, ErrInsufficientDeposit, err)
		mockCarRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}
