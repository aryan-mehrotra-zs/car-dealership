package car

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/types"
)

func TestService_CarCreate(t *testing.T) {
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

	car2 := models.Car{
		ID:                uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb3"),
		Model:             "X",
		YearOfManufacture: 0020,
		Brand:             "BMW",
		FuelType:          0,
		Engine: models.Engine{
			Displacement: 100,
			NCylinder:    2,
			Range:        0,
		},
	}

	car3 := models.Car{
		ID:                uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5"),
		Model:             "X",
		YearOfManufacture: 2020,
		Brand:             "",
		FuelType:          0,
		Engine: models.Engine{
			Displacement: 100,
			NCylinder:    2,
			Range:        0,
		},
	}

	cases := []struct {
		desc  string
		input models.Car
		resp  uuid.UUID
		err   error
	}{
		{"create successful", car, uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"), nil},
		{"invalid parameter", car2, uuid.Nil, errors.InvalidParam{}},
		{"missing parameter", car3, uuid.Nil, errors.MissingParam{}},
	}

	s := New(mockEngine{}, mockCar{})

	for i, tc := range cases {
		resp, err := s.car.Create(tc.input)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if tc.resp != resp {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_CarGetAll(t *testing.T) {
	car := []models.Car{
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

	carWithoutEngine := []models.Car{
		{
			ID:                uuid.Nil,
			Model:             "X",
			YearOfManufacture: 2020,
			Brand:             "BMW",
			FuelType:          types.Petrol,
			Engine:            models.Engine{},
		},
	}

	cases := []struct {
		desc  string
		input filters.Car
		resp  []models.Car
		err   error
	}{
		{"received all cars", filters.Car{Brand: "BMW", Engine: true}, car, nil},
		{"received all cars without enginge", filters.Car{Brand: "BMW", Engine: false}, carWithoutEngine, nil},
	}

	s := New(mockEngine{}, mockCar{})

	for i, tc := range cases {
		resp, err := s.car.GetAll(tc.input)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if reflect.DeepEqual(tc.resp, resp) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_CarGetByID(t *testing.T) {
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
		desc string
		id   uuid.UUID
		resp models.Car
		err  error
	}{
		{"received car", uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"), car, nil},
		{"id not found", uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb6"), models.Car{}, errors.EntityNotFound{}},
	}

	s := New(mockEngine{}, mockCar{})

	for i, tc := range cases {
		resp, err := s.car.GetByID(tc.id)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if tc.resp != resp {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_CarUpdate(t *testing.T) {
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

	car2 := models.Car{
		ID: uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5"),
	}

	cases := []struct {
		desc string
		car  models.Car
		err  error
	}{
		{"update successful", car, nil},
		{"id not found", car2, errors.EntityNotFound{}},
	}

	s := New(mockEngine{}, mockCar{})

	for i, tc := range cases {
		err := s.car.Update(tc.car)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestService_CarDelete(t *testing.T) {
	cases := []struct {
		desc string
		id   uuid.UUID
		err  error
	}{
		{"delete success", uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"), nil},
		{"invalid id", uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5"), errors.EntityNotFound{}},
	}

	s := New(mockEngine{}, mockCar{})

	for i, tc := range cases {
		err := s.car.Delete(tc.id)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}
