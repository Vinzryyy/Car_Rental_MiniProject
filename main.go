package main

import (
	"log"

	"github.com/joho/godotenv"
	"car_rental_miniproject/app"
	"car_rental_miniproject/app/config"

	_ "car_rental_miniproject/docs"
)

// @title Rental Car API
// @version 1.0
// @description A REST API for rental car service built with Go, Echo, PostgreSQL, and JWT authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@rentalcar.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	// Load environment variables from .env file (override system env vars)
	if err := godotenv.Overload(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Create and start application
	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Start the server
	app.Start()
}
