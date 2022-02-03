package engine

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/stores"
)

type store struct {
	db *sql.DB
}

func New(db *sql.DB) stores.Engine {
	return store{db: db}
}

// Create inserts a new engine in the database
func (s store) Create(engine *models.Engine) (uuid.UUID, error) {
	_, err := s.db.Exec(insertEngine, engine.ID, engine.Displacement, engine.NCylinder, engine.Range)
	if err != nil {
		return uuid.Nil, errors.DB{Err: err}
	}

	return engine.ID, nil
}

// GetByID fetches the engine from database of the given id
func (s store) GetByID(id uuid.UUID) (models.Engine, error) {
	var engine models.Engine

	err := s.db.QueryRow(getEngine, id).
		Scan(&engine.ID, &engine.Displacement, &engine.NCylinder, &engine.Range)
	if err != nil {
		return models.Engine{}, errors.DB{Err: err}
	}

	return engine, nil
}

// Update modifies engine of the given id
func (s store) Update(engine *models.Engine) error {
	_, err := s.db.Exec(updateEngine, engine.Displacement, engine.NCylinder, engine.Range, engine.ID.String())
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}

// Delete removes engine with the given id
func (s store) Delete(id uuid.UUID) error {
	_, err := s.db.Exec(deleteEngine, id.String())
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
