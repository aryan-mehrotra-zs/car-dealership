package engine

import (
	"database/sql"
	goError "errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/stores"
)

func initializeTests(t *testing.T) (*sql.DB, sqlmock.Sqlmock, stores.Engine) {
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

	engine := models.Engine{
		ID:           id,
		Displacement: 200,
		NCylinder:    2,
		Range:        0,
	}

	queryError := goError.New("error in inserting")

	mock.ExpectExec(insertEngine).WithArgs(engine.ID, engine.Displacement, engine.NCylinder, engine.Range).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(insertEngine).WithArgs(engine.ID, engine.Displacement, engine.NCylinder, engine.Range).WillReturnError(queryError)

	cases := []struct {
		desc   string
		input  models.Engine
		output uuid.UUID
		err    error
	}{
		{"success case", engine, id, nil},
		{desc: "failure case", input: engine, output: uuid.Nil, err: errors.DB{queryError}},
	}

	for i, tc := range cases {
		id, err := s.Create(tc.input)

		if id != tc.output {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, id, tc.output)
		}

		if err != tc.err {
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

	queryError := goError.New("error in inserting")

	engine := models.Engine{
		ID:           id,
		Displacement: 200,
		NCylinder:    2,
		Range:        0,
	}

	rows := sqlmock.NewRows([]string{"id", "displacement", "no_of_cylinder", "range"}).
		AddRow(engine.ID, engine.Displacement, engine.NCylinder, engine.Range)

	mock.ExpectQuery(getEngine).WithArgs(id).WillReturnRows(rows)
	mock.ExpectQuery(getEngine).WithArgs(id).WillReturnError(queryError)

	cases := []struct {
		desc   string
		input  uuid.UUID
		output models.Engine
		err    error
	}{
		{"success", id, engine, nil},
		{"failure", id, models.Engine{}, errors.DB{Err: queryError}},
	}

	for i, tc := range cases {
		engine, err := s.GetByID(tc.input)

		if engine != tc.output {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, id, tc.output)
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

	insertError := goError.New("error in updating")

	engine := models.Engine{
		ID:           id,
		Displacement: 200,
		NCylinder:    2,
		Range:        0,
	}

	// why error cant be nil here
	mock.ExpectExec(updateEngine).WithArgs(engine.Displacement, engine.NCylinder, engine.Range, engine.ID.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(updateEngine).WithArgs(engine.Displacement, engine.NCylinder, engine.Range, engine.ID.String()).
		WillReturnError(insertError)

	cases := []struct {
		desc  string
		input models.Engine
		err   error
	}{
		{"success", engine, nil},
		{"failure", models.Engine{}, errors.DB{Err: insertError}},
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

	deleteError := goError.New("error in inserting")

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	mock.ExpectExec(deleteEngine).WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(deleteEngine).WithArgs(uuid.Nil).WillReturnError(deleteError)

	cases := []struct {
		desc string
		id   uuid.UUID
		err  error
	}{
		{"success", id, nil},
		{"failure", uuid.Nil, errors.DB{Err: deleteError}},
	}

	for i, tc := range cases {
		err := s.Delete(tc.id)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}

}
