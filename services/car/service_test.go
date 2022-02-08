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

//nolint
var engine = models.Engine{
	Displacement: 100,
	NCylinder:    2,
}

// nolint:gochecknoglobals // to remove redundant declaration in test file
var car = models.Car{
	Model:           "X",
	ManufactureYear: 2020,
	Brand:           "BMW",
	Engine:          engine,
}

func TestService_Create(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("error in creating id : %v", err)
	}

	car.ID = id
	engine.ID = id

	s, mockCar, mockEngine := initializeTest(t)

	mockEngine.EXPECT().Create(gomock.Any()).Return(nil)
	mockCar.EXPECT().Create(gomock.Any()).Return(nil)
	mockCar.EXPECT().GetByID(gomock.Any()).Return(car, nil)
	mockEngine.EXPECT().GetByID(gomock.Any()).Return(engine, nil)

	resp, err := s.Create(&car)

	if err != nil {
		t.Errorf("\n[TEST] Failed \nDesc create successful\nGot %v\n Expected %v", err, nil)
	}

	if !reflect.DeepEqual(resp, &car) {
		t.Errorf("\n[TEST] Failed \nDesc create successful\nGot %v\n Expected %v", resp, &car)
	}
}

func TestService_CreateInvalidCar(t *testing.T) {
	s, _, _ := initializeTest(t)

	resp, err := s.Create(&models.Car{})

	if !reflect.DeepEqual(err, errors.InvalidParam{Param: []string{"model"}}) {
		t.Errorf("\n[TEST] Failed \nDesc invalid car model\nGot %v\n Expected %v", err, errors.InvalidParam{})
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc invalid car model\nGot %v\n Expected %v", resp, nil)
	}
}

func TestService_CreateInvalidEngine(t *testing.T) {
	car = models.Car{
		ID:              uuid.Nil,
		Model:           "X",
		ManufactureYear: 2021,
		Brand:           "BMW",
		FuelType:        types.Petrol,
		Engine:          models.Engine{NCylinder: -11},
	}

	s, _, _ := initializeTest(t)

	resp, err := s.Create(&car)

	if !reflect.DeepEqual(err, errors.InvalidParam{Param: []string{"noOfCylinder"}}) {
		t.Errorf("\n[TEST] Failed \nDesc invalid engine parameter\nGot %v\n Expected %v", err,
			errors.InvalidParam{Param: []string{"noOfCylinder"}})
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc invalid engine parameter\nGot %v\n Expected %v", resp, nil)
	}
}

func TestService_CreateEngineDBError(t *testing.T) {
	s, _, mockEngine := initializeTest(t)

	mockEngine.EXPECT().Create(gomock.Any()).Return(errors.DB{})

	resp, err := s.Create(&car)

	if !reflect.DeepEqual(err, errors.DB{}) {
		t.Errorf("\n[TEST] Failed \nDesc db error when creating engine\nGot %v\n Expected %v", err, errors.DB{})
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc db error when creating engine\nGot %v\n Expected %v", resp, &car)
	}
}

func TestService_CreateVerificationError(t *testing.T) {
	s, mockCar, mockEngine := initializeTest(t)

	mockEngine.EXPECT().Create(gomock.Any()).Return(nil)
	mockCar.EXPECT().Create(gomock.Any()).Return(nil)
	mockCar.EXPECT().GetByID(gomock.Any()).Return(car, errors.DB{})

	resp, err := s.Create(&car)

	if !reflect.DeepEqual(err, errors.DB{}) {
		t.Errorf("\n[TEST] Failed \nDesc create successful\nGot %v\n Expected %v", err, errors.DB{})
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc create successful\nGot %v\n Expected %v", resp, &car)
	}
}

func TestService_CreateCarDBError(t *testing.T) {
	s, mockCar, mockEngine := initializeTest(t)

	mockEngine.EXPECT().Create(gomock.Any()).Return(nil)
	mockCar.EXPECT().Create(gomock.Any()).Return(errors.DB{})

	resp, err := s.Create(&car)

	if !reflect.DeepEqual(err, errors.DB{}) {
		t.Errorf("\n[TEST] Failed \nDesc create successful\nGot %v\n Expected %v", err, errors.DB{})
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc create successful\nGot %v\n Expected %v", resp, &car)
	}
}

func TestService_GetAllWithEngine(t *testing.T) {
	cars := []models.Car{
		{
			Model:           "X",
			ManufactureYear: 2020,
			Brand:           "BMW",
			FuelType:        types.Petrol,
			Engine:          models.Engine{Displacement: 100, NCylinder: 2},
		},
	}

	s, mockCar, mockEngine := initializeTest(t)

	mockCar.EXPECT().GetAll(gomock.Any()).Return(cars, nil)
	mockEngine.EXPECT().GetByID(gomock.Any()).Return(engine, nil)

	resp, err := s.GetAll(filters.Car{Brand: "BMW", Engine: true})

	if err != nil {
		t.Errorf("\n[TEST] Failed \nDesc received all cars\nGot %v\n Expected %v", err, nil)
	}

	if !reflect.DeepEqual(cars, resp) {
		t.Errorf("\n[TEST] Failed \nDesc received all cars\nGot %v\n Expected %v", resp, nil)
	}
}

