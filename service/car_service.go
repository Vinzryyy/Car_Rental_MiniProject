package service

import (
	"context"
	"errors"
	"time"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/model"
	"car_rental_miniproject/repository"

	"github.com/google/uuid"
)

var (
	ErrCarNotFound      = errors.New("car not found")
	ErrCarNotAvailable  = errors.New("car not available")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type CarService interface {
	CreateCar(ctx context.Context, req dto.CreateCarRequest) (*model.Car, error)
	GetAllCars(ctx context.Context, category string, availableOnly bool) ([]model.Car, error)
	GetCarByID(ctx context.Context, id uuid.UUID) (*model.Car, error)
	UpdateCar(ctx context.Context, id uuid.UUID, req dto.UpdateCarRequest) (*model.Car, error)
	DeleteCar(ctx context.Context, id uuid.UUID) error
}

type carService struct {
	carRepo repository.CarRepository
}

func NewCarService(carRepo repository.CarRepository) CarService {
	return &carService{
		carRepo: carRepo,
	}
}

func (s *carService) CreateCar(ctx context.Context, req dto.CreateCarRequest) (*model.Car, error) {
	car := &model.Car{
		ID:                uuid.New(),
		Name:              req.Name,
		Availability:      req.StockAvailability > 0,
		StockAvailability: req.StockAvailability,
		RentalCosts:       req.RentalCosts,
		Category:          req.Category,
		Description:       req.Description,
		ImageURL:          req.ImageURL,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.carRepo.Create(ctx, car); err != nil {
		return nil, err
	}

	return car, nil
}

func (s *carService) GetAllCars(ctx context.Context, category string, availableOnly bool) ([]model.Car, error) {
	return s.carRepo.GetAll(ctx, category, availableOnly)
}

func (s *carService) GetCarByID(ctx context.Context, id uuid.UUID) (*model.Car, error) {
	car, err := s.carRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCarNotFound
	}
	return car, nil
}

func (s *carService) UpdateCar(ctx context.Context, id uuid.UUID, req dto.UpdateCarRequest) (*model.Car, error) {
	car, err := s.carRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCarNotFound
	}

	// Update fields if provided
	if req.Name != "" {
		car.Name = req.Name
	}
	if req.StockAvailability >= 0 {
		car.StockAvailability = req.StockAvailability
		car.Availability = req.StockAvailability > 0
	}
	if req.RentalCosts > 0 {
		car.RentalCosts = req.RentalCosts
	}
	if req.Category != "" {
		car.Category = req.Category
	}
	if req.Description != "" {
		car.Description = req.Description
	}
	if req.ImageURL != "" {
		car.ImageURL = req.ImageURL
	}

	car.UpdatedAt = time.Now()

	if err := s.carRepo.Update(ctx, car); err != nil {
		return nil, err
	}

	return car, nil
}

func (s *carService) DeleteCar(ctx context.Context, id uuid.UUID) error {
	_, err := s.carRepo.GetByID(ctx, id)
	if err != nil {
		return ErrCarNotFound
	}

	return s.carRepo.Delete(ctx, id)
}
