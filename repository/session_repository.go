package repository

import (
	"context"
	"time"

	"car_rental_miniproject/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.UserSession) error
	GetByToken(ctx context.Context, token string) (*model.UserSession, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type sessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) SessionRepository {
	return &sessionRepository{pool: pool}
}

func (r *sessionRepository) Create(ctx context.Context, session *model.UserSession) error {
	query := `INSERT INTO user_sessions (id, user_id, token, expires_at, created_at) 
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, query, session.ID, session.UserID, session.Token, session.ExpiresAt, session.CreatedAt)
	return err
}

func (r *sessionRepository) GetByToken(ctx context.Context, token string) (*model.UserSession, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM user_sessions WHERE token = $1`
	session := &model.UserSession{}
	err := r.pool.QueryRow(ctx, query, token).Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *sessionRepository) DeleteByToken(ctx context.Context, token string) error {
	query := `DELETE FROM user_sessions WHERE token = $1`
	_, err := r.pool.Exec(ctx, query, time.Now())
	return err
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM user_sessions WHERE expires_at < $1`
	_, err := r.pool.Exec(ctx, query, time.Now())
	return err
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}
