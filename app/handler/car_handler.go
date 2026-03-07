package handler

import (
	"errors"
	"net/http"
	"strconv"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/app/middleware"
	"car_rental_miniproject/repository"
	"car_rental_miniproject/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CarHandler struct {
	carService   service.CarService
	imageService service.ImageService
	validator    *middleware.CustomValidator
}

func NewCarHandler(carService service.CarService, imageService service.ImageService, validator *middleware.CustomValidator) *CarHandler {
	return &CarHandler{
		carService:   carService,
		imageService: imageService,
		validator:    validator,
	}
}

// UploadCarImage godoc
// @Summary Upload car image
// @Description Upload an image for a car to Cloudinary (admin only)
// @Tags cars
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param image formData file true "Car image file"
// @Success 200 {object} dto.APIResponse{data=dto.PaymentResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/cars/upload [post]
func (h *CarHandler) UploadCarImage(c echo.Context) error {
	if h.imageService == nil {
		return c.JSON(http.StatusServiceUnavailable, dto.APIResponse{
			Success: false,
			Message: "image service not configured",
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "failed to get image file",
			Error:   err.Error(),
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to open image file",
			Error:   err.Error(),
		})
	}
	defer src.Close()

	url, err := h.imageService.UploadImage(c.Request().Context(), src, "cars")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to upload image",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "image uploaded successfully",
		Data: map[string]string{
			"url": url,
		},
	})
}

// GetAllCars godoc
// @Summary Get all cars
// @Description Get list of all available cars with optional filtering, searching, sorting and pagination
// @Tags cars
// @Produce json
// @Param category query string false "Filter by category"
// @Param available query bool false "Filter by availability"
// @Param search query string false "Search by name or description"
// @Param sort_by query string false "Sort by field (name, rental_costs, category, created_at, stock_availability)"
// @Param sort_order query string false "Sort order (ASC, DESC)"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} dto.APIResponse{data=dto.PaginatedCarsResponse}
// @Router /api/cars [get]
func (h *CarHandler) GetAllCars(c echo.Context) error {
	category := c.QueryParam("category")
	available := c.QueryParam("available")
	search := c.QueryParam("search")
	sortBy := c.QueryParam("sort_by")
	sortOrder := c.QueryParam("sort_order")
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	availableOnly := available == "true"

	filter := repository.CarFilter{
		Category:      category,
		AvailableOnly: availableOnly,
		Search:        search,
		SortBy:        sortBy,
		SortOrder:     sortOrder,
		Limit:         limit,
		Offset:        offset,
	}

	cars, total, err := h.carService.GetAllCars(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to retrieve cars",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "cars retrieved successfully",
		Data: dto.PaginatedCarsResponse{
			Cars:  cars,
			Total: total,
		},
	})
}

// GetCarByID godoc
// @Summary Get car by ID
// @Description Get details of a specific car by ID
// @Tags cars
// @Produce json
// @Param id path string true "Car ID"
// @Success 200 {object} dto.APIResponse{data=model.Car}
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /api/cars/:id [get]
func (h *CarHandler) GetCarByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid car ID",
			Error:   err.Error(),
		})
	}

	car, err := h.carService.GetCarByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			return c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Message: "car not found",
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to retrieve car",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "car retrieved successfully",
		Data:    car,
	})
}

// CreateCar godoc
// @Summary Create a new car
// @Description Create a new car (admin only)
// @Tags cars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCarRequest true "Car details"
// @Success 201 {object} dto.APIResponse{data=model.Car}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Router /api/cars [post]
func (h *CarHandler) CreateCar(c echo.Context) error {
	var req dto.CreateCarRequest
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

	car, err := h.carService.CreateCar(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to create car",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "car created successfully",
		Data:    car,
	})
}

// UpdateCar godoc
// @Summary Update a car
// @Description Update an existing car (admin only)
// @Tags cars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Car ID"
// @Param request body dto.UpdateCarRequest true "Car details"
// @Success 200 {object} dto.APIResponse{data=model.Car}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /api/cars/:id [put]
func (h *CarHandler) UpdateCar(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid car ID",
			Error:   err.Error(),
		})
	}

	var req dto.UpdateCarRequest
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

	car, err := h.carService.UpdateCar(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			return c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Message: "car not found",
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to update car",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "car updated successfully",
		Data:    car,
	})
}

// DeleteCar godoc
// @Summary Delete a car
// @Description Delete a car by ID (admin only)
// @Tags cars
// @Produce json
// @Security BearerAuth
// @Param id path string true "Car ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /api/cars/:id [delete]
func (h *CarHandler) DeleteCar(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid car ID",
			Error:   err.Error(),
		})
	}

	if err := h.carService.DeleteCar(c.Request().Context(), id); err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			return c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Message: "car not found",
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to delete car",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "car deleted successfully",
	})
}
