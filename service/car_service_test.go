package service

import (
	"context"
	"testing"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/model"
	"car_rental_miniproject/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCarRepository is a mock implementation of CarRepository
type MockCarRepository struct {
	mock.Mock
}

func (m *MockCarRepository) Create(ctx context.Context, car *model.Car) error {
	args := m.Called(ctx, car)
	return args.Error(0)
}

func (m *MockCarRepository) GetAll(ctx context.Context, filter repository.CarFilter) ([]model.Car, int, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Car), args.Int(1), args.Error(2)
}

func (m *MockCarRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Car, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Car), args.Error(1)
}

func (m *MockCarRepository) Update(ctx context.Context, car *model.Car) error {
	args := m.Called(ctx, car)
	return args.Error(0)
}

func (m *MockCarRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCarRepository) DecreaseStock(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCarRepository) IncreaseStock(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCarService_CreateCar(t *testing.T) {
	mockRepo := new(MockCarRepository)
	service := NewCarService(mockRepo)

	t.Run("successful car creation", func(t *testing.T) {
		req := dto.CreateCarRequest{
			Name:              "Test Car",
			StockAvailability: 5,
			RentalCosts:       100.00,
			Category:          "Sedan",
			Description:       "A test car",
			ImageURL:          "https://example.com/car.jpg",
		}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Car")).Return(nil)

		car, err := service.CreateCar(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, car)
		assert.Equal(t, req.Name, car.Name)
		assert.Equal(t, req.Category, car.Category)
		assert.True(t, car.Availability)
		mockRepo.AssertExpectations(t)
	})

	t.Run("car with zero stock is not available", func(t *testing.T) {
		req := dto.CreateCarRequest{
			Name:              "Test Car",
			StockAvailability: 0,
			RentalCosts:       100.00,
			Category:          "Sedan",
			Description:       "A test car",
		}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Car")).Return(nil)

		car, err := service.CreateCar(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, car)
		assert.False(t, car.Availability)
		mockRepo.AssertExpectations(t)
	})
}

func TestCarService_GetAllCars(t *testing.T) {
	mockRepo := new(MockCarRepository)
	service := NewCarService(mockRepo)

	t.Run("get all cars without filter", func(t *testing.T) {
		cars := []model.Car{
			{ID: uuid.New(), Name: "Car 1", Category: "Sedan", Availability: true, StockAvailability: 5},
			{ID: uuid.New(), Name: "Car 2", Category: "SUV", Availability: true, StockAvailability: 3},
		}
		filter := repository.CarFilter{}

		mockRepo.On("GetAll", mock.Anything, filter).Return(cars, 2, nil)

		result, total, err := service.GetAllCars(context.Background(), filter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get cars by category", func(t *testing.T) {
		cars := []model.Car{
			{ID: uuid.New(), Name: "Car 1", Category: "Sedan", Availability: true, StockAvailability: 5},
		}
		filter := repository.CarFilter{Category: "Sedan"}

		mockRepo.On("GetAll", mock.Anything, filter).Return(cars, 1, nil)

		result, total, err := service.GetAllCars(context.Background(), filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, 1, total)
		assert.Equal(t, "Sedan", result[0].Category)
		mockRepo.AssertExpectations(t)
	})
}

func TestCarService_GetCarByID(t *testing.T) {
	mockRepo := new(MockCarRepository)
	service := NewCarService(mockRepo)

	t.Run("successful get car by ID", func(t *testing.T) {
		carID := uuid.New()
		car := &model.Car{
			ID:                carID,
			Name:              "Test Car",
			Category:          "Sedan",
			Availability:      true,
			StockAvailability: 5,
			RentalCosts:       100.00,
		}

		mockRepo.On("GetByID", mock.Anything, carID).Return(car, nil)

		result, err := service.GetCarByID(context.Background(), carID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, carID, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("car not found", func(t *testing.T) {
		carID := uuid.New()

		mockRepo.On("GetByID", mock.Anything, carID).Return((*model.Car)(nil), ErrCarNotFound)

		result, err := service.GetCarByID(context.Background(), carID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrCarNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCarService_UpdateCar(t *testing.T) {
	mockRepo := new(MockCarRepository)
	service := NewCarService(mockRepo)

	t.Run("successful update", func(t *testing.T) {
		carID := uuid.New()
		existingCar := &model.Car{
			ID:                carID,
			Name:              "Old Name",
			Category:          "Sedan",
			Availability:      true,
			StockAvailability: 5,
			RentalCosts:       100.00,
		}

		req := dto.UpdateCarRequest{
			Name: "New Name",
		}

		mockRepo.On("GetByID", mock.Anything, carID).Return(existingCar, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Car")).Return(nil)

		result, err := service.UpdateCar(context.Background(), carID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		mockRepo.AssertExpectations(t)
	})
}

func TestCarService_DeleteCar(t *testing.T) {
	mockRepo := new(MockCarRepository)
	service := NewCarService(mockRepo)

	t.Run("successful delete", func(t *testing.T) {
		carID := uuid.New()
		car := &model.Car{ID: carID, Name: "Test Car"}

		mockRepo.On("GetByID", mock.Anything, carID).Return(car, nil)
		mockRepo.On("Delete", mock.Anything, carID).Return(nil)

		err := service.DeleteCar(context.Background(), carID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete non-existent car", func(t *testing.T) {
		carID := uuid.New()

		mockRepo.On("GetByID", mock.Anything, carID).Return((*model.Car)(nil), ErrCarNotFound)

		err := service.DeleteCar(context.Background(), carID)

		assert.Error(t, err)
		assert.Equal(t, ErrCarNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