func TestService_GetAllWithEngineDBError(t *testing.T) {
	cars := []models.Car{
		{
			Model:           "X",
			ManufactureYear: 2020,
			Brand:           "BMW",
			FuelType:        types.Petrol,
			Engine:          models.Engine{Displacement: 100, NCylinder: 2},
		},
	}

	s, mockCar, mockEngine := initializeTest(t)

	mockCar.EXPECT().GetAll(gomock.Any()).Return(cars, nil)
	mockEngine.EXPECT().GetByID(gomock.Any()).Return(engine, errors.DB{})

	resp, err := s.GetAll(filters.Car{Brand: "BMW", Engine: true})

	if !reflect.DeepEqual(err, errors.DB{}) {
		t.Errorf("\n[TEST] Failed \nDesc received all cars\nGot %v\n Expected %v", err, nil)
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc received all cars\nGot %v\n Expected %v", resp, nil)
	}
}

func TestService_GetAllWithoutEngine(t *testing.T) {
	cars := []models.Car{
		{
			ID:              uuid.Nil,
			Model:           "X",
			ManufactureYear: 2020,
			Brand:           "BMW",
			FuelType:        types.Petrol,
		},
	}

	cases := []struct {
		desc   string
		filter filters.Car
		resp   []models.Car
		err    error
	}{
		{"received all cars", filters.Car{Brand: "BMW"}, cars, nil},
		{"received all cars", filters.Car{Brand: ""}, cars, nil},
	}

	for i, tc := range cases {
		s, mockCar, _ := initializeTest(t)

		mockCar.EXPECT().GetAll(gomock.Any()).Return(cars, nil)

		resp, err := s.GetAll(tc.filter)

		if err != nil {
			t.Errorf("\n[TEST %d] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if !reflect.DeepEqual(cars, resp) {
			t.Errorf("\n[TEST %d] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_GetAllWithoutEngineDBError(t *testing.T) {
	s, mockCar, _ := initializeTest(t)

	mockCar.EXPECT().GetAll(gomock.Any()).Return(nil, errors.DB{})

	resp, err := s.GetAll(filters.Car{Brand: "BMW"})

	if !reflect.DeepEqual(err, errors.DB{}) {
		t.Errorf("\n[TEST] Failed \nDesc error in getting cars\nGot %v\n Expected %v", err, errors.DB{})
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc error in getting cars\nGot %v\n Expected %v", resp, nil)
	}
}

func TestService_GetAllInvalidBrand(t *testing.T) {
	s, _, _ := initializeTest(t)

	resp, err := s.GetAll(filters.Car{Brand: "Aryan"})

	if !reflect.DeepEqual(err, errors.InvalidParam{}) {
		t.Errorf("\n[TEST] Failed \nDesc received all cars\nGot %v\n Expected %v", err, errors.InvalidParam{})
	}

	if resp != nil {
		t.Errorf("\n[TEST] Failed \nDesc received all cars\nGot %v\n Expected %v", resp, nil)
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
	car = models.Car{
		Model:           "X",
		ManufactureYear: 2020,
		Brand:           "Aryan",
		Engine:          engine,
	}

	s, _, mockEngine := initializeTest(t)

	mockEngine.EXPECT().Update(&engine).Return(nil)

	resp, err := s.Update(&car)

	if !reflect.DeepEqual(err, errors.InvalidParam{Param: []string{"brand"}}) {
		t.Errorf("\n[TEST] Failed \nDesc invalid param\nGot %v\n Expected %v", err, errors.InvalidParam{Param: []string{"brand"}})
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
		err := checkCar(&tc.input)

		if reflect.DeepEqual(err, tc.err) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func Test_checkEngine(t *testing.T) {
	invalidEngine1 := models.Engine{Displacement: 10, NCylinder: 10, Range: 10}
	invalidEngine2 := models.Engine{Displacement: 0, NCylinder: 0, Range: 0}
	invalidEngine3 := models.Engine{Displacement: -1, NCylinder: -1, Range: -1}
	invalidEngine4 := models.Engine{Displacement: 0, NCylinder: 0, Range: -1}

	validEngine1 := models.Engine{Displacement: 10, NCylinder: 20}
	validEngine2 := models.Engine{Displacement: 0, NCylinder: 0, Range: 10}

	cases := []struct {
		desc  string
		input models.Engine
		err   error
	}{
		{"all value more than 0", invalidEngine1, errors.InvalidParam{Param: []string{"displacement", "noOfCylinder", "range"}}},
		{"all value equal to 0", invalidEngine2, errors.InvalidParam{Param: []string{"displacement", "noOfCylinder", "range"}}},
		{"all value less than 0", invalidEngine3, errors.InvalidParam{Param: []string{"displacement"}}},
		{"all value less than 0", invalidEngine4, errors.InvalidParam{Param: []string{"range"}}},
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
