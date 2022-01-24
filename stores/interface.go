package stores

import (
	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/models"
)

type Car interface {
	Create(car models.Car) (uuid.UUID, error)
	GetAll(brand string, engine bool) (models.Car, error)
	GetByID(id uuid.UUID) (models.Car, error)
	Update(car models.Car) error
	Delete(id uuid.UUID) error
}

type Engine interface {
	Create(engine models.Engine) (uuid.UUID, error)
	GetByID(id uuid.UUID) (models.Engine, error)
	Update(car models.Engine) error
	Delete(id uuid.UUID) error
}
