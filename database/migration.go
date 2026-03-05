package database

import (
	"fmt"
)

// RunMigrations creates the necessary tables by reading from ddl.sql
// Note: This is deprecated - manage schema via Supabase Dashboard
func (d *Database) RunMigrations() error {
	return fmt.Errorf("migrations are disabled - manage schema via Supabase Dashboard")
}

// SeedData inserts initial car data by reading from dml.sql
// Note: This is deprecated - manage data via Supabase Dashboard
func (d *Database) SeedData() error {
	return fmt.Errorf("seeding is disabled - manage data via Supabase Dashboard")
}