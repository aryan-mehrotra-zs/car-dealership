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

func (e store) Create(engine models.Engine) (uuid.UUID, error) {
	_, err := e.db.Exec(insertEngine, engine.ID, engine.Displacement, engine.NCylinder, engine.Range)
	if err != nil {
		return uuid.Nil, errors.DB{Err: err}
	}

	return engine.ID, nil

}
func (e store) GetByID(id uuid.UUID) (models.Engine, error) {
	var engine models.Engine

	err := e.db.QueryRow(getEngine, id).
		Scan(&engine.ID, &engine.Displacement, &engine.NCylinder, &engine.Range)
	if err != nil {
		return models.Engine{}, errors.DB{Err: err}
	}

	return engine, nil
}

func (e store) Update(engine models.Engine) error {
	_, err := e.db.Exec(updateEngine, engine.Displacement, engine.NCylinder, engine.Range, engine.ID.String())
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
func (e store) Delete(id uuid.UUID) error {
	_, err := e.db.Exec(deleteEngine, id.String())
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
