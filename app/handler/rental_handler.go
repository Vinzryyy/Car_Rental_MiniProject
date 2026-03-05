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

type RentalHandler struct {
	rentalService service.RentalService
	topUpService  service.TopUpService
	validator     *middleware.CustomValidator
}

func NewRentalHandler(rentalService service.RentalService, topUpService service.TopUpService, validator *middleware.CustomValidator) *RentalHandler {
	return &RentalHandler{
		rentalService: rentalService,
		topUpService:  topUpService,
		validator:     validator,
	}
}

// RentCar godoc
// @Summary Rent a car
// @Description Rent a car using user's deposit
// @Tags rentals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RentCarRequest true "Rental details"
// @Success 201 {object} dto.APIResponse{data=model.RentalHistory}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 422 {object} dto.APIResponse
// @Router /api/rentals [post]
func (h *RentalHandler) RentCar(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "invalid user ID",
			Error:   err.Error(),
		})
	}

	var req dto.RentCarRequest
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
			Errors:  middleware.FormatValidationErrors(err),
		})
	}

	rental, err := h.rentalService.RentCar(c.Request().Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCarNotFound):
			return c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Message: "car not found",
				Error:   err.Error(),
			})
		case errors.Is(err, service.ErrCarNotAvailable):
			return c.JSON(http.StatusUnprocessableEntity, dto.APIResponse{
				Success: false,
				Message: "car not available",
				Error:   err.Error(),
			})
		case errors.Is(err, service.ErrInsufficientDeposit):
			return c.JSON(http.StatusUnprocessableEntity, dto.APIResponse{
				Success: false,
				Message: "insufficient deposit. Please top up your balance",
				Error:   err.Error(),
			})
		case errors.Is(err, service.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Message: "user not found",
				Error:   err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.APIResponse{
				Success: false,
				Message: "failed to rent car",
				Error:   err.Error(),
			})
		}
	}

	return c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "car rented successfully",
		Data:    rental,
	})
}

// GetMyRentals godoc
// @Summary Get my rentals
// @Description Get rental history for authenticated user
// @Tags rentals
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=[]dto.RentalHistoryResponse}
// @Failure 401 {object} dto.APIResponse
// @Router /api/rentals/my [get]
func (h *RentalHandler) GetMyRentals(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "invalid user ID",
			Error:   err.Error(),
		})
	}

	rentals, err := h.rentalService.GetRentalsByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to retrieve rentals",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "rentals retrieved successfully",
		Data:    rentals,
	})
}

// GetBookingReport godoc
// @Summary Get booking report
// @Description Get comprehensive booking report for authenticated user
// @Tags rentals
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=dto.BookingReportResponse}
// @Failure 401 {object} dto.APIResponse
// @Router /api/rentals/booking-report [get]
func (h *RentalHandler) GetBookingReport(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "invalid user ID",
			Error:   err.Error(),
		})
	}

	report, err := h.rentalService.GetBookingReport(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to retrieve booking report",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "booking report retrieved successfully",
		Data:    report,
	})
}

// TopUp godoc
// @Summary Top up deposit
// @Description Add funds to user deposit balance
// @Tags topup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TopUpRequest true "Top-up amount"
// @Success 201 {object} dto.APIResponse{data=model.TopUpTransaction}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Router /api/topup [post]
func (h *RentalHandler) TopUp(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "invalid user ID",
			Error:   err.Error(),
		})
	}

	var req dto.TopUpRequest
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
			Errors:  middleware.FormatValidationErrors(err),
		})
	}

	transaction, err := h.topUpService.CreateTopUp(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to create top-up transaction",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "top-up transaction created. Please complete payment.",
		Data:    transaction,
	})
}

// GetTopUpHistory godoc
// @Summary Get top-up history
// @Description Get top-up transaction history for authenticated user
// @Tags topup
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=[]model.TopUpTransaction}
// @Failure 401 {object} dto.APIResponse
// @Router /api/topup/history [get]
func (h *RentalHandler) GetTopUpHistory(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "invalid user ID",
			Error:   err.Error(),
		})
	}

	transactions, err := h.topUpService.GetTopUpsByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to retrieve top-up history",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "top-up history retrieved successfully",
		Data:    transactions,
	})
}

// ConfirmPayment godoc
// @Summary Confirm payment
// @Description Confirm payment for a rental (callback from payment gateway)
// @Tags rentals
// @Produce json
// @Security BearerAuth
// @Param id path string true "Rental ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /api/rentals/:id/confirm-payment [post]
func (h *RentalHandler) ConfirmPayment(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "invalid user ID",
			Error:   err.Error(),
		})
	}

	rentalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid rental ID",
			Error:   err.Error(),
		})
	}

	if err := h.rentalService.ConfirmPayment(c.Request().Context(), rentalID, userID); err != nil {
		switch {
		case errors.Is(err, service.ErrRentalNotFound):
			return c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Message: "rental not found",
				Error:   err.Error(),
			})
		case errors.Is(err, service.ErrUnauthorizedAccess):
			return c.JSON(http.StatusForbidden, dto.APIResponse{
				Success: false,
				Message: "you do not have permission to confirm this payment",
				Error:   err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.APIResponse{
				Success: false,
				Message: "failed to confirm payment",
				Error:   err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "payment confirmed successfully",
	})
}

// CancelRental godoc
// @Summary Cancel rental
// @Description Cancel a pending rental
// @Tags rentals
// @Produce json
// @Security BearerAuth
// @Param id path string true "Rental ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /api/rentals/:id/cancel [post]
func (h *RentalHandler) CancelRental(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "invalid user ID",
			Error:   err.Error(),
		})
	}

	rentalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid rental ID",
			Error:   err.Error(),
		})
	}

	if err := h.rentalService.CancelRental(c.Request().Context(), rentalID, userID); err != nil {
		switch {
		case errors.Is(err, service.ErrRentalNotFound):
			return c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Message: "rental not found",
				Error:   err.Error(),
			})
		case errors.Is(err, service.ErrUnauthorizedAccess):
			return c.JSON(http.StatusForbidden, dto.APIResponse{
				Success: false,
				Message: "you do not have permission to cancel this rental",
				Error:   err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.APIResponse{
				Success: false,
				Message: "failed to cancel rental",
				Error:   err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "rental cancelled successfully",
	})
}
