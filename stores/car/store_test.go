package car

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/stores"
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

	tests := []struct {
		desc  string
		input models.Car
		id    int
		err   error
	}{}

}
