package handler

import (
	"errors"
	"net/http"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/app/middleware"
	"car_rental_miniproject/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *middleware.CustomValidator
}

func NewAuthHandler(authService service.AuthService, validator *middleware.CustomValidator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 409 {object} dto.APIResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "validation failed",
			Error:   err.Error(),
		})
	}

	user, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return c.JSON(http.StatusConflict, dto.APIResponse{
				Success: false,
				Message: "user already exists",
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to register user",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "registration successful",
		Data: dto.UserResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			DepositAmount: user.DepositAmount,
		},
	})
}

// Login godoc
// @Summary Login user
// @Description Login with email and password to get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.APIResponse{data=dto.LoginResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "validation failed",
			Error:   err.Error(),
		})
	}

	token, response, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "invalid credentials",
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to login",
			Error:   err.Error(),
		})
	}

	_ = token // token is already in response

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "login successful",
		Data:    response,
	})
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh JWT token using a valid token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.APIResponse{data=dto.LoginResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "validation failed",
			Error:   err.Error(),
		})
	}

	token, response, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			return c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "invalid or expired token",
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to refresh token",
			Error:   err.Error(),
		})
	}

	_ = token // token is already in response

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "token refreshed successfully",
		Data:    response,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user (client should remove stored token)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	token := c.Get("token").(string)
	if err := h.authService.Logout(c.Request().Context(), token); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to logout",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "logout successful",
	})
}

// Me godoc
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 401 {object} dto.APIResponse
// @Router /api/auth/me [get]
func (h *AuthHandler) Me(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "invalid user ID",
			Error:   err.Error(),
		})
	}

	user, err := h.authService.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "user not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "user retrieved successfully",
		Data: dto.UserResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			DepositAmount: user.DepositAmount,
		},
	})
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Request a password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Email address"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /api/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var req dto.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "validation failed",
			Error:   err.Error(),
		})
	}

	token, err := h.authService.ForgotPassword(c.Request().Context(), req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to process request",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "password reset instructions sent",
		Data:    map[string]string{"token": token},
	})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using a valid reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Router /api/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req dto.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "validation failed",
			Error:   err.Error(),
		})
	}

	err := h.authService.ResetPassword(c.Request().Context(), req.Token, req.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrTokenExpired) {
			return c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "invalid or expired reset token",
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to reset password",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "password reset successful",
	})
}

// ChangePassword godoc
// @Summary Change password
// @Description Change password for authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "Old and new password"
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Router /api/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "validation failed",
			Error:   err.Error(),
		})
	}

	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "invalid user ID",
			Error:   err.Error(),
		})
	}

	err = h.authService.ChangePassword(c.Request().Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOldPassword) {
			return c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "invalid old password",
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to change password",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "password changed successfully",
	})
}

// UpdateProfile godoc
// @Summary Update profile
// @Description Update user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 409 {object} dto.APIResponse
// @Router /api/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "validation failed",
			Error:   err.Error(),
		})
	}

	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "invalid user ID",
			Error:   err.Error(),
		})
	}

	user, err := h.authService.UpdateProfile(c.Request().Context(), userID, req.Email)
	if err != nil {
		if err.Error() == "email already in use" {
			return c.JSON(http.StatusConflict, dto.APIResponse{
				Success: false,
				Message: "email already in use",
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to update profile",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "profile updated successfully",
		Data: dto.UserResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			DepositAmount: user.DepositAmount,
		},
	})
}
