package car

import (
	"bytes"
	"database/sql"
	goError "errors"
	"fmt"
	"log"
	"reflect"
	"strings"
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
		ID:              id,
		Model:           "X",
		ManufactureYear: 2020,
		Brand:           "BMW",
		Engine: models.Engine{
			ID:           id,
			Displacement: 100,
			NCylinder:    2,
		},
	}

	queryErr := goError.New("query error")

	mock.ExpectExec(insertCar).
		WithArgs(car.ID, car.Model, car.ManufactureYear, car.Brand, car.FuelType, car.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(insertCar).
		WithArgs(car.ID, car.Model, car.ManufactureYear, car.Brand, car.FuelType, car.ID).
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
		err := s.Create(&car)

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
			ID:              id,
			Model:           "X",
			ManufactureYear: 2020,
			Brand:           "BMW",
			FuelType:        types.Petrol,
			Engine:          models.Engine{ID: id},
		},
	}

	queryError := goError.New("query error")

	row1 := sqlmock.NewRows([]string{"id", "model", "year_of_manufacture", "brand", "fuel_type", "engine_id"}).
		AddRow(id.String(), "X", 2020, "BMW", []byte("petrol"), id.String())

	row2 := sqlmock.NewRows([]string{"id", "model", "year_of_manufacture", "brand", "fuel_type", "engine_id", "scan_error"}).
		AddRow(id.String(), "X", 2020, "BMW", []byte("petrol"), id.String(), "scan_error")

	mock.ExpectQuery(getCarsWithBrand).WithArgs("BMW").WillReturnRows(row1)
	mock.ExpectQuery(getCars).WillReturnError(queryError)
	mock.ExpectQuery(getCars).WillReturnRows(row2)

	cases := []struct {
		desc   string
		filter filters.Car // input
		output []models.Car
		err    error
	}{
		{"success case", filters.Car{Brand: "BMW"}, cars, nil},
		{"query error", filters.Car{}, nil, errors.DB{Err: queryError}},
		{"scan error", filters.Car{}, nil, errors.DB{Err: fmt.Errorf("sql: expected %d destination arguments in Scan, not %d", 7, 6)}},
	}

	for i, tc := range cases {
		car, err := s.GetAll(tc.filter)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if !reflect.DeepEqual(car, tc.output) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestStore_GetAllCloseErr(t *testing.T) {
	db, mock, s := initializeTests(t)
	defer db.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	closeError := goError.New("close error")
	rowError := goError.New("row error")

	closeRow := sqlmock.NewRows([]string{"id", "model", "year_of_manufacture", "brand", "fuel_type", "engine_id"}).
		AddRow(id.String(), "X", 2020, "BMW", []byte("petro"), id.String()).CloseError(errors.DB{Err: closeError})

	errRow := sqlmock.NewRows([]string{"id", "model", "year_of_manufacture", "brand", "fuel_type", "engine_id"}).
		AddRow(id.String(), "X", 2020, "BMW", []byte("petro"), id.String()).RowError(0, errors.DB{Err: rowError})

	mock.ExpectQuery(getCars).WillReturnRows(closeRow)
	mock.ExpectQuery(getCars).WillReturnRows(errRow)

	cases := []struct {
		desc string
		row  *sqlmock.Rows
		err  string
	}{
		{"row close error", closeRow, "close error"},
		{"row error", errRow, "row error"},
	}

	for i, tc := range cases {
		var b bytes.Buffer

		log.SetOutput(&b)

		car, err := s.GetAll(filters.Car{})

		if !strings.Contains(b.String(), tc.err) {
			t.Errorf("\n[TEST %d] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if len(car) != 0 {
			t.Errorf("\n[TEST %d] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, len(car), 0)
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
		ID:              id,
		Model:           "X",
		ManufactureYear: 2020,
		Brand:           "BMW",
		Engine:          models.Engine{ID: id},
	}

	queryErr := goError.New("query error")

	rows := sqlmock.NewRows([]string{"id", "model", "year_of_manufacture", "brand", "fuel_type", "engine_id"}).
		AddRow(id.String(), "X", 2020, "BMW", []byte("diesel"), id.String())

	mock.ExpectQuery(getCar).WithArgs(id).WillReturnRows(rows)
	mock.ExpectQuery(getCar).WithArgs(uuid.Nil).WillReturnError(queryErr)

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
		ID:              id,
		Model:           "X",
		ManufactureYear: 2020,
		Brand:           "BMW",
		Engine:          models.Engine{ID: id},
	}

	mock.ExpectExec(updateCar).WithArgs(car.Model, car.ManufactureYear, car.Brand, car.FuelType, car.ID, car.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(updateCar).WithArgs(car.Model, car.ManufactureYear, car.Brand, car.FuelType, car.ID, car.ID).
		WillReturnError(updateFailed)

	cases := []struct {
		desc  string
		input models.Car
		err   error
	}{
		{"success", car, nil},
		{"failure", car, errors.DB{Err: updateFailed}},
	}

	for i, tc := range cases {
		err := s.Update(&tc.input)

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
