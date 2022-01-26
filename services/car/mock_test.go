package car

import (
	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
)

type mockCar struct{}

type mockEngine struct{}

func (c mockCar) Create(car models.Car) (uuid.UUID, error) {
	return uuid.Nil, nil
}
func (c mockCar) GetAll(filter filters.Car) ([]models.Car, error) {
	return nil, nil
}
func (c mockCar) GetByID(id uuid.UUID) (models.Car, error) {
	return models.Car{}, nil
}
func (c mockCar) Update(car models.Car) error {
	return nil
}
func (c mockCar) Delete(id uuid.UUID) error {
	return nil
}

func (e mockEngine) Create(engine models.Engine) (uuid.UUID, error) {
	return uuid.Nil, nil
}
func (e mockEngine) GetByID(id uuid.UUID) (models.Engine, error) {
	return models.Engine{}, nil
}
func (e mockEngine) Update(car models.Engine) error {
	return nil
}
func (e mockEngine) Delete(id uuid.UUID) error {
	return nil
}
