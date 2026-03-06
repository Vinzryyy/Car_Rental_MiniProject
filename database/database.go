package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"car_rental_miniproject/app/config"
	"car_rental_miniproject/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxpoolAdapter struct {
	pool *pgxpool.Pool
}

func (a *pgxpoolAdapter) Begin(ctx context.Context) (pgx.Tx, error) {
	return a.pool.Begin(ctx)
}

func (a *pgxpoolAdapter) Close() {
	a.pool.Close()
}

func (a *pgxpoolAdapter) Ping(ctx context.Context) error {
	return a.pool.Ping(ctx)
}

// For accessing the raw pool if needed (e.g. for healthchecks or migrations)
func (a *pgxpoolAdapter) RawPool() *pgxpool.Pool {
	return a.pool
}

type Database struct {
	Pool repository.DBPool
}

var dbInstance *Database

func Initialize(cfg *config.DatabaseConfig) (*Database, error) {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = cfg.DSN()
	}

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	// Create a context with timeout for initialization
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	adapter := &pgxpoolAdapter{pool: pool}

	// Test connection with timeout
	if err := adapter.Ping(ctx); err != nil {
		log.Printf("Warning: database ping failed during init: %v", err)
	}

	dbInstance = &Database{Pool: adapter}
	log.Println("Database connection setup complete")

	return dbInstance, nil
}

func GetDB() *Database {
	return dbInstance
}

func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
		log.Println("Database connection closed")
	}
}
