package repository

import (
	"context"
	"time"

	"car_rental_miniproject/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TopUpRepository interface {
	Create(ctx context.Context, transaction *model.TopUpTransaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.TopUpTransaction, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.TopUpTransaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdatePaymentURL(ctx context.Context, id uuid.UUID, paymentURL string) error
	Update(ctx context.Context, transaction *model.TopUpTransaction) error
}

type topUpRepository struct {
	pool *pgxpool.Pool
}

func NewTopUpRepository(pool *pgxpool.Pool) TopUpRepository {
	return &topUpRepository{pool: pool}
}

func (r *topUpRepository) Create(ctx context.Context, transaction *model.TopUpTransaction) error {
	query := `INSERT INTO top_up_transactions (id, user_id, amount, status, payment_method, payment_url, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, transaction.ID, transaction.UserID, transaction.Amount, transaction.Status, transaction.PaymentMethod, transaction.PaymentURL, transaction.CreatedAt, transaction.UpdatedAt)
	return err
}

func (r *topUpRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.TopUpTransaction, error) {
	query := `SELECT id, user_id, amount, status, payment_method, payment_url, created_at, updated_at FROM top_up_transactions WHERE id = $1`
	transaction := &model.TopUpTransaction{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Status, &transaction.PaymentMethod, &transaction.PaymentURL, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *topUpRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.TopUpTransaction, error) {
	query := `SELECT id, user_id, amount, status, payment_method, payment_url, created_at, updated_at FROM top_up_transactions WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.TopUpTransaction
	for rows.Next() {
		var transaction model.TopUpTransaction
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Status, &transaction.PaymentMethod, &transaction.PaymentURL, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *topUpRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE top_up_transactions SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, status, time.Now(), id)
	return err
}

func (r *topUpRepository) UpdatePaymentURL(ctx context.Context, id uuid.UUID, paymentURL string) error {
	query := `UPDATE top_up_transactions SET payment_url = $1, updated_at = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, paymentURL, time.Now(), id)
	return err
}

func (r *topUpRepository) Update(ctx context.Context, transaction *model.TopUpTransaction) error {
	query := `UPDATE top_up_transactions SET amount = $1, status = $2, payment_method = $3, payment_url = $4, updated_at = $5 WHERE id = $6`
	_, err := r.pool.Exec(ctx, query, transaction.Amount, transaction.Status, transaction.PaymentMethod, transaction.PaymentURL, transaction.UpdatedAt, transaction.ID)
	return err
}
