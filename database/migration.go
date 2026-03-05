package database

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)


var sqlFiles embed.FS

// RunMigrations creates the necessary tables by reading from ddl.sql
func (d *Database) RunMigrations() error {
	content, err := sqlFiles.ReadFile("ddl.sql")
	if err != nil {
		return fmt.Errorf("failed to read ddl.sql: %w", err)
	}

	// Split by semicolon to execute individual statements
	statements := strings.Split(string(content), ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := d.Pool.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// SeedData inserts initial car data by reading from dml.sql
func (d *Database) SeedData() error {
	content, err := os.ReadFile(filepath.Join("database", "dml.sql"))
	if err != nil {
		return fmt.Errorf("failed to read dml.sql: %w", err)
	}

	// Split by semicolon to execute individual statements
	statements := strings.Split(string(content), ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := d.Pool.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("seed failed: %w", err)
		}
	}

	log.Println("Database seeded successfully")
	return nil
}

