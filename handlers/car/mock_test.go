package car

import (
	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
)

type mockService struct{}

func (m mockService) Create(car models.Car) (models.Car, error) {
	return models.Car{}, nil
}

func (m mockService) GetAll(filter filters.Car) ([]models.Car, error) {
	return nil, nil
}

func (m mockService) GetByID(id uuid.UUID) (models.Car, error) {
	return models.Car{}, nil
}

func (m mockService) Update(car models.Car) (models.Car, error) {
	return models.Car{}, nil
}

func (m mockService) Delete(id uuid.UUID) error {
	return nil
}

type mockReader struct{}

func (m mockReader) Read(p []byte) (n int, err error) {
	return 0, errors.InvalidParam{}
}
