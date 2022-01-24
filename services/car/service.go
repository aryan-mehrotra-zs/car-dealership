package car

import (
	"github.com/google/uuid"

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

func (s service) Create(car models.Car) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (s service) GetAll(brand string, engine bool) (models.Car, error) {
	return models.Car{}, nil
}

func (s service) GetByID(id uuid.UUID) (models.Car, error) {
	return models.Car{}, nil
}

func (s service) Update(car models.Car) error {
	return nil
}

func (s service) Delete(id uuid.UUID) error {
	return nil
}
