package car

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/services"
	"github.com/amehrotra/car-dealership/stores"
	"github.com/amehrotra/car-dealership/types"
)

func initializeTest(t *testing.T) (services.Car, *stores.MockCar, *stores.MockEngine) {
	ctrl := gomock.NewController(t)

	mockCar := stores.NewMockCar(ctrl)
	mockEngine := stores.NewMockEngine(ctrl)

	service := New(mockEngine, mockCar)

	return service, mockCar, mockEngine
}

var engine = models.Engine{
	Displacement: 100,
	NCylinder:    2,
}

var car = models.Car{
	Model:           "X",
	ManufactureYear: 2020,
	Brand:           "BMW",
	Engine:          engine,
}

//func TestService_Create(t *testing.T) {
//	id, err := uuid.NewRandom()
//	if err != nil {
//		t.Errorf("error in creating id : %v", err)
//	}
//
//	car.ID = id
//	engine.ID = id
//
//	cases := []struct {
//		desc       string
//		mockCar    *models.Car
//		mockEngine *models.Engine
//		resp       *models.Car
//		err        error
//	}{
//		{"create successful", &car, &engine, &car, nil},
//		//{"invalid parameter", car2, models.Car{}, errors.InvalidParam{}},
//		//{"missing parameter", car3, models.Car{}, errors.MissingParam{}},
//	}
//
//	for i, tc := range cases {
//		s, mockCar, mockEngine := initializeTest(t)
//
//		mockEngine.EXPECT().Create(tc.mockEngine).Return(tc.mockEngine.ID, tc.err)
//		mockEngine.EXPECT().GetByID(tc.mockEngine.ID).Return(*tc.mockEngine, tc.err)
//
//		mockCar.EXPECT().Create(tc.mockCar).Return(tc.mockCar.ID, tc.err)
//
//		resp, err := s.Create(tc.mockCar)
//
//		if err != tc.err {
//			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
//		}
//
//		if tc.resp != resp {
//			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
//		}
//	}
//}

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
		},
	}

	cases := []struct {
		desc       string
		mockFilter filters.Car
		resp       []models.Car
		err        error
	}{
		{"received all cars", filters.Car{Brand: "BMW", Engine: true}, car, nil},
		{"received all cars without enginge", filters.Car{Brand: "BMW", Engine: false}, carWithoutEngine, nil},
	}

	for i, tc := range cases {
		s, mockCar, _ := initializeTest(t)

		mockCar.EXPECT().GetAll(tc.mockFilter).Return(tc.resp)

		resp, err := s.GetAll(tc.mockFilter)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if reflect.DeepEqual(tc.resp, resp) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_CarGetByID(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("error in creating id : %v", err)
	}

	car.ID = id
	engine.ID = id

	s, mockCar, mockEngine := initializeTest(t)

	mockCar.EXPECT().GetByID(id).Return(car, nil)
	mockEngine.EXPECT().GetByID(id).Return(engine, nil)

	car.Engine = engine

	resp, err := s.GetByID(id)

	if err != nil {
		t.Errorf("\n[TEST] Failed \nDesc received car\nGot %v\n Expected %v", err, nil)
	}

	if !reflect.DeepEqual(&car, resp) {
		t.Errorf("\n[TEST] Failed \nDesc received car\nGot %v\n Expected %v", resp, car)
	}
}

func TestService_CarGetByIDInvalidCar(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("error in creating id : %v", err)
	}

	car.ID = id
	engine.ID = id

	s, mockCar, _ := initializeTest(t)

	mockCar.EXPECT().GetByID(id).Return(car, errors.EntityNotFound{})

	car.Engine = engine

	resp, err := s.GetByID(id)

	if !reflect.DeepEqual(err, errors.EntityNotFound{}) {
		t.Errorf("\n[TEST] Failed \nDesc received car\nGot %v\n Expected %v", err, nil)
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc received car\nGot %v\n Expected %v", resp, nil)
	}
}

func TestService_CarGetByIDInvalidEngine(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("error in creating id : %v", err)
	}

	car.ID = id
	engine.ID = id

	s, mockCar, mockEngine := initializeTest(t)

	mockCar.EXPECT().GetByID(id).Return(car, nil)
	mockEngine.EXPECT().GetByID(id).Return(engine, errors.EntityNotFound{})

	car.Engine = engine

	resp, err := s.GetByID(id)

	if !reflect.DeepEqual(err, errors.EntityNotFound{}) {
		t.Errorf("\n[TEST] Failed \nDesc received car\nGot %v\n Expected %v", err, nil)
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc received car\nGot %v\n Expected %v", resp, nil)
	}
}

func TestService_Update(t *testing.T) {
	s, mockCar, mockEngine := initializeTest(t)

	mockEngine.EXPECT().Update(&engine).Return(nil)
	mockCar.EXPECT().Update(&car).Return(nil)

	resp, err := s.Update(&car)

	if err != nil {
		t.Errorf("\n[TEST] Failed \nDesc update successful\nGot %v\n Expected %v", err, nil)
	}

	if &car != resp {
		t.Errorf("\n[TEST] Failed \nDesc update successful\nGot %v\n Expected %v", resp, car)
	}
}

func TestService_UpdateInvalidEngine(t *testing.T) {
	s, _, mockEngine := initializeTest(t)

	mockEngine.EXPECT().Update(&engine).Return(errors.EntityNotFound{})

	resp, err := s.Update(&car)

	if !reflect.DeepEqual(err, errors.EntityNotFound{}) {
		t.Errorf("\n[TEST] Failed \nDesc update successful\nGot %v\n Expected %v", err, errors.EntityNotFound{})
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc update successful\nGot %v\n Expected %v", resp, car)
	}
}

