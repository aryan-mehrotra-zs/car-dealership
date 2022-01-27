package car

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
)

type store struct {
	db *sql.DB
}

func New(db *sql.DB) store {
	return store{db: db}
}

func (c store) Create(car models.Car) (uuid.UUID, error) {
	_, err := c.db.Exec("INSERT INTO cars (id,model,year_of_manufacture,brand,fuel_type,engine_id) VALUES (?,?,?,?,?,?)",
		car.ID, car.Model, car.YearOfManufacture, car.Brand, car.FuelType, car.ID)
	if err != nil {
		return uuid.Nil, errors.DB{}
	}
	return car.ID, nil
}

func (c store) GetAll(filter filters.Car) ([]models.Car, error) {
	var rows *sql.Rows
	var err error

	if filter.Brand == "" {
		rows, err = c.db.Query("SELECT * FROM cars WHERE brand=?;", filter.Brand)
	} else {
		rows, err = c.db.Query("SELECT * FROM cars;")
	}

	if err != nil {
		return nil, errors.DB{}
	}

	defer func() {
		err := rows.Close()
		if err != nil {

		}
	}()

	cars := make([]models.Car, 0)

	for rows.Next() {
		var car models.Car

		if err := rows.Scan(car.ID, car.Model, car.YearOfManufacture, car.Brand, car.FuelType, car.Engine.ID); err != nil {
			return nil, errors.DB{}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, errors.DB{}
	}

	return cars, nil
}

func (c store) GetByID(id uuid.UUID) (models.Car, error) {
	var car models.Car

	err := c.db.QueryRow("SELECT * FROM cars WHERE id = ?;", id.String())
	if err != nil {
		return models.Car{}, errors.DB{}
	}

	return car, nil
}

func (c store) Update(car models.Car) error {
	_, err := c.db.Exec("UPDATE cars SET `model=?,year_of_manufacture=?,brand=?,fuel_type=?,engine_id=?` WHERE id=?", car.Model, car.YearOfManufacture, car.Brand, car.FuelType, car.ID)

	if err != nil {
		return errors.DB{}
	}

	return nil
}

func (c store) Delete(id uuid.UUID) error {
	_, err := c.db.Exec("DELETE FROM cars WHERE id = ?;", id.String())
	if err != nil {
		return errors.DB{}
	}

	return nil
}
