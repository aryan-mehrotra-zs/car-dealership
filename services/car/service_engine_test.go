package car

import (
	"testing"

	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/models"
)

func TestService_EngineCreate(t *testing.T) {
	engine := models.Engine{
		ID:           uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Displacement: 200,
		NCylinder:    3,
		Range:        0,
	}

	cases := []struct {
		desc  string
		input models.Engine
		resp  uuid.UUID
		err   error
	}{
		{"create successful", engine, uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"), nil},
		{"invalid parameter", models.Engine{ID: uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5")}, uuid.Nil, errors.InvalidParam{}},
		{"missing parameter", models.Engine{ID: uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb6")}, uuid.Nil, errors.MissingParam{}},
	}

	s := New(mockEngine{}, mockCar{})

	for i, tc := range cases {
		resp, err := s.engine.Create(tc.input)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if tc.resp != resp {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_EngineGetByID(t *testing.T) {
	engine := models.Engine{
		ID:           uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Displacement: 200,
		NCylinder:    3,
		Range:        0,
	}

	cases := []struct {
		desc string
		id   uuid.UUID
		resp models.Engine
		err  error
	}{
		{"received car", uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"), engine, nil},
		{"id not found", uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb5"), models.Engine{}, errors.EntityNotFound{}},
	}

	s := New(mockEngine{}, mockCar{})

	for i, tc := range cases {
		resp, err := s.engine.GetByID(tc.id)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if tc.resp != resp {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_EngineUpdate(t *testing.T) {
	engine := models.Engine{
		ID:           uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Displacement: 200,
		NCylinder:    3,
		Range:        0,
	}

	cases := []struct {
		desc string
		car  models.Engine
		err  error
	}{
		{"update successful", engine, nil},
		{"id not found", models.Engine{ID: uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb6")}, errors.EntityNotFound{}},
	}

	s := New(mockEngine{}, mockCar{})

	for i, tc := range cases {
		err := s.engine.Update(tc.car)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestService_EngineDelete(t *testing.T) {
	cases := []struct {
		desc string
		id   uuid.UUID
		err  error
	}{
		{"delete success", uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"), nil},
		{"invalid id", uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb6"), errors.EntityNotFound{}},
	}

	s := New(mockEngine{}, mockCar{})

	for i, tc := range cases {
		err := s.engine.Delete(tc.id)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}
