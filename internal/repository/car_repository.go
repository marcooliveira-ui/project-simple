package repository

import (
	"errors"
	"project-simple/internal/domain/dto"
	"project-simple/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CarRepository interface {
	Create(car *entity.Car) error
	FindByID(id uuid.UUID) (*entity.Car, error)
	FindAll(pagination *dto.PaginationRequest) ([]entity.Car, int64, error)
	Update(car *entity.Car) error
	Delete(id uuid.UUID) error
	ExistsByID(id uuid.UUID) (bool, error)
}

type carRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) CarRepository {
	return &carRepository{db: db}
}

func (r *carRepository) Create(car *entity.Car) error {
	return r.db.Create(car).Error
}

func (r *carRepository) FindByID(id uuid.UUID) (*entity.Car, error) {
	var car entity.Car
	err := r.db.Where("id = ?", id).First(&car).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCarNotFound
		}
		return nil, err
	}
	return &car, nil
}

func (r *carRepository) FindAll(pagination *dto.PaginationRequest) ([]entity.Car, int64, error) {
	var cars []entity.Car
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Car{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records with sorting
	err := r.db.
		Order(pagination.GetOrderBy()).
		Limit(pagination.PageSize).
		Offset(pagination.GetOffset()).
		Find(&cars).Error

	if err != nil {
		return nil, 0, err
	}

	return cars, total, nil
}

func (r *carRepository) Update(car *entity.Car) error {
	result := r.db.Model(&entity.Car{}).
		Where("id = ?", car.ID).
		Updates(map[string]interface{}{
			"name":           car.Name,
			"engine_version": car.EngineVersion,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrCarNotFound
	}

	return nil
}

func (r *carRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entity.Car{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrCarNotFound
	}

	return nil
}

func (r *carRepository) ExistsByID(id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Car{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

var (
	ErrCarNotFound = errors.New("car not found")
)
