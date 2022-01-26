package car

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
)

type car struct {
	db *sql.DB
}

func New(db *sql.DB) car {
	return car{db: db}
}

func (c car) Create(car models.Car) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (c car) GetAll(filter filters.Car) ([]models.Car, error) {

	return nil, nil
}

func (c car) GetByID(id uuid.UUID) (models.Car, error) {
	return models.Car{}, nil
}

func (c car) Update(car models.Car) error {
	return nil
}

func (c car) Delete(id uuid.UUID) error {
	return nil
}
