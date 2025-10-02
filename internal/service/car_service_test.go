package service

import (
	"errors"
	"testing"
	"time"

	"project-simple/internal/domain/dto"
	"project-simple/internal/domain/entity"
	"project-simple/internal/repository"
	"project-simple/internal/repository/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCarService_CreateCar(t *testing.T) {
	t.Run("Success - Create car with valid data", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		req := &dto.CreateCarRequest{
			Name:          "Honda Civic",
			EngineVersion: "2.0",
		}

		mockRepo.On("Create", mock.AnythingOfType("*entity.Car")).Return(nil).Run(func(args mock.Arguments) {
			car := args.Get(0).(*entity.Car)
			car.ID = uuid.New()
			car.CreatedAt = time.Now()
			car.UpdatedAt = time.Now()
		})

		result, err := service.CreateCar(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, req.EngineVersion, result.EngineVersion)
		assert.NotEqual(t, uuid.Nil, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository create fails", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		req := &dto.CreateCarRequest{
			Name:          "Honda Civic",
			EngineVersion: "2.0",
		}

		expectedError := errors.New("database error")
		mockRepo.On("Create", mock.AnythingOfType("*entity.Car")).Return(expectedError)

		result, err := service.CreateCar(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCarService_GetCarByID(t *testing.T) {
	t.Run("Success - Get existing car", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		carID := uuid.New()
		expectedCar := &entity.Car{
			ID:            carID,
			Name:          "Toyota Corolla",
			EngineVersion: "1.8",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		mockRepo.On("FindByID", carID).Return(expectedCar, nil)

		result, err := service.GetCarByID(carID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, carID, result.ID)
		assert.Equal(t, expectedCar.Name, result.Name)
		assert.Equal(t, expectedCar.EngineVersion, result.EngineVersion)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Car not found", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		carID := uuid.New()
		mockRepo.On("FindByID", carID).Return(nil, repository.ErrCarNotFound)

		result, err := service.GetCarByID(carID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrCarNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository error", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		carID := uuid.New()
		expectedError := errors.New("database connection error")
		mockRepo.On("FindByID", carID).Return(nil, expectedError)

		result, err := service.GetCarByID(carID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCarService_GetAllCars(t *testing.T) {
	t.Run("Success - Get paginated cars", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		pagination := &dto.PaginationRequest{
			Page:     1,
			PageSize: 10,
		}

		expectedCars := []entity.Car{
			{
				ID:            uuid.New(),
				Name:          "Honda Civic",
				EngineVersion: "2.0",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				ID:            uuid.New(),
				Name:          "Toyota Corolla",
				EngineVersion: "1.8",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}

		mockRepo.On("FindAll", mock.AnythingOfType("*dto.PaginationRequest")).Return(expectedCars, int64(25), nil)

		result, err := service.GetAllCars(pagination)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, 1, result.Pagination.CurrentPage)
		assert.Equal(t, 10, result.Pagination.PageSize)
		assert.Equal(t, int64(25), result.Pagination.TotalRecords)
		assert.Equal(t, 3, result.Pagination.TotalPages) // 25 records / 10 per page = 3 pages
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Empty result", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		pagination := &dto.PaginationRequest{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.On("FindAll", mock.AnythingOfType("*dto.PaginationRequest")).Return([]entity.Car{}, int64(0), nil)

		result, err := service.GetAllCars(pagination)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Data)
		assert.Equal(t, int64(0), result.Pagination.TotalRecords)
		assert.Equal(t, 0, result.Pagination.TotalPages)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository error", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		pagination := &dto.PaginationRequest{
			Page:     1,
			PageSize: 10,
		}

		expectedError := errors.New("database error")
		mockRepo.On("FindAll", mock.AnythingOfType("*dto.PaginationRequest")).Return([]entity.Car{}, int64(0), expectedError)

		result, err := service.GetAllCars(pagination)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCarService_UpdateCar(t *testing.T) {
	t.Run("Success - Update existing car", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		carID := uuid.New()
		existingCar := &entity.Car{
			ID:            carID,
			Name:          "Honda Civic",
			EngineVersion: "2.0",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		req := &dto.UpdateCarRequest{
			Name:          "Honda Civic Sport",
			EngineVersion: "2.0",
		}

		updatedCar := &entity.Car{
			ID:            carID,
			Name:          req.Name,
			EngineVersion: req.EngineVersion,
			CreatedAt:     existingCar.CreatedAt,
			UpdatedAt:     time.Now(),
		}

		mockRepo.On("FindByID", carID).Return(existingCar, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*entity.Car")).Return(nil)
		mockRepo.On("FindByID", carID).Return(updatedCar, nil).Once()

		result, err := service.UpdateCar(carID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, carID, result.ID)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, req.EngineVersion, result.EngineVersion)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Partial update (name only)", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		carID := uuid.New()
		existingCar := &entity.Car{
			ID:            carID,
			Name:          "Honda Civic",
			EngineVersion: "2.0",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		req := &dto.UpdateCarRequest{
			Name: "Honda Civic Sport",
		}

		updatedCar := &entity.Car{
			ID:            carID,
			Name:          req.Name,
			EngineVersion: existingCar.EngineVersion,
			CreatedAt:     existingCar.CreatedAt,
			UpdatedAt:     time.Now(),
		}

		mockRepo.On("FindByID", carID).Return(existingCar, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*entity.Car")).Return(nil)
		mockRepo.On("FindByID", carID).Return(updatedCar, nil).Once()

		result, err := service.UpdateCar(carID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, existingCar.EngineVersion, result.EngineVersion)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Car not found", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		carID := uuid.New()
		req := &dto.UpdateCarRequest{
			Name: "Updated Name",
		}

		mockRepo.On("FindByID", carID).Return(nil, repository.ErrCarNotFound)

		result, err := service.UpdateCar(carID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrCarNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Update fails", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		carID := uuid.New()
		existingCar := &entity.Car{
			ID:            carID,
			Name:          "Honda Civic",
			EngineVersion: "2.0",
		}

		req := &dto.UpdateCarRequest{
			Name: "Updated Name",
		}

		expectedError := errors.New("database error")
		mockRepo.On("FindByID", carID).Return(existingCar, nil)
		mockRepo.On("Update", mock.AnythingOfType("*entity.Car")).Return(expectedError)

		result, err := service.UpdateCar(carID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCarService_DeleteCar(t *testing.T) {
	t.Run("Success - Delete existing car", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		carID := uuid.New()
		mockRepo.On("Delete", carID).Return(nil)

		err := service.DeleteCar(carID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Car not found", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		carID := uuid.New()
		mockRepo.On("Delete", carID).Return(repository.ErrCarNotFound)

		err := service.DeleteCar(carID)

		assert.Error(t, err)
		assert.Equal(t, ErrCarNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository error", func(t *testing.T) {
		mockRepo := new(mocks.MockCarRepository)
		service := NewCarService(mockRepo)

		carID := uuid.New()
		expectedError := errors.New("database error")
		mockRepo.On("Delete", carID).Return(expectedError)

		err := service.DeleteCar(carID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
