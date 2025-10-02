package handler

import (
	"errors"
	"project-simple/internal/domain/dto"
	"project-simple/internal/service"
	"project-simple/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CarHandler struct {
	carService service.CarService
}

func NewCarHandler(carService service.CarService) *CarHandler {
	return &CarHandler{
		carService: carService,
	}
}

// CreateCar godoc
// @Summary Create a new car
// @Description Create a new car with the provided information
// @Tags cars
// @Accept json
// @Produce json
// @Param car body dto.CreateCarRequest true "Car information"
// @Success 201 {object} response.Response{data=dto.CarResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/cars [post]
func (h *CarHandler) CreateCar(c *gin.Context) {
	var req dto.CreateCarRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := h.formatValidationErrors(err)
		if validationErrors != nil {
			response.UnprocessableEntity(c, "Validation failed", validationErrors)
			return
		}
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	car, err := h.carService.CreateCar(&req)
	if err != nil {
		response.InternalServerError(c, "Failed to create car")
		return
	}

	response.Created(c, "Car created successfully", car)
}

// GetCarByID godoc
// @Summary Get a car by ID
// @Description Get detailed information about a specific car
// @Tags cars
// @Accept json
// @Produce json
// @Param id path string true "Car ID (UUID)"
// @Success 200 {object} response.Response{data=dto.CarResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/cars/{id} [get]
func (h *CarHandler) GetCarByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, "Invalid car ID format", nil)
		return
	}

	car, err := h.carService.GetCarByID(id)
	if err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			response.NotFound(c, "Car not found")
			return
		}
		response.InternalServerError(c, "Failed to retrieve car")
		return
	}

	response.Success(c, "Car retrieved successfully", car)
}

// GetAllCars godoc
// @Summary Get all cars with pagination
// @Description Get a paginated list of all cars with optional sorting
// @Tags cars
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param page_size query int false "Items per page (default: 10, max: 100)" minimum(1) maximum(100)
// @Param sort_by query string false "Sort by field (name, engine_version, created_at)" Enums(name, engine_version, created_at)
// @Param sort_dir query string false "Sort direction (asc, desc)" Enums(asc, desc)
// @Success 200 {object} response.Response{data=dto.PaginatedResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/cars [get]
func (h *CarHandler) GetAllCars(c *gin.Context) {
	var pagination dto.PaginationRequest

	if err := c.ShouldBindQuery(&pagination); err != nil {
		validationErrors := h.formatValidationErrors(err)
		if validationErrors != nil {
			response.UnprocessableEntity(c, "Validation failed", validationErrors)
			return
		}
		response.BadRequest(c, "Invalid query parameters", err.Error())
		return
	}

	result, err := h.carService.GetAllCars(&pagination)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve cars")
		return
	}

	response.Success(c, "Cars retrieved successfully", result)
}

// UpdateCar godoc
// @Summary Update a car
// @Description Update an existing car's information
// @Tags cars
// @Accept json
// @Produce json
// @Param id path string true "Car ID (UUID)"
// @Param car body dto.UpdateCarRequest true "Updated car information"
// @Success 200 {object} response.Response{data=dto.CarResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/cars/{id} [put]
func (h *CarHandler) UpdateCar(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, "Invalid car ID format", nil)
		return
	}

	var req dto.UpdateCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := h.formatValidationErrors(err)
		if validationErrors != nil {
			response.UnprocessableEntity(c, "Validation failed", validationErrors)
			return
		}
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	car, err := h.carService.UpdateCar(id, &req)
	if err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			response.NotFound(c, "Car not found")
			return
		}
		response.InternalServerError(c, "Failed to update car")
		return
	}

	response.Success(c, "Car updated successfully", car)
}

// DeleteCar godoc
// @Summary Delete a car
// @Description Delete a car by ID (soft delete)
// @Tags cars
// @Accept json
// @Produce json
// @Param id path string true "Car ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/cars/{id} [delete]
func (h *CarHandler) DeleteCar(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, "Invalid car ID format", nil)
		return
	}

	err = h.carService.DeleteCar(id)
	if err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			response.NotFound(c, "Car not found")
			return
		}
		response.InternalServerError(c, "Failed to delete car")
		return
	}

	response.NoContent(c)
}

func (h *CarHandler) formatValidationErrors(err error) []response.ValidationError {
	var validationErrors []response.ValidationError

	if validatorErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validatorErrs {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   e.Field(),
				Message: h.getErrorMessage(e),
			})
		}
		return validationErrors
	}

	return nil
}

func (h *CarHandler) getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value is too short or small"
	case "max":
		return "Value is too long or large"
	case "oneof":
		return "Invalid value. Allowed values: " + err.Param()
	default:
		return "Invalid value"
	}
}

