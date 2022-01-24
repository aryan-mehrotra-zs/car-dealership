package engine

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/models"
)

type engine struct {
	db *sql.DB
}

func New(db *sql.DB) engine {
	return engine{db: db}
}

func (e engine) Create(engine models.Engine) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
func (e engine) GetByID(id uuid.UUID) (models.Engine, error) {
	return models.Engine{}, nil
}
func (e engine) Update(car models.Engine) error {
	return nil
}
func (e engine) Delete(id uuid.UUID) error {
	return nil
}
