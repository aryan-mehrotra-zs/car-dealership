package engine

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/models"
)

type store struct {
	db *sql.DB
}

func New(db *sql.DB) store {
	return store{db: db}
}

func (e store) Create(engine models.Engine) (uuid.UUID, error) {
	_, err := e.db.Exec("INSERT INTO engines (id,displacement,no_of_cylinder,`range`) VALUES (?,?,?,?)", engine.ID, engine.Displacement, engine.NCylinder, engine.Range)
	if err != nil {
		return uuid.Nil, errors.DB{}
	}

	return engine.ID, nil

}
func (e store) GetByID(id uuid.UUID) (models.Engine, error) {
	var engine models.Engine

	// except id store.ID will panic
	err := e.db.QueryRow("SELECT * FROM engines WHERE ID=?", id).
		Scan(engine.ID, engine.Displacement, engine.NCylinder, engine.Range)
	if err != nil {
		return models.Engine{}, errors.DB{}
	}

	return engine, nil
}

func (e store) Update(engine models.Engine) error {
	_, err := e.db.Exec("UPDATE store SET `displacement=?,no_of_cylinder=?,'range'=?` WHERE id=?", engine.Displacement, engine.NCylinder, engine.Range, engine.ID.String())
	if err != nil {
		return errors.DB{}
	}

	return nil
}
func (e store) Delete(id uuid.UUID) error {
	_, err := e.db.Exec("DELETE FROM engines WHERE id = ?;", id.String())
	if err != nil {
		return errors.DB{}
	}

	return nil
}
