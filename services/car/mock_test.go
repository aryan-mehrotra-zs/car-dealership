package car

import (
	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/types"
)

type mockCar struct{}

type mockEngine struct{}

func (c mockCar) Create(car models.Car) (uuid.UUID, error) {
	switch car.ID {
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"):
		return car.ID, nil
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb3"):
		return uuid.Nil, errors.InvalidParam{}
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5"):
		return uuid.Nil, errors.MissingParam{}
	default:
		return uuid.Nil, errors.DB{}
	}
}
func (c mockCar) GetAll(filter filters.Car) ([]models.Car, error) {
	car := []models.Car{
		{
			ID:              uuid.Nil,
			Model:           "X",
			ManufactureYear: 2020,
			Brand:           "BMW",
			FuelType:        types.Petrol,
			Engine: models.Engine{
				ID:           uuid.Nil,
				Displacement: 100,
				NCylinder:    2,
				Range:        0,
			},
		},
	}

	carWithoutEngine := []models.Car{
		{
			ID:              uuid.Nil,
			Model:           "X",
			ManufactureYear: 2020,
			Brand:           "BMW",
			FuelType:        types.Petrol,
			Engine:          models.Engine{},
		},
	}

	switch {
	case filter.Engine == true && filter.Brand == "BMW":
		return car, nil
	case filter.Engine == false && filter.Brand == "BMW":
		return carWithoutEngine, nil
	default:
		return []models.Car{}, nil
	}
}
func (c mockCar) GetByID(id uuid.UUID) (models.Car, error) {
	car := models.Car{
		ID:              uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Model:           "X",
		ManufactureYear: 2020,
		Brand:           "BMW",
		FuelType:        0,
		Engine: models.Engine{
			Displacement: 100,
			NCylinder:    2,
			Range:        0,
		},
	}

	switch id {
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"):
		return car, nil
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb6"):
		return models.Car{}, errors.EntityNotFound{}
	default:
		return models.Car{}, nil
	}
}
func (c mockCar) Update(car models.Car) error {
	switch car.ID {
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"):
		return nil
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5"):
		return errors.EntityNotFound{}
	}

	return nil
}
func (c mockCar) Delete(id uuid.UUID) error {
	switch id {
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"):
		return nil
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5"):
		return errors.EntityNotFound{}
	default:
		return nil
	}
}

func (e mockEngine) Create(engine models.Engine) (uuid.UUID, error) {
	switch engine.ID {
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"):
		return uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"), nil
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5"):
		return uuid.Nil, errors.InvalidParam{}
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb6"):
		return uuid.Nil, errors.MissingParam{}
	default:
		return uuid.Nil, nil
	}
}
func (e mockEngine) GetByID(id uuid.UUID) (models.Engine, error) {
	engine := models.Engine{
		ID:           uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Displacement: 100,
		NCylinder:    2,
		Range:        0,
	}

	switch id {
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"):
		return engine, nil
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5"):
		return models.Engine{}, errors.EntityNotFound{}
	default:
		return models.Engine{}, nil
	}
}

func (e mockEngine) Update(engine models.Engine) error {
	switch engine.ID {
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"):
		return nil
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb6"):
		return errors.EntityNotFound{}
	default:
		return nil
	}
}

func (e mockEngine) Delete(id uuid.UUID) error {
	switch id {
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"):
		return nil
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb6"):
		return errors.EntityNotFound{}
	default:
		return nil
	}
}
