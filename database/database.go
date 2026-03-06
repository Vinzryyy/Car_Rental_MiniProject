package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"car_rental_miniproject/app/config"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Database struct {
	Pool *pgxpool.Pool
}

var db *Database

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

	// Test connection with timeout
	if err := pool.Ping(ctx); err != nil {
		log.Printf("Warning: database ping failed during init: %v", err)
		// We still return the pool because it might recover later, 
		// or at least allow the healthcheck to respond.
	}

	db = &Database{Pool: pool}
	log.Println("Database connection setup complete")

	return db, nil
}

func GetDB() *Database {
	return db
}

func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
		log.Println("Database connection closed")
	}
}
