package car

import (
	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/types"
)

type mockService struct{}

func (m mockService) Create(car models.Car) (models.Car, error) {
	switch car.Model {
	case "x":
		return car, nil
	case "y":
		return models.Car{}, errors.EntityAlreadyExists{}
	case "z":
		return models.Car{}, errors.DB{}
	default:
		return models.Car{}, nil
	}

}

func (m mockService) GetAll(filter filters.Car) ([]models.Car, error) {
	withEngine := []models.Car{
		{
			ID:                uuid.Nil,
			Model:             "X",
			YearOfManufacture: 2020,
			Brand:             "BMW",
			FuelType:          types.Petrol,
			Engine: models.Engine{
				ID:           uuid.Nil,
				Displacement: 100,
				NCylinder:    2,
				Range:        0,
			},
		},
	}

	withoutEngine := []models.Car{
		{
			ID:                uuid.Nil,
			Model:             "X",
			YearOfManufacture: 2020,
			Brand:             "BMW",
			FuelType:          types.Petrol,
			Engine:            models.Engine{},
		},
	}

	switch filter {
	case filters.Car{Brand: "BMW", Engine: true}:
		return withEngine, nil
	case filters.Car{Brand: "BMW", Engine: false}:
		return withoutEngine, nil
	case filters.Car{Brand: "xyz", Engine: true}:
		return nil, errors.InvalidParam{}
	default:
		return nil, nil
	}

}

func (m mockService) GetByID(id uuid.UUID) (models.Car, error) {
	car := models.Car{
		ID:                uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Model:             "X",
		YearOfManufacture: 2020,
		Brand:             "BMW",
		FuelType:          0,
		Engine: models.Engine{
			Displacement: 100,
			NCylinder:    2,
			Range:        0,
		},
	}

	switch id {
	case uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"):
		return car, nil
	case uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"):
		return models.Car{}, errors.EntityNotFound{}
	case uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"):
		return models.Car{}, errors.InvalidParam{}
	case uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"):
		return models.Car{}, errors.DB{}
	default:
		return models.Car{}, errors.DB{}

	}
}

func (m mockService) Update(car models.Car) (models.Car, error) {
	car2 := models.Car{
		ID:                uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Model:             "X",
		YearOfManufacture: 2020,
		Brand:             "BMW",
		FuelType:          0,
		Engine: models.Engine{
			Displacement: 100,
			NCylinder:    2,
			Range:        0,
		},
	}

	switch car.ID {
	case uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"):
		return car2, nil
	case uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"):
		return models.Car{}, errors.EntityNotFound{}
	case uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"):
		return models.Car{}, errors.DB{}
	default:
		return models.Car{}, errors.DB{}
	}
}

func (m mockService) Delete(id uuid.UUID) error {
	switch id {
	case uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"):
		return nil
	case uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"):
		return errors.EntityNotFound{}
	case uuid.Nil:
		return errors.InvalidParam{}
	default:
		return errors.DB{}
	}
}
