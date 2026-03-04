package repository

import (
	"context"
	"time"

	"car_rental_miniproject/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarRepository interface {
	Create(ctx context.Context, car *model.Car) error
	GetAll(ctx context.Context, category string, availableOnly bool) ([]model.Car, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Car, error)
	Update(ctx context.Context, car *model.Car) error
	Delete(ctx context.Context, id uuid.UUID) error
	DecreaseStock(ctx context.Context, id uuid.UUID) error
	IncreaseStock(ctx context.Context, id uuid.UUID) error
}

type carRepository struct {
	pool *pgxpool.Pool
}

func NewCarRepository(pool *pgxpool.Pool) CarRepository {
	return &carRepository{pool: pool}
}

func (r *carRepository) Create(ctx context.Context, car *model.Car) error {
	query := `INSERT INTO cars (id, name, availability, stock_availability, rental_costs, category, description, image_url, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.pool.Exec(ctx, query, car.ID, car.Name, car.Availability, car.StockAvailability, car.RentalCosts, car.Category, car.Description, car.ImageURL, car.CreatedAt, car.UpdatedAt)
	return err
}

func (r *carRepository) GetAll(ctx context.Context, category string, availableOnly bool) ([]model.Car, error) {
	query := `SELECT id, name, availability, stock_availability, rental_costs, category, description, image_url, created_at, updated_at FROM cars WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if category != "" {
		query += ` AND category = $` + string(rune('0'+argIndex))
		args = append(args, category)
		argIndex++
	}

	if availableOnly {
		query += ` AND availability = $` + string(rune('0'+argIndex)) + ` AND stock_availability > 0`
		args = append(args, true)
		argIndex++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []model.Car
	for rows.Next() {
		var car model.Car
		err := rows.Scan(&car.ID, &car.Name, &car.Availability, &car.StockAvailability, &car.RentalCosts, &car.Category, &car.Description, &car.ImageURL, &car.CreatedAt, &car.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func (r *carRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Car, error) {
	query := `SELECT id, name, availability, stock_availability, rental_costs, category, description, image_url, created_at, updated_at FROM cars WHERE id = $1`
	car := &model.Car{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&car.ID, &car.Name, &car.Availability, &car.StockAvailability, &car.RentalCosts, &car.Category, &car.Description, &car.ImageURL, &car.CreatedAt, &car.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return car, nil
}

func (r *carRepository) Update(ctx context.Context, car *model.Car) error {
	query := `UPDATE cars SET name = $1, availability = $2, stock_availability = $3, rental_costs = $4, 
			  category = $5, description = $6, image_url = $7, updated_at = $8 WHERE id = $9`
	_, err := r.pool.Exec(ctx, query, car.Name, car.Availability, car.StockAvailability, car.RentalCosts, car.Category, car.Description, car.ImageURL, time.Now(), car.ID)
	return err
}

func (r *carRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM cars WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *carRepository) DecreaseStock(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE cars SET stock_availability = stock_availability - 1, 
			  availability = (stock_availability - 1 > 0), updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	return err
}

func (r *carRepository) IncreaseStock(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE cars SET stock_availability = stock_availability + 1, 
			  availability = true, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	return err
}
