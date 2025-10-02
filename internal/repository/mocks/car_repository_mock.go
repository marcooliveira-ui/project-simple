package mocks

import (
	"project-simple/internal/domain/dto"
	"project-simple/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockCarRepository struct {
	mock.Mock
}

func (m *MockCarRepository) Create(car *entity.Car) error {
	args := m.Called(car)
	return args.Error(0)
}

func (m *MockCarRepository) FindByID(id uuid.UUID) (*entity.Car, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Car), args.Error(1)
}

func (m *MockCarRepository) FindAll(pagination *dto.PaginationRequest) ([]entity.Car, int64, error) {
	args := m.Called(pagination)
	return args.Get(0).([]entity.Car), args.Get(1).(int64), args.Error(2)
}

func (m *MockCarRepository) Update(car *entity.Car) error {
	args := m.Called(car)
	return args.Error(0)
}

func (m *MockCarRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCarRepository) ExistsByID(id uuid.UUID) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}
