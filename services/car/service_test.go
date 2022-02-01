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

func TestService_Create(t *testing.T) {
	car := models.Car{
		ID:              uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Model:           "X",
		ManufactureYear: 2020,
		Brand:           "BMW",
		FuelType:        0,
		Engine: models.Engine{
			ID:           uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
			Displacement: 100,
			NCylinder:    2,
			Range:        0,
		},
	}

	car2 := models.Car{
		ID:              uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb3"),
		Model:           "X",
		ManufactureYear: 0020,
		Brand:           "BMW",
		FuelType:        0,
		Engine: models.Engine{
			Displacement: 100,
			NCylinder:    2,
			Range:        0,
		},
	}

	car3 := models.Car{
		ID:              uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5"),
		Model:           "X",
		ManufactureYear: 2020,
		Brand:           "Tesla",
		FuelType:        1,
		Engine: models.Engine{
			Displacement: 100,
			NCylinder:    2,
			Range:        0,
		},
	}

	cases := []struct {
		desc  string
		input models.Car
		resp  models.Car
		err   error
	}{
		{"create successful", car, car, nil},
		{"invalid parameter", car2, models.Car{}, errors.InvalidParam{}},
		{"missing parameter", car3, models.Car{}, errors.MissingParam{}},
	}

	s := New(mockEngine{}, mockCar{})

	for i, tc := range cases {
		resp, err := s.Create(tc.input)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if tc.resp != resp {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_GetAll(t *testing.T) {
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
		resp, err := s.GetAll(tc.input)

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
		ID:              uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Model:           "X",
		ManufactureYear: 2020,
		Brand:           "BMW",
		FuelType:        0,
		Engine: models.Engine{
			ID:           uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
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
		resp, err := s.GetByID(tc.id)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if tc.resp != resp {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_Update(t *testing.T) {
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

	car2 := models.Car{
		ID: uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5"),
	}

	cases := []struct {
		desc string
		car  models.Car
		resp models.Car
		err  error
	}{
		{"update successful", car, car, nil},
		{"id not found", car2, models.Car{}, errors.EntityNotFound{}},
	}

	s := New(mockEngine{}, mockCar{})

	for i, tc := range cases {
		resp, err := s.Update(tc.car)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if tc.resp != resp {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_Delete(t *testing.T) {
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
		err := s.Delete(tc.id)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func Test_CheckCar(t *testing.T) {
	invalidEngine := models.Car{Model: "A", ManufactureYear: 2000, Brand: "tesla", FuelType: 3,
		Engine: models.Engine{Displacement: 200, NCylinder: 10, Range: 10}}
	invalidEngine2 := models.Car{Model: "A", ManufactureYear: 2000, Brand: "tesla", FuelType: 3,
		Engine: models.Engine{Displacement: -1, NCylinder: -1, Range: -1}}
	invalidEngine3 := models.Car{Model: "A", ManufactureYear: 2000, Brand: "tesla", FuelType: 3,
		Engine: models.Engine{Displacement: 0, NCylinder: 0, Range: 0}}

	cases := []struct {
		desc  string
		input models.Car
		err   error
	}{
		{"invalid model", models.Car{Model: ""}, errors.InvalidParam{}},
		{"invalid year", models.Car{Model: "X", ManufactureYear: 1800}, errors.InvalidParam{}},
		{"invalid Brand", models.Car{Model: "Y", ManufactureYear: 2000, Brand: "suzuki"}, errors.InvalidParam{}},
		{"invalid Fuel", models.Car{Model: "Z", ManufactureYear: 2000, Brand: "tesla", FuelType: 5}, errors.InvalidParam{}},
		{"invalid engine for petrol", invalidEngine, errors.InvalidParam{}},
		{"invalid engine for ev", invalidEngine2, errors.InvalidParam{}},
		{"invalid engine ", invalidEngine3, errors.InvalidParam{}},
	}

	for i, tc := range cases {
		err := checkCar(tc.input)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}
