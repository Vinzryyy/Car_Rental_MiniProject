package repository

import (
	"context"
	"fmt"
	"time"

	"car_rental_miniproject/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarFilter struct {
	Category      string
	AvailableOnly bool
	Search        string
	SortBy        string
	SortOrder     string
	Limit         int
	Offset        int
}

type CarRepository interface {
	Create(ctx context.Context, car *model.Car) error
	GetAll(ctx context.Context, filter CarFilter) ([]model.Car, int, error)
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

func (r *carRepository) GetAll(ctx context.Context, filter CarFilter) ([]model.Car, int, error) {
	// Base query for counting total results (without limit/offset)
	countQuery := `SELECT COUNT(*) FROM cars WHERE 1=1`
	
	// Base query for fetching data
	dataQuery := `SELECT id, name, availability, stock_availability, rental_costs, category, description, image_url, created_at, updated_at FROM cars WHERE 1=1`
	
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if filter.Category != "" {
		whereClause += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, filter.Category)
		argIndex++
	}

	if filter.AvailableOnly {
		whereClause += fmt.Sprintf(" AND availability = $%d AND stock_availability > 0", argIndex)
		args = append(args, true)
		argIndex++
	}

	if filter.Search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	countQuery += whereClause
	dataQuery += whereClause

	// Sorting
	sortBy := "created_at"
	if filter.SortBy != "" {
		// Basic validation to prevent SQL injection (only allow specific columns)
		allowedSortColumns := map[string]bool{
			"name":               true,
			"rental_costs":       true,
			"category":           true,
			"created_at":         true,
			"stock_availability": true,
		}
		if allowedSortColumns[filter.SortBy] {
			sortBy = filter.SortBy
		}
	}

	sortOrder := "DESC"
	if filter.SortOrder != "" && (filter.SortOrder == "ASC" || filter.SortOrder == "DESC") {
		sortOrder = filter.SortOrder
	}

	dataQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Pagination
	limit := 10
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	
	offset := 0
	if filter.Offset > 0 {
		offset = filter.Offset
	}

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args[:argIndex-1]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cars []model.Car
	for rows.Next() {
		var car model.Car
		err := rows.Scan(&car.ID, &car.Name, &car.Availability, &car.StockAvailability, &car.RentalCosts, &car.Category, &car.Description, &car.ImageURL, &car.CreatedAt, &car.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		cars = append(cars, car)
	}

	return cars, total, nil
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
			  availability = (stock_availability - 1 > 0), updated_at = $2 
			  WHERE id = $1 AND stock_availability > 0`
	result, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("car not available or out of stock")
	}
	
	return nil
}

func (r *carRepository) IncreaseStock(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE cars SET stock_availability = stock_availability + 1, 
			  availability = true, updated_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	return err
}
