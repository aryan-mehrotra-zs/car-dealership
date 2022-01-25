package car

import (
	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/stores"
)

type service struct {
	engine stores.Engine
	car    stores.Car
}

func New(engine stores.Engine, car stores.Car) service {
	return service{engine: engine, car: car}
}

func (s service) Create(car models.Car) (models.Car, error) {
	return car, nil
}

func (s service) GetAll(car filters.Car) ([]models.Car, error) {
	return nil, nil
}

func (s service) GetByID(id uuid.UUID) (models.Car, error) {
	return models.Car{}, nil
}

func (s service) Update(car models.Car) (models.Car, error) {
	return models.Car{}, nil
}

func (s service) Delete(id uuid.UUID) error {
	return nil
}
