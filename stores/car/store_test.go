package car

import (
	"database/sql"
	goError "errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/stores"
	"github.com/amehrotra/car-dealership/types"
)

func initializeTests(t *testing.T) (*sql.DB, sqlmock.Sqlmock, stores.Car) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error %s was not expected when opening a stub database connection", err)
	}

	s := New(db)

	return db, mock, s
}

func TestStore_Create(t *testing.T) {
	db, mock, s := initializeTests(t)
	defer db.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	car := models.Car{
		ID:                id,
		Model:             "X",
		YearOfManufacture: 2020,
		Brand:             "BMW",
		FuelType:          0,
		Engine: models.Engine{
			ID:           id,
			Displacement: 100,
			NCylinder:    2,
			Range:        0,
		},
	}

	queryErr := goError.New("query error")

	mock.ExpectExec(insertCar).
		WithArgs(car.ID, car.Model, car.YearOfManufacture, car.Brand, car.FuelType, car.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(insertCar).
		WithArgs(car.ID, car.Model, car.YearOfManufacture, car.Brand, car.FuelType, car.ID).
		WillReturnError(queryErr)

	cases := []struct {
		desc string
		id   uuid.UUID
		err  error
	}{
		{"success case", car.ID, nil},
		{"failure case", uuid.Nil, errors.DB{Err: queryErr}},
	}

	for i, tc := range cases {
		id, err := s.Create(car)

		if id != tc.id {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, id, tc.id)
		}

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStore_GetAll(t *testing.T) {
	db, mock, s := initializeTests(t)
	defer db.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	cars := []models.Car{
		{
			ID:                id,
			Model:             "X",
			YearOfManufacture: 2020,
			Brand:             "BMW",
			FuelType:          types.Petrol,
			Engine: models.Engine{
				ID: id,
			},
		},
	}

	rows := sqlmock.NewRows([]string{"id", "model", "year_of_manufacture", "brand", "fuel_type", "engine_id"}).
		AddRow(id.String(), "X", 2020, "BMW", 1, id.String())

	mock.ExpectQuery(getCarsWithBrand).WithArgs("BMW").WillReturnRows(rows)
	//mock.ExpectQuery(getCars).WithArgs().WillReturnError(goError.New("error reading the row"))
	//mock.ExpectQuery(getCars).WithArgs().WillReturnError(goError.New("query error"))
	//mock.ExpectQuery(getCars).WithArgs().WillReturnError(goError.New("scan error"))

	cases := []struct {
		desc   string
		filter filters.Car // input
		output []models.Car
		err    error
	}{
		{"success case", filters.Car{Brand: "BMW"}, cars, nil},
		//{"rows error", filters.Car{}, nil, errors.DB{Err: goError.New("error reading the row")}},
		//{"query error", filters.Car{}, nil, errors.DB{Err: goError.New("query error")}},
		//{"scan error", filters.Car{}, nil, errors.DB{Err: goError.New("scan error")}},
	}

	for i, tc := range cases {
		cars, err := s.GetAll(tc.filter)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if !reflect.DeepEqual(cars, tc.output) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

	}
}

func TestStore_GetByID(t *testing.T) {
	db, mock, s := initializeTests(t)
	defer db.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	car := models.Car{
		ID:                id,
		Model:             "X",
		YearOfManufacture: 2020,
		Brand:             "BMW",
		FuelType:          0,
		Engine: models.Engine{
			ID:           id,
			Displacement: 100,
			NCylinder:    2,
			Range:        0,
		},
	}

	queryErr := goError.New("query error")

	rows := sqlmock.NewRows([]string{"id", "model", "year_of_manufacture", "brand", "fuel_type", "engine_id"}).
		AddRow(id.String(), "X", 2020, "BMW", 1, id.String())

	mock.ExpectQuery(getCar).WithArgs(id).WillReturnRows(rows)
	mock.ExpectQuery(getCar).WithArgs(id).WillReturnError(queryErr)

	cases := []struct {
		desc   string
		input  uuid.UUID
		output models.Car
		err    error
	}{
		{"success case", car.ID, car, nil},
		{"failure case", uuid.Nil, models.Car{}, errors.DB{Err: queryErr}},
	}

	for i, tc := range cases {
		resp, err := s.GetByID(tc.input)

		if resp != tc.output {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.output)
		}

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestStore_Update(t *testing.T) {
	db, mock, s := initializeTests(t)
	defer db.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	updateFailed := goError.New("update failed")

	car := models.Car{
		ID:                id,
		Model:             "X",
		YearOfManufacture: 2020,
		Brand:             "BMW",
		FuelType:          0,
		Engine: models.Engine{
			ID: id,
		},
	}

	mock.ExpectExec(updateCar).WithArgs(car.Model, car.YearOfManufacture, car.Brand, car.FuelType, car.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(updateCar).WithArgs(car.Model, car.YearOfManufacture, car.Brand, car.FuelType, car.ID).
		WillReturnError(updateFailed)

	cases := []struct {
		desc  string
		input models.Car
		err   error
	}{
		{"success", car, nil},
		{"failure", car, errors.DB{updateFailed}},
	}

	for i, tc := range cases {
		err := s.Update(tc.input)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestStore_Delete(t *testing.T) {
	db, mock, s := initializeTests(t)
	defer db.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	deleteErr := goError.New("delete failed")

	mock.ExpectExec(deleteCar).WithArgs(id.String()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(deleteCar).WithArgs(id.String()).WillReturnError(deleteErr)

	cases := []struct {
		desc string
		id   uuid.UUID
		err  error
	}{
		{"Delete Success", id, nil},
		{"Delete Failed", id, errors.DB{Err: deleteErr}},
	}

	for i, tc := range cases {
		err := s.Delete(tc.id)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}
