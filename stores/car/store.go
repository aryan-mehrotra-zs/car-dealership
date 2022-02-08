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

// Create inserts a new car in the database
func (s store) Create(car *models.Car) error {
	_, err := s.db.Exec(insertCar, car.ID, car.Model, car.ManufactureYear, car.Brand, car.FuelType, car.ID)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}

// GetAll fetches cars based on filter
func (s store) GetAll(filter filters.Car) ([]models.Car, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if filter.Brand != "" {
		rows, err = s.db.Query(getCarsWithBrand, filter.Brand)
	} else {
		rows, err = s.db.Query(getCars)
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

		if err := rows.Scan(&car.ID, &car.Model, &car.ManufactureYear, &car.Brand, &car.FuelType, &car.Engine.ID); err != nil {
			return nil, errors.DB{Err: err}
		}

		cars = append(cars, car)
	}

	return cars, nil
}

// GetByID fetches the car from database of the given id
func (s store) GetByID(id uuid.UUID) (models.Car, error) {
	var car models.Car

	err := s.db.QueryRow(getCar, id.String()).
		Scan(&car.ID, &car.Model, &car.ManufactureYear, &car.Brand, &car.FuelType, &car.Engine.ID)
	if err != nil {
		return models.Car{}, errors.DB{Err: err}
	}

	return car, nil
}

// Update modifies car of the given id
func (s store) Update(car *models.Car) error {
	_, err := s.db.Exec(updateCar, car.Model, car.ManufactureYear, car.Brand, car.FuelType, car.ID, car.ID)

	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}

// Delete removes car with the given id
func (s store) Delete(id uuid.UUID) error {
	_, err := s.db.Exec(deleteCar, id.String())
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