func TestService_UpdateInvalidParam(t *testing.T) {
	s, _, mockEngine := initializeTest(t)

	mockEngine.EXPECT().Update(&engine).Return(nil)

	car.Brand = "Aryan"

	resp, err := s.Update(&car)

	if !reflect.DeepEqual(err, errors.InvalidParam{}) {
		t.Errorf("\n[TEST] Failed \nDesc invalid param\nGot %v\n Expected %v", err, errors.InvalidParam{})
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc invalid param\nGot %v\n Expected %v", resp, car)
	}
}

func TestService_UpdateInvalidCar(t *testing.T) {
	s, mockCar, mockEngine := initializeTest(t)

	mockEngine.EXPECT().Update(&engine).Return(nil)
	mockCar.EXPECT().Update(&car).Return(errors.EntityNotFound{})

	resp, err := s.Update(&car)

	if !reflect.DeepEqual(err, errors.EntityNotFound{}) {
		t.Errorf("\n[TEST] Failed \nDesc update successful\nGot %v\n Expected %v", err, errors.EntityNotFound{})
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc update successful\nGot %v\n Expected %v", resp, car)
	}
}

func TestService_Delete(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("error in creating id : %v", err)
	}

	s, mockCar, mockEngine := initializeTest(t)

	mockCar.EXPECT().Delete(id).Return(nil)
	mockEngine.EXPECT().Delete(id).Return(nil)

	err = s.Delete(id)

	if err != nil {
		t.Errorf("\n[TEST] Failed \nDesc delete success\nGot %v\n Expected nil", err)
	}
}

func TestService_DeleteInvalidCar(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("error in creating id : %v", err)
	}

	s, mockCar, _ := initializeTest(t)

	mockCar.EXPECT().Delete(id).Return(errors.EntityNotFound{})

	err = s.Delete(id)

	if !reflect.DeepEqual(err, errors.EntityNotFound{}) {
		t.Errorf("\n[TEST] Failed \nDesc update successful\nGot %v\n Expected %v", err, errors.EntityNotFound{})
	}
}

func TestService_DeleteInvalidEngine(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("error in creating id : %v", err)
	}

	s, mockCar, mockEngine := initializeTest(t)

	mockCar.EXPECT().Delete(id).Return(nil)
	mockEngine.EXPECT().Delete(id).Return(errors.EntityNotFound{})

	err = s.Delete(id)

	if !reflect.DeepEqual(err, errors.EntityNotFound{}) {
		t.Errorf("\n[TEST] Failed \nDesc update successful\nGot %v\n Expected %v", err, errors.EntityNotFound{})
	}
}

func Test_CheckCar(t *testing.T) {
	invalidEngine := models.Car{Model: "A", ManufactureYear: 2000, Brand: "tesla", FuelType: 3,
		Engine: models.Engine{Displacement: 200, NCylinder: 10, Range: 10}}
	invalidEngine2 := models.Car{Model: "A", ManufactureYear: 2000, Brand: "tesla", FuelType: 3,
		Engine: models.Engine{Displacement: -1, NCylinder: -1, Range: -1}}
	invalidEngine3 := models.Car{Model: "A", ManufactureYear: 2000, Brand: "tesla", FuelType: 3,
		Engine: models.Engine{Range: 0}}

	cases := []struct {
		desc  string
		input models.Car
		err   error
	}{
		{"invalid model", models.Car{Model: ""}, errors.InvalidParam{}},
		{"invalid year", models.Car{Model: "X", ManufactureYear: 1800}, errors.InvalidParam{}},
		{"invalid Brand", models.Car{Model: "Y", ManufactureYear: 2000, Brand: "suzuki"}, errors.InvalidParam{}},
		{"invalid fuel", models.Car{Model: "Z", ManufactureYear: 2000, Brand: "tesla", FuelType: 5}, errors.InvalidParam{}},
		{"invalid engine for petrol", invalidEngine, errors.InvalidParam{}},
		{"invalid engine for ev", invalidEngine2, errors.InvalidParam{}},
		{"invalid engine ", invalidEngine3, errors.InvalidParam{}},
	}

	for i, tc := range cases {
		err := checkCar(tc.input)

		if reflect.DeepEqual(err, tc.err) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func Test_checkEngine(t *testing.T) {
	invalidEngine1 := models.Engine{Displacement: 10, NCylinder: 10, Range: 10}
	invalidEngine2 := models.Engine{Displacement: 0, NCylinder: 0, Range: 0}
	invalidEngine3 := models.Engine{Displacement: -1, NCylinder: -1, Range: -1}

	validEngine1 := models.Engine{Displacement: 10, NCylinder: 20}
	validEngine2 := models.Engine{Displacement: 0, NCylinder: 0, Range: 10}

	cases := []struct {
		desc  string
		input models.Engine
		err   error
	}{
		{"all value more than 0", invalidEngine1, errors.InvalidParam{}},
		{"all value equal to 0", invalidEngine2, errors.InvalidParam{}},
		{"all value less than 0", invalidEngine3, errors.InvalidParam{}},
		{"valid engine for non electric", validEngine1, nil},
		{"valid engine for electric", validEngine2, nil},
	}

	for i, tc := range cases {
		err := checkEngine(tc.input)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}
