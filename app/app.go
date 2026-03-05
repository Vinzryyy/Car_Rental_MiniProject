package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	echoSwagger "github.com/swaggo/echo-swagger"

	"car_rental_miniproject/app/config"
	"car_rental_miniproject/app/handler"
	"car_rental_miniproject/app/middleware"
	"car_rental_miniproject/database"
	"car_rental_miniproject/repository"
	"car_rental_miniproject/service"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type App struct {
	echo     *echo.Echo
	config   *config.Config
	database *database.Database
}

func NewApp(cfg *config.Config) (*App, error) {
	// Initialize database
	db, err := database.Initialize(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Note: Migrations and seeding are skipped - manage schema via Supabase Dashboard
	// Run migrations
	// if err := db.RunMigrations(); err != nil {
	// 	return nil, fmt.Errorf("failed to run migrations: %w", err)
	// }

	// Seed initial data
	// if err := db.SeedData(); err != nil {
	// 	log.Printf("Warning: failed to seed data: %v", err)
	// }

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Use(echoMiddleware.RequestID())

	// Set custom validator
	e.Validator = middleware.NewCustomValidator()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.Pool)
	carRepo := repository.NewCarRepository(db.Pool)
	rentalRepo := repository.NewRentalRepository(db.Pool)
	topUpRepo := repository.NewTopUpRepository(db.Pool)
	sessionRepo := repository.NewSessionRepository(db.Pool)

	// Initialize services
	emailService := service.NewEmailService(cfg)
	authService := service.NewAuthService(userRepo, sessionRepo, &cfg.JWT, emailService)
	carService := service.NewCarService(carRepo)
	paymentService := service.NewXenditPaymentService(cfg)
	rentalService := service.NewRentalService(rentalRepo, carRepo, userRepo, paymentService, emailService)
	topUpService := service.NewTopUpService(topUpRepo, userRepo, paymentService, emailService)

	// Initialize middleware
	jwtMiddleware := middleware.NewJWTMiddleware(authService, &cfg.JWT)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, middleware.NewCustomValidator())
	carHandler := handler.NewCarHandler(carService, middleware.NewCustomValidator())
	rentalHandler := handler.NewRentalHandler(rentalService, topUpService, middleware.NewCustomValidator())
	webhookHandler := handler.NewPaymentWebhookHandler(rentalService, topUpService, paymentService, emailService)

	// Setup routes
	setupRoutes(e, authHandler, carHandler, rentalHandler, jwtMiddleware, webhookHandler)

	return &App{
		echo:     e,
		config:   cfg,
		database: db,
	}, nil
}

func setupRoutes(e *echo.Echo, authHandler *handler.AuthHandler, carHandler *handler.CarHandler, rentalHandler *handler.RentalHandler, jwtMiddleware *middleware.JWTMiddleware, webhookHandler *handler.PaymentWebhookHandler) {
	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// API routes
	api := e.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/forgot-password", authHandler.ForgotPassword)
	auth.POST("/reset-password", authHandler.ResetPassword)

	// Auth routes (protected)
	authProtected := api.Group("/auth", jwtMiddleware.Authenticate)
	authProtected.GET("/me", authHandler.Me)
	authProtected.POST("/logout", authHandler.Logout)
	authProtected.POST("/refresh", authHandler.RefreshToken)
	authProtected.POST("/change-password", authHandler.ChangePassword)
	authProtected.PUT("/profile", authHandler.UpdateProfile)

	// Car routes
	cars := api.Group("/cars")
	cars.GET("", carHandler.GetAllCars)
	cars.GET("/:id", carHandler.GetCarByID)
	cars.POST("", carHandler.CreateCar, jwtMiddleware.Authenticate)
	cars.PUT("/:id", carHandler.UpdateCar, jwtMiddleware.Authenticate)
	cars.DELETE("/:id", carHandler.DeleteCar, jwtMiddleware.Authenticate)

	// Rental routes (protected)
	rentals := api.Group("/rentals", jwtMiddleware.Authenticate)
	rentals.POST("", rentalHandler.RentCar)
	rentals.GET("/my", rentalHandler.GetMyRentals)
	rentals.GET("/booking-report", rentalHandler.GetBookingReport)
	rentals.POST("/:id/confirm-payment", rentalHandler.ConfirmPayment)
	rentals.POST("/:id/cancel", rentalHandler.CancelRental)

	// Top-up routes (protected)
	topup := api.Group("/topup", jwtMiddleware.Authenticate)
	topup.POST("", rentalHandler.TopUp)
	topup.GET("/history", rentalHandler.GetTopUpHistory)

	// Webhook routes (public, but should be secured by Xendit callback verification)
	webhook := api.Group("/webhook")
	webhook.POST("/payment", webhookHandler.PaymentNotification)
}

func (a *App) Start() {
	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", a.config.Server.Port)
		log.Printf("Starting server on %s", addr)
		log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", a.config.Server.Port)
		if err := a.echo.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.echo.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	a.database.Close()

	log.Println("Server exited gracefully")
}

func (a *App) Echo() *echo.Echo {
	return a.echo
}
