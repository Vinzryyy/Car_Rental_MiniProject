package handler

import (
	"net/http"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/service"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	rentalService service.RentalService
}

func NewAdminHandler(rentalService service.RentalService) *AdminHandler {
	return &AdminHandler{
		rentalService: rentalService,
	}
}

// GetDashboardStats godoc
// @Summary Get admin dashboard statistics
// @Description Get overall statistics including revenue, total rentals, total users, and popular cars (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/admin/dashboard [get]
func (h *AdminHandler) GetDashboardStats(c echo.Context) error {
	stats, popularCars, err := h.rentalService.GetAdminDashboardStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to retrieve admin statistics",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "admin statistics retrieved successfully",
		Data: map[string]interface{}{
			"stats":        stats,
			"popular_cars": popularCars,
		},
	})
}
