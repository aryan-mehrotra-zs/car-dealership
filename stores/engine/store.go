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
	_, err := e.db.Exec("INSERT INTO engines (id,displacement,no_of_cylinder,`range`) VALUES (?,?,?,?)", engine.ID, engine.Displacement, engine.NCylinder, engine.Range)
	if err != nil {
		return uuid.Nil, err
	}

	return engine.ID, nil

}
func (e engine) GetByID(id uuid.UUID) (models.Engine, error) {
	var engine models.Engine

	// except id engine.ID will panic
	err := e.db.QueryRow("SELECT * FROM engines WHERE ID=?", id).
		Scan(engine.ID, engine.Displacement, engine.NCylinder, engine.Range)
	if err != nil {
		return models.Engine{}, err
	}

	return engine, nil
}

func (e engine) Update(engine models.Engine) error {
	_, err := e.db.Exec("UPDATE engine SET `displacement=?,no_of_cylinder=?,'range'=?` WHERE id=?", engine.Displacement, engine.NCylinder, engine.Range, engine.ID.String())
	if err != nil {
		return err
	}

	return nil
}
func (e engine) Delete(id uuid.UUID) error {
	_, err := e.db.Exec("DELETE FROM engines WHERE id = ?;", id.String())
	if err != nil {
		return err
	}

	return nil
}
