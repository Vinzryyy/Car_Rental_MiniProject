package database

import (
	"context"
	"fmt"
	"log"

	"car_rental_miniproject/app/config"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Database struct {
	Pool *pgxpool.Pool
}

var db *Database

func Initialize(cfg *config.DatabaseConfig) (*Database, error) {
	connString := cfg.DSN()

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	db = &Database{Pool: pool}
	log.Println("Database connection established successfully")

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
