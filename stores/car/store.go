package car

import (
	"database/sql"
	"log"

	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/stores"
)

type store struct {
	db *sql.DB
}

func New(db *sql.DB) stores.Car {
	return store{db: db}
}

func (c store) Create(car models.Car) (uuid.UUID, error) {
	_, err := c.db.Exec(insertCar, car.ID, car.Model, car.YearOfManufacture, car.Brand, car.FuelType, car.ID)
	if err != nil {
		return uuid.Nil, errors.DB{Err: err}
	}
	return car.ID, nil
}

func (c store) GetAll(filter filters.Car) ([]models.Car, error) {
	var rows *sql.Rows
	var err error

	if filter.Brand != "" {
		rows, err = c.db.Query(getCarsWithBrand, filter.Brand)
	} else {
		rows, err = c.db.Query(getCars)
	}

	if err != nil {
		return nil, errors.DB{Err: err}
	}

	defer func() {
		if err := rows.Err(); err != nil {
			log.Printf("error in accessing all rows: %v", err)
		}
	}()

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("error in closing rows : %v", err)
		}
	}()

	cars := make([]models.Car, 0)

	for rows.Next() {
		var car models.Car

		if err := rows.Scan(&car.ID, &car.Model, &car.YearOfManufacture, &car.Brand, &car.FuelType, &car.Engine.ID); err != nil {
			return nil, errors.DB{Err: err}
		}

		cars = append(cars, car)
	}

	return cars, nil
}

func (c store) GetByID(id uuid.UUID) (models.Car, error) {
	var car models.Car

	err := c.db.QueryRow(getCar, id.String()).
		Scan(&car.ID, &car.Model, &car.YearOfManufacture, &car.Brand, &car.FuelType, &car.Engine.ID)
	if err != nil {
		return models.Car{}, errors.DB{Err: err}
	}

	return car, nil
}

func (c store) Update(car models.Car) error {
	_, err := c.db.Exec(updateCar, car.Model, car.YearOfManufacture, car.Brand, car.FuelType, car.ID)

	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}

func (c store) Delete(id uuid.UUID) error {
	_, err := c.db.Exec(deleteCar, id.String())
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
