package service

import (
	"errors"
	"math"
	"project-simple/internal/domain/dto"
	"project-simple/internal/domain/entity"
	"project-simple/internal/repository"

	"github.com/google/uuid"
)

type CarService interface {
	CreateCar(req *dto.CreateCarRequest) (*dto.CarResponse, error)
	GetCarByID(id uuid.UUID) (*dto.CarResponse, error)
	GetAllCars(pagination *dto.PaginationRequest) (*dto.PaginatedResponse, error)
	UpdateCar(id uuid.UUID, req *dto.UpdateCarRequest) (*dto.CarResponse, error)
	DeleteCar(id uuid.UUID) error
}

type carService struct {
	carRepo repository.CarRepository
}

func NewCarService(carRepo repository.CarRepository) CarService {
	return &carService{
		carRepo: carRepo,
	}
}

func (s *carService) CreateCar(req *dto.CreateCarRequest) (*dto.CarResponse, error) {
	car := &entity.Car{
		Name:          req.Name,
		EngineVersion: req.EngineVersion,
	}

	if err := s.carRepo.Create(car); err != nil {
		return nil, err
	}

	return s.entityToResponse(car), nil
}

func (s *carService) GetCarByID(id uuid.UUID) (*dto.CarResponse, error) {
	car, err := s.carRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrCarNotFound) {
			return nil, ErrCarNotFound
		}
		return nil, err
	}

	return s.entityToResponse(car), nil
}

func (s *carService) GetAllCars(pagination *dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	pagination.SetDefaults()

	cars, total, err := s.carRepo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	carResponses := make([]dto.CarResponse, len(cars))
	for i, car := range cars {
		carResponses[i] = *s.entityToResponse(&car)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &dto.PaginatedResponse{
		Data: carResponses,
		Pagination: dto.PaginationMeta{
			CurrentPage:  pagination.Page,
			PageSize:     pagination.PageSize,
			TotalPages:   totalPages,
			TotalRecords: total,
		},
	}, nil
}

func (s *carService) UpdateCar(id uuid.UUID, req *dto.UpdateCarRequest) (*dto.CarResponse, error) {
	// Check if car exists
	car, err := s.carRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrCarNotFound) {
			return nil, ErrCarNotFound
		}
		return nil, err
	}

	// Update only provided fields
	if req.Name != "" {
		car.Name = req.Name
	}
	if req.EngineVersion != "" {
		car.EngineVersion = req.EngineVersion
	}

	if err := s.carRepo.Update(car); err != nil {
		if errors.Is(err, repository.ErrCarNotFound) {
			return nil, ErrCarNotFound
		}
		return nil, err
	}

	// Fetch updated car to get the new UpdatedAt timestamp
	updatedCar, err := s.carRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(updatedCar), nil
}

func (s *carService) DeleteCar(id uuid.UUID) error {
	err := s.carRepo.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrCarNotFound) {
			return ErrCarNotFound
		}
		return err
	}
	return nil
}

func (s *carService) entityToResponse(car *entity.Car) *dto.CarResponse {
	return &dto.CarResponse{
		ID:            car.ID,
		Name:          car.Name,
		EngineVersion: car.EngineVersion,
		CreatedAt:     car.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     car.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

var (
	ErrCarNotFound = errors.New("car not found")
)
