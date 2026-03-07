package service

import (
	"context"
	"testing"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/model"
	"car_rental_miniproject/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (m *MockRentalRepository) WithTx(tx pgx.Tx) repository.RentalRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.RentalRepository)
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

func (m *MockRentalRepository) GetOverdueRentals(ctx context.Context) ([]model.RentalWithCarDetails, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RentalWithCarDetails), args.Error(1)
}

func (m *MockRentalRepository) GetAdminStats(ctx context.Context) (*model.AdminStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AdminStats), args.Error(1)
}

func (m *MockRentalRepository) GetPopularCars(ctx context.Context, limit int) ([]model.PopularCar, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PopularCar), args.Error(1)
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

// MockDBPool is a mock implementation of DBPool
type MockDBPool struct {
	mock.Mock
}

func (m *MockDBPool) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockDBPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (m *MockDBPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	callArgs := m.Called(ctx, sql, args)
	return callArgs.Get(0).(pgx.Rows), callArgs.Error(1)
}

func (m *MockDBPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	callArgs := m.Called(ctx, sql, args)
	return callArgs.Get(0).(pgx.Row)
}

func (m *MockDBPool) Close() {
	m.Called()
}

func (m *MockDBPool) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockTx is a mock implementation of pgx.Tx
type MockTx struct {
	mock.Mock
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	args := m.Called(ctx, tableName, columnNames, rowSrc)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	args := m.Called(ctx, b)
	return args.Get(0).(pgx.BatchResults)
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	args := m.Called()
	return args.Get(0).(pgx.LargeObjects)
}

func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	args := m.Called(ctx, name, sql)
	return args.Get(0).(*pgconn.StatementDescription), args.Error(1)
}

func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	callArgs := m.Called(ctx, sql, args)
	return callArgs.Get(0).(pgx.Rows), callArgs.Error(1)
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	callArgs := m.Called(ctx, sql, args)
	return callArgs.Get(0).(pgx.Row)
}

func (m *MockTx) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}

func TestRentalService_RentCar(t *testing.T) {
	t.Run("successful car rental", func(t *testing.T) {
		mockPool := new(MockDBPool)
		mockRentalRepo := new(MockRentalRepository)
		mockCarRepo := new(MockCarRepository)
		mockUserRepo := new(MockUserRepository)
		mockPaymentService := new(MockPaymentService)
		service := NewRentalService(mockPool, mockRentalRepo, mockCarRepo, mockUserRepo, mockPaymentService, nil)

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

		mockTx := new(MockTx)
		mockPool.On("Begin", mock.Anything).Return(mockTx, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Commit", mock.Anything).Return(nil)

		mockCarRepo.On("WithTx", mockTx).Return(mockCarRepo)
		mockRentalRepo.On("WithTx", mockTx).Return(mockRentalRepo)
		mockUserRepo.On("WithTx", mockTx).Return(mockUserRepo)

		mockCarRepo.On("GetByIDForUpdate", mock.Anything, carID).Return(car, nil)
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
		mockPool.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockRentalRepo.AssertExpectations(t)
		mockPaymentService.AssertExpectations(t)
	})

	t.Run("rental fails if car not available", func(t *testing.T) {
		mockPool := new(MockDBPool)
		mockRentalRepo := new(MockRentalRepository)
		mockCarRepo := new(MockCarRepository)
		mockUserRepo := new(MockUserRepository)
		mockPaymentService := new(MockPaymentService)
		service := NewRentalService(mockPool, mockRentalRepo, mockCarRepo, mockUserRepo, mockPaymentService, nil)

		userID := uuid.New()
		carID := uuid.New()
		req := dto.RentCarRequest{CarID: carID, RentalDays: 1}

		car := &model.Car{ID: carID, Availability: false, StockAvailability: 0}

		mockTx := new(MockTx)
		mockPool.On("Begin", mock.Anything).Return(mockTx, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)

		mockCarRepo.On("WithTx", mockTx).Return(mockCarRepo)
		mockRentalRepo.On("WithTx", mockTx).Return(mockRentalRepo)
		mockUserRepo.On("WithTx", mockTx).Return(mockUserRepo)

		mockCarRepo.On("GetByIDForUpdate", mock.Anything, carID).Return(car, nil)

		rental, err := service.RentCar(context.Background(), userID, req)

		assert.Error(t, err)
		assert.Nil(t, rental)
		assert.Equal(t, ErrCarNotAvailable, err)
		mockPool.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})

	t.Run("successful rental even with insufficient deposit", func(t *testing.T) {
		mockPool := new(MockDBPool)
		mockRentalRepo := new(MockRentalRepository)
		mockCarRepo := new(MockCarRepository)
		mockUserRepo := new(MockUserRepository)
		mockPaymentService := new(MockPaymentService)
		service := NewRentalService(mockPool, mockRentalRepo, mockCarRepo, mockUserRepo, mockPaymentService, nil)

		userID := uuid.New()
		carID := uuid.New()
		req := dto.RentCarRequest{CarID: carID, RentalDays: 5}

		car := &model.Car{ID: carID, Name: "Expensive Car", Availability: true, StockAvailability: 5, RentalCosts: 100.00}
		user := &model.User{ID: userID, Email: "poor@example.com", DepositAmount: 50.00}

		mockTx := new(MockTx)
		mockPool.On("Begin", mock.Anything).Return(mockTx, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Commit", mock.Anything).Return(nil)

		mockCarRepo.On("WithTx", mockTx).Return(mockCarRepo)
		mockRentalRepo.On("WithTx", mockTx).Return(mockRentalRepo)
		mockUserRepo.On("WithTx", mockTx).Return(mockUserRepo)

		mockCarRepo.On("GetByIDForUpdate", mock.Anything, carID).Return(car, nil)
		mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
		mockCarRepo.On("DecreaseStock", mock.Anything, carID).Return(nil)
		mockRentalRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.RentalHistory")).Return(nil)
		mockPaymentService.On("CreateInvoice", mock.Anything, mock.Anything, 500.00, user.Email, mock.Anything).Return("https://payment.url", nil)
		mockRentalRepo.On("UpdatePaymentURL", mock.Anything, mock.Anything, "https://payment.url").Return(nil)

		rental, err := service.RentCar(context.Background(), userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, rental)
		assert.Equal(t, 500.00, rental.TotalCost)
		assert.Equal(t, "pending", rental.Status)
		mockPool.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}
