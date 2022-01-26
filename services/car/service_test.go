package car

import (
	"testing"

	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/models"
)

func TestService_Create(t *testing.T) {
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

	cases := []struct {
		desc   string
		input  models.Car
		output models.Car
		err    error
	}{
		{"create successful", car, car, nil},
		{"invalid parameter", car, models.Car{}, errors.InvalidParam{}},
		{"missing parameter", car, models.Car{}, errors.MissingParam{}},
	}

}

func TestService_GetAll(t *testing.T) {

}

func TestService_GetByID(t *testing.T) {

}

func TestService_Update(t *testing.T) {

}

func TestService_Delete(t *testing.T) {

}
