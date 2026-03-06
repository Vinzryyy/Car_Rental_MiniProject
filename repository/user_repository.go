package repository

import (
	"context"
	"fmt"
	"time"

	"car_rental_miniproject/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	WithTx(tx pgx.Tx) UserRepository
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateDeposit(ctx context.Context, id uuid.UUID, amount float64) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	Update(ctx context.Context, user *model.User) error
}

type userRepository struct {
	pool DBPool
	tx   pgx.Tx
}

func NewUserRepository(pool DBPool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) WithTx(tx pgx.Tx) UserRepository {
	return &userRepository{pool: r.pool, tx: tx}
}

func (r *userRepository) getQuerier() Querier {
	if r.tx != nil {
		return r.tx
	}
	return r.pool
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (id, email, password, deposit_amount, role, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.getQuerier().Exec(ctx, query, user.ID, user.Email, user.Password, user.DepositAmount, user.Role, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, password, deposit_amount, role, created_at, updated_at FROM users WHERE email = $1`
	user := &model.User{}
	err := r.getQuerier().QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.DepositAmount, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `SELECT id, email, password, deposit_amount, role, created_at, updated_at FROM users WHERE id = $1`
	user := &model.User{}
	err := r.getQuerier().QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password, &user.DepositAmount, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `SELECT id, email, password, deposit_amount, role, created_at, updated_at FROM users WHERE id = $1 FOR UPDATE`
	user := &model.User{}
	err := r.getQuerier().QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password, &user.DepositAmount, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateDeposit(ctx context.Context, id uuid.UUID, amount float64) error {
	query := `UPDATE users SET deposit_amount = deposit_amount + $1, updated_at = $2 
			  WHERE id = $3 AND (deposit_amount + $1) >= 0`
	result, err := r.getQuerier().Exec(ctx, query, amount, time.Now(), id)
	if err != nil {
		return err
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("insufficient balance or user not found")
	}
	
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	query := `UPDATE users SET password = $1, updated_at = $2 WHERE id = $3`
	_, err := r.getQuerier().Exec(ctx, query, password, time.Now(), id)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET email = $1, deposit_amount = $2, role = $3, updated_at = $4 WHERE id = $5`
	_, err := r.getQuerier().Exec(ctx, query, user.Email, user.DepositAmount, user.Role, user.UpdatedAt, user.ID)
	return err
}
