package services

import (
	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
)

type Car interface {
	Create(car *models.Car) (*models.Car, error)
	GetAll(filter filters.Car) ([]models.Car, error)
	GetByID(id uuid.UUID) (*models.Car, error)
	Update(car *models.Car) (*models.Car, error)
	Delete(id uuid.UUID) error
}
