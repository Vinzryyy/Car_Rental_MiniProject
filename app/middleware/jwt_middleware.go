package middleware

import (
	"net/http"
	"strings"

	"car_rental_miniproject/app/config"
	"car_rental_miniproject/service"

	"github.com/labstack/echo/v4"
)

type JWTMiddleware struct {
	authService service.AuthService
	cfg         *config.JWTConfig
}

func NewJWTMiddleware(authService service.AuthService, cfg *config.JWTConfig) *JWTMiddleware {
	return &JWTMiddleware{
		authService: authService,
		cfg:         cfg,
	}
}

func (m *JWTMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "missing authorization header",
			})
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid authorization format",
			})
		}

		tokenString := parts[1]
		userID, err := m.authService.ValidateToken(c.Request().Context(), tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid or expired token",
			})
		}

		// Set user ID in context for handlers to use
		c.Set("user_id", userID.String())
		c.Set("token", tokenString)

		return next(c)
	}
}

// OptionalAuth is similar to Authenticate but doesn't fail if token is missing/invalid
func (m *JWTMiddleware) OptionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return next(c)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			userID, err := m.authService.ValidateToken(c.Request().Context(), parts[1])
			if err == nil {
				c.Set("user_id", userID.String())
				c.Set("token", parts[1])
			}
		}

		return next(c)
	}
}
