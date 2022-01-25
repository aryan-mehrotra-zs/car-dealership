package car

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
)

type mockService struct{}

func (m mockService) Create(car models.Car) (models.Car, error) {
	return models.Car{}, nil
}

func (m mockService) GetAll(param filters.Car) ([]models.Car, error) {
	return nil, nil
}

func (m mockService) GetByID(id uuid.UUID) (models.Car, error) {
	return models.Car{}, nil
}

func (m mockService) Update(car models.Car) (models.Car, error) {
	return models.Car{}, nil
}

func (m mockService) Delete(id uuid.UUID) error {
	return nil
}

type mockReader struct{}

func (m mockReader) Read(p []byte) (n int, err error) {
	return 0, errors.InvalidParam{}
}

func TestHandler_Create(t *testing.T) {
	resp := []byte(`("ID":"123e4567-e89b-12d3-a456-426614174000","Model":"BMW","YearOfManufacture":2022,"Brand":"BMW",
					"FuelType":"Petrol",{"Displacement":20,"NCylinder":2,"Range":0}}`)

	cases := []struct {
		desc       string
		body       io.Reader
		car        []byte
		statusCode int
	}{
		{"created successfully", bytes.NewReader(resp), resp, http.StatusCreated},
		{"entity already exists", bytes.NewReader(resp), resp, http.StatusOK},
		{"unmarshal error", bytes.NewReader([]byte(`Invalid Body`)), []byte(""), http.StatusBadRequest},
		{"database error", bytes.NewReader([]byte(`{"Model":"BMW","YearOfManufacture":2022,"Brand":"BMW","FuelType":"Petrol",{20,2,0}}`)), []byte(""), http.StatusInternalServerError},
		{"missing parameter", bytes.NewReader([]byte(`{"Model":"BMW","YearOfManufacture":2022,"Brand":"BMW"}`)), []byte(""), http.StatusBadRequest},
		{"invalid parameter", bytes.NewReader([]byte(`{"Model":BMW,"YearOfManufacture":2022,"Brand":"BMW","FuelType":"Petrol",{20,2,0}}`)), []byte(""), http.StatusBadRequest},
		{"unable to read body", mockReader{}, []byte(""), http.StatusBadRequest},
		{"database connectivity error", bytes.NewReader(resp), resp, http.StatusInternalServerError},
	}

	h := New(mockService{})

	for i, tc := range cases {
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8000/car", tc.body)
		w := httptest.NewRecorder()

		h.Create(w, r)

		resp := w.Result()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("cannot read response")
		}

		err = resp.Body.Close()
		if err != nil {
			t.Errorf("Error in closing body")
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if reflect.DeepEqual(body, resp.Body) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), string(tc.car))
		}
	}
}

func TestHandler_GetAll(t *testing.T) {
	withoutEngine := []byte(`[{"ID":"123e4567-e89b-12d3-a456-426614174000","Model":"BMW","YearOfManufacture":2022,"Brand":"BMW",
						"FuelType":"Petrol"}]`)
	withEngine := []byte(`[{"ID":"123e4567-e89b-12d3-a456-426614174000","Model":"BMW","YearOfManufacture":2022,"Brand":"BMW",
						"FuelType":"Petrol",{"Displacement":20,"NCylinder":2,"Range":0}}]`)
	allWithoutEngine := []byte(`[{"ID":"123e4567-e89b-12d3-a456-426614174000","Model":"BMW","YearOfManufacture":2022,"Brand":"BMW",
						"FuelType":"Petrol"}
						{"ID":"123e4567-e89b-12d3-a457-426614174000","Model":"Mercedes","YearOfManufacture":2022,"Brand":"BMW",
						"FuelType":"Petrol"}]`)
	allWithEngine := []byte(`[{"ID":"123e4567-e89b-12d3-a456-426614174000","Model":"BMW","YearOfManufacture":2022,"Brand":"BMW",
						"FuelType":"Petrol",{"Displacement":20,"NCylinder":2,"Range":0}}
						{"ID":"123e4567-e89b-12d3-a456-427614174000","Model":"BMW","YearOfManufacture":2022,"Brand":"BMW",
						"FuelType":"Petrol"}
						{"ID":"123e4567-e89b-12d3-a457-426614174000","Model":"Mercedes","YearOfManufacture":2022,"Brand":"BMW",
						"FuelType":"Petrol",{"Displacement":20,"NCylinder":2,"Range":0}}]`)

	cases := []struct {
		desc       string
		car        filters.Car
		resp       []byte
		statusCode int
	}{
		{"get all cars of a brand with engine", filters.Car{Brand: "BMW", Engine: true}, withoutEngine, http.StatusOK},
		{"get all cars without engine", filters.Car{Brand: "BMW", Engine: false}, withEngine, http.StatusOK},
		{"get all cars from all brands without engine", filters.Car{Brand: "", Engine: true}, allWithoutEngine, http.StatusOK},
		{"get all cars from all brands with engine", filters.Car{Brand: "", Engine: true}, allWithEngine, http.StatusOK},
		{"invalid brand name", filters.Car{Brand: "xyz", Engine: true}, []byte(""), http.StatusBadRequest},
		{"database connectivity error", filters.Car{Brand: "", Engine: true}, []byte(""), http.StatusInternalServerError},
	}

	h := New(mockService{})

	for i, tc := range cases {
		req := httptest.NewRequest(http.MethodPost, "http://localhost:8000/car", http.NoBody)
		r := mux.SetURLVars(req, map[string]string{"brand": tc.car.Brand, "engine": strconv.FormatBool(tc.car.Engine)})
		w := httptest.NewRecorder()

		h.GetAll(w, r)

		resp := w.Result()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("cannot read response")
		}

		resp.Body.Close()

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if reflect.DeepEqual(body, tc.resp) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), string(tc.resp))
		}
	}
}

