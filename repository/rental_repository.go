package repository

import (
	"context"
	"time"

	"car_rental_miniproject/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RentalRepository interface {
	Create(ctx context.Context, rental *model.RentalHistory) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.RentalHistory, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.RentalHistory, error)
	GetByUserIDWithCarDetails(ctx context.Context, userID uuid.UUID) ([]model.RentalWithCarDetails, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status, paymentStatus string) error
	UpdatePaymentURL(ctx context.Context, id uuid.UUID, paymentURL string) error
	GetBookingReport(ctx context.Context, userID uuid.UUID) (*model.BookingReport, error)
	GetOverdueRentals(ctx context.Context) ([]model.RentalWithCarDetails, error)
	GetAdminStats(ctx context.Context) (*model.AdminStats, error)
	GetPopularCars(ctx context.Context, limit int) ([]model.PopularCar, error)
}

type rentalRepository struct {
	pool *pgxpool.Pool
}

func NewRentalRepository(pool *pgxpool.Pool) RentalRepository {
	return &rentalRepository{pool: pool}
}

func (r *rentalRepository) Create(ctx context.Context, rental *model.RentalHistory) error {
	query := `INSERT INTO rental_histories (id, user_id, car_id, rental_date, return_date, total_cost, status, payment_status, payment_url, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(ctx, query, rental.ID, rental.UserID, rental.CarID, rental.RentalDate, rental.ReturnDate, rental.TotalCost, rental.Status, rental.PaymentStatus, rental.PaymentURL, rental.CreatedAt, rental.UpdatedAt)
	return err
}

func (r *rentalRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.RentalHistory, error) {
	query := `SELECT id, user_id, car_id, rental_date, return_date, total_cost, status, payment_status, payment_url, created_at, updated_at FROM rental_histories WHERE id = $1`
	rental := &model.RentalHistory{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&rental.ID, &rental.UserID, &rental.CarID, &rental.RentalDate, &rental.ReturnDate, &rental.TotalCost, &rental.Status, &rental.PaymentStatus, &rental.PaymentURL, &rental.CreatedAt, &rental.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return rental, nil
}

func (r *rentalRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.RentalHistory, error) {
	query := `SELECT id, user_id, car_id, rental_date, return_date, total_cost, status, payment_status, payment_url, created_at, updated_at FROM rental_histories WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []model.RentalHistory
	for rows.Next() {
		var rental model.RentalHistory
		err := rows.Scan(&rental.ID, &rental.UserID, &rental.CarID, &rental.RentalDate, &rental.ReturnDate, &rental.TotalCost, &rental.Status, &rental.PaymentStatus, &rental.PaymentURL, &rental.CreatedAt, &rental.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rentals = append(rentals, rental)
	}

	return rentals, nil
}

func (r *rentalRepository) GetByUserIDWithCarDetails(ctx context.Context, userID uuid.UUID) ([]model.RentalWithCarDetails, error) {
	query := `SELECT rh.id, rh.user_id, rh.car_id, rh.rental_date, rh.return_date, rh.total_cost, rh.status, rh.payment_status, rh.payment_url, rh.created_at, rh.updated_at, c.name as car_name 
			  FROM rental_histories rh 
			  JOIN cars c ON rh.car_id = c.id 
			  WHERE rh.user_id = $1 ORDER BY rh.created_at DESC`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []model.RentalWithCarDetails
	for rows.Next() {
		var rental model.RentalWithCarDetails
		err := rows.Scan(&rental.RentalHistory.ID, &rental.RentalHistory.UserID, &rental.RentalHistory.CarID, &rental.RentalHistory.RentalDate, &rental.RentalHistory.ReturnDate, &rental.RentalHistory.TotalCost, &rental.RentalHistory.Status, &rental.RentalHistory.PaymentStatus, &rental.RentalHistory.PaymentURL, &rental.RentalHistory.CreatedAt, &rental.RentalHistory.UpdatedAt, &rental.CarName)
		if err != nil {
			return nil, err
		}
		rentals = append(rentals, rental)
	}

	return rentals, nil
}

func (r *rentalRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status, paymentStatus string) error {
	query := `UPDATE rental_histories SET status = $1, payment_status = $2, updated_at = $3 WHERE id = $4`
	_, err := r.pool.Exec(ctx, query, status, paymentStatus, time.Now(), id)
	return err
}

func (r *rentalRepository) UpdatePaymentURL(ctx context.Context, id uuid.UUID, paymentURL string) error {
	query := `UPDATE rental_histories SET payment_url = $1, updated_at = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, paymentURL, time.Now(), id)
	return err
}

func (r *rentalRepository) GetBookingReport(ctx context.Context, userID uuid.UUID) (*model.BookingReport, error) {
	query := `SELECT 
			  COUNT(rh.id) as total_rentals,
			  COUNT(CASE WHEN rh.status = 'active' THEN 1 END) as active_rentals,
			  COALESCE(SUM(rh.total_cost), 0) as total_spent
			  FROM rental_histories rh
			  WHERE rh.user_id = $1`
	
	report := &model.BookingReport{UserID: userID}
	err := r.pool.QueryRow(ctx, query, userID).Scan(&report.TotalRentals, &report.ActiveRentals, &report.TotalSpent)
	if err != nil {
		return nil, err
	}

	// Get user email and deposit
	userQuery := `SELECT email, deposit_amount FROM users WHERE id = $1`
	err = r.pool.QueryRow(ctx, userQuery, userID).Scan(&report.Email, &report.CurrentDeposit)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r *rentalRepository) GetOverdueRentals(ctx context.Context) ([]model.RentalWithCarDetails, error) {
	query := `SELECT rh.id, rh.user_id, rh.car_id, rh.rental_date, rh.return_date, rh.total_cost, rh.status, rh.payment_status, rh.payment_url, rh.created_at, rh.updated_at, c.name as car_name 
			  FROM rental_histories rh 
			  JOIN cars c ON rh.car_id = c.id 
			  WHERE rh.status = 'active' AND rh.return_date < $1`
	
	rows, err := r.pool.Query(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []model.RentalWithCarDetails
	for rows.Next() {
		var rental model.RentalWithCarDetails
		err := rows.Scan(&rental.RentalHistory.ID, &rental.RentalHistory.UserID, &rental.RentalHistory.CarID, &rental.RentalHistory.RentalDate, &rental.RentalHistory.ReturnDate, &rental.RentalHistory.TotalCost, &rental.RentalHistory.Status, &rental.RentalHistory.PaymentStatus, &rental.RentalHistory.PaymentURL, &rental.RentalHistory.CreatedAt, &rental.RentalHistory.UpdatedAt, &rental.CarName)
		if err != nil {
			return nil, err
		}
		rentals = append(rentals, rental)
	}

	return rentals, nil
}

func (r *rentalRepository) GetAdminStats(ctx context.Context) (*model.AdminStats, error) {
	stats := &model.AdminStats{}
	
	// Total Revenue (only from paid rentals)
	queryRevenue := `SELECT COALESCE(SUM(total_cost), 0) FROM rental_histories WHERE payment_status = 'paid'`
	err := r.pool.QueryRow(ctx, queryRevenue).Scan(&stats.TotalRevenue)
	if err != nil {
		return nil, err
	}

	// Total Rentals
	queryRentals := `SELECT COUNT(*) FROM rental_histories`
	err = r.pool.QueryRow(ctx, queryRentals).Scan(&stats.TotalRentals)
	if err != nil {
		return nil, err
	}

	// Total Users
	queryUsers := `SELECT COUNT(*) FROM users`
	err = r.pool.QueryRow(ctx, queryUsers).Scan(&stats.TotalUsers)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *rentalRepository) GetPopularCars(ctx context.Context, limit int) ([]model.PopularCar, error) {
	query := `SELECT c.id, c.name, COUNT(rh.id) as rental_count 
			  FROM cars c 
			  LEFT JOIN rental_histories rh ON c.id = rh.car_id 
			  GROUP BY c.id, c.name 
			  ORDER BY rental_count DESC 
			  LIMIT $1`
	
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var popularCars []model.PopularCar
	for rows.Next() {
		var pc model.PopularCar
		err := rows.Scan(&pc.CarID, &pc.CarName, &pc.RentalCount)
		if err != nil {
			return nil, err
		}
		popularCars = append(popularCars, pc)
	}

	return popularCars, nil
}
