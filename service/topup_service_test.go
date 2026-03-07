package service

import (
	"context"
	"testing"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/model"
	"car_rental_miniproject/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTopUpRepository is a mock implementation of TopUpRepository
type MockTopUpRepository struct {
	mock.Mock
}

func (m *MockTopUpRepository) Create(ctx context.Context, transaction *model.TopUpTransaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTopUpRepository) WithTx(tx pgx.Tx) repository.TopUpRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.TopUpRepository)
}

func (m *MockTopUpRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.TopUpTransaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TopUpTransaction), args.Error(1)
}

func (m *MockTopUpRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.TopUpTransaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.TopUpTransaction), args.Error(1)
}

func (m *MockTopUpRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockTopUpRepository) UpdatePaymentURL(ctx context.Context, id uuid.UUID, paymentURL string) error {
	args := m.Called(ctx, id, paymentURL)
	return args.Error(0)
}

func (m *MockTopUpRepository) Update(ctx context.Context, transaction *model.TopUpTransaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func TestTopUpService_CreateTopUp(t *testing.T) {
	t.Run("successful top-up request", func(t *testing.T) {
		mockPool := new(MockDBPool)
		mockTopUpRepo := new(MockTopUpRepository)
		mockUserRepo := new(MockUserRepository)
		mockPaymentService := new(MockPaymentService)
		service := NewTopUpService(mockPool, mockTopUpRepo, mockUserRepo, mockPaymentService, nil)

		userID := uuid.New()
		req := dto.TopUpRequest{Amount: 500.00}
		user := &model.User{ID: userID, Email: "test@example.com"}

		mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
		mockTopUpRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.TopUpTransaction")).Return(nil)
		mockPaymentService.On("CreateInvoice", mock.Anything, mock.Anything, 500.00, user.Email, mock.Anything).Return("https://payment.url", nil)
		mockTopUpRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.TopUpTransaction")).Return(nil)

		transaction, err := service.CreateTopUp(context.Background(), userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, 500.00, transaction.Amount)
		assert.Equal(t, "pending", transaction.Status)
		assert.Equal(t, "https://payment.url", transaction.PaymentURL)
		mockUserRepo.AssertExpectations(t)
		mockTopUpRepo.AssertExpectations(t)
		mockPaymentService.AssertExpectations(t)
	})
}

func TestTopUpService_ConfirmTopUp(t *testing.T) {
	t.Run("successful top-up confirmation", func(t *testing.T) {
		mockPool := new(MockDBPool)
		mockTopUpRepo := new(MockTopUpRepository)
		mockUserRepo := new(MockUserRepository)
		service := NewTopUpService(mockPool, mockTopUpRepo, mockUserRepo, nil, nil)

		txID := uuid.New()
		userID := uuid.New()
		transaction := &model.TopUpTransaction{
			ID:     txID,
			UserID: userID,
			Amount: 500.00,
			Status: "pending",
		}

		mockTx := new(MockTx)
		mockPool.On("Begin", mock.Anything).Return(mockTx, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)
		mockTx.On("Commit", mock.Anything).Return(nil)

		mockTopUpRepo.On("WithTx", mockTx).Return(mockTopUpRepo)
		mockUserRepo.On("WithTx", mockTx).Return(mockUserRepo)

		mockTopUpRepo.On("GetByID", mock.Anything, txID).Return(transaction, nil)
		mockUserRepo.On("UpdateDeposit", mock.Anything, userID, 500.00).Return(nil)
		mockTopUpRepo.On("UpdateStatus", mock.Anything, txID, "completed").Return(nil)

		err := service.ConfirmTopUp(context.Background(), txID)

		assert.NoError(t, err)
		mockPool.AssertExpectations(t)
		mockTopUpRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("confirmation fails if already completed", func(t *testing.T) {
		mockPool := new(MockDBPool)
		mockTopUpRepo := new(MockTopUpRepository)
		mockUserRepo := new(MockUserRepository)
		service := NewTopUpService(mockPool, mockTopUpRepo, mockUserRepo, nil, nil)

		txID := uuid.New()
		transaction := &model.TopUpTransaction{ID: txID, Status: "completed"}

		mockTx := new(MockTx)
		mockPool.On("Begin", mock.Anything).Return(mockTx, nil)
		mockTx.On("Rollback", mock.Anything).Return(nil)

		mockTopUpRepo.On("WithTx", mockTx).Return(mockTopUpRepo)
		mockUserRepo.On("WithTx", mockTx).Return(mockUserRepo)

		mockTopUpRepo.On("GetByID", mock.Anything, txID).Return(transaction, nil)

		err := service.ConfirmTopUp(context.Background(), txID)

		assert.NoError(t, err) // Should return nil (idempotent)
		mockPool.AssertExpectations(t)
		mockTopUpRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}