func TestHandler_GetByID(t *testing.T) {
	resp := []byte(`{"ID":"123e4567-e89b-12d3-a456-426614174000","Model":"BMW","YearOfManufacture":2022,"Brand":"BMW",
					"FuelType":"Petrol"}`)

	cases := []struct {
		desc       string
		id         uuid.UUID
		resp       []byte
		statusCode int
	}{
		{"request successful", uuid.UUID{}, resp, http.StatusOK},
		{"entity not found", uuid.UUID{}, []byte(""), http.StatusBadRequest},
		{"invalid id", uuid.Nil, []byte(""), http.StatusBadRequest},
		{"database connectivity error", uuid.UUID{}, []byte(""), http.StatusInternalServerError},
	}

	h := New(mockService{})

	for i, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8000/car", http.NoBody)
		r := mux.SetURLVars(req, map[string]string{"id": tc.id.URN()})
		w := httptest.NewRecorder()

		h.GetByID(w, r)

		resp := w.Result()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("cannot read response")
		}

		resp.Body.Close()

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if reflect.DeepEqual(body, tc.resp) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), string(tc.resp))
		}
	}
}

func TestHandler_Update(t *testing.T) {
	fields := []byte(`{"Model":"BMW","YearOfManufacture":2022,"Brand":"BMW",
					"FuelType":"Petrol"}`)

	cases := []struct {
		desc       string
		id         uuid.UUID
		body       io.Reader
		statusCode int
	}{
		{"entity updated successfully", uuid.UUID{}, bytes.NewReader(fields), http.StatusOK},
		{"entity does not exist", uuid.Nil, bytes.NewReader([]byte("")), http.StatusNotFound},
		{"unable to read body", uuid.UUID{}, mockReader{}, http.StatusBadRequest},
		{"unmarshal error", uuid.UUID{}, bytes.NewReader([]byte("invalid Body")), http.StatusBadRequest},
		{"database connectivity error", uuid.UUID{}, bytes.NewReader([]byte("")), http.StatusInternalServerError},
	}

	h := New(mockService{})

	for i, tc := range cases {
		req := httptest.NewRequest(http.MethodPut, "http://localhost:8000/car/{id}", tc.body)
		r := mux.SetURLVars(req, map[string]string{"id": tc.id.URN()})
		w := httptest.NewRecorder()

		h.Update(w, r)

		resp := w.Result()

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}
	}
}

func TestHandler_Delete(t *testing.T) {
	cases := []struct {
		desc       string
		id         uuid.UUID
		statusCode int
	}{
		{"delete successful", uuid.UUID{}, http.StatusNoContent},
		{"entity does not exist", uuid.UUID{}, http.StatusNotFound},
		{"invalid id", uuid.Nil, http.StatusBadRequest},
		{"database connectivity error", uuid.UUID{}, http.StatusInternalServerError},
	}

	h := New(mockService{})

	for i, tc := range cases {
		req := httptest.NewRequest(http.MethodPut, "http://localhost:8000/car/{id}", http.NoBody)
		r := mux.SetURLVars(req, map[string]string{"id": tc.id.URN()})
		w := httptest.NewRecorder()

		h.Delete(w, r)

		resp := w.Result()

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}
	}
}
