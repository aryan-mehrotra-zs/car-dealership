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

	"github.com/amehrotra/car-dealership/filters"
)

func initializeTest(method string, body io.Reader, queryParams map[string]string) (handler, *http.Request, *httptest.ResponseRecorder) {
	h := New(mockService{})

	req := httptest.NewRequest(method, "http://car", body)
	r := mux.SetURLVars(req, queryParams)
	w := httptest.NewRecorder()

	return h, r, w
}

// should we do error handling here?
func readBody(w *httptest.ResponseRecorder, t *testing.T) ([]byte, *http.Response) {
	resp := w.Result()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("cannot read response")
	}

	err = resp.Body.Close()
	if err != nil {
		t.Errorf("Error in closing body")
	}

	return body, resp
}

func TestHandler_Create(t *testing.T) {
	res := []byte(`("ID":"123e4567-e89b-12d3-a456-426614174000","Model":"BMW","YearOfManufacture":2022,"Brand":"BMW",
					"FuelType":"Petrol",{"Displacement":20,"NCylinder":2,"Range":0}}`)

	cases := []struct {
		desc       string
		body       io.Reader
		car        []byte
		statusCode int
	}{
		{"created successfully", bytes.NewReader(res), res, http.StatusCreated},
		{"entity already exists", bytes.NewReader(res), res, http.StatusOK},
		{"unmarshal error", bytes.NewReader([]byte(`Invalid Body`)), []byte(""), http.StatusBadRequest},
		{"missing parameter", bytes.NewReader([]byte(`{"Model":"BMW","YearOfManufacture":2022,"Brand":"BMW"}`)), []byte(""), http.StatusBadRequest},
		{"invalid parameter", bytes.NewReader([]byte(`{"Model":BMW,"YearOfManufacture":2022,"Brand":"BMW","FuelType":"Petrol",{20,2,0}}`)), []byte(""), http.StatusBadRequest},
		{"unable to read body", mockReader{}, []byte(""), http.StatusBadRequest},
		{"database connectivity error", bytes.NewReader(res), res, http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodPost, tc.body, nil)

		h.Create(w, r)

		body, resp := readBody(w, t)

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

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodPost, http.NoBody, map[string]string{"brand": tc.car.Brand, "engine": strconv.FormatBool(tc.car.Engine)})

		h.GetAll(w, r)

		body, resp := readBody(w, t)

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if reflect.DeepEqual(body, tc.resp) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), string(tc.resp))
		}
	}
}

func TestHandler_GetByID(t *testing.T) {
	res := []byte(`{"ID":"123e4567-e89b-12d3-a456-426614174000","Model":"BMW","YearOfManufacture":2022,"Brand":"BMW",
					"FuelType":"Petrol"}`)

	cases := []struct {
		desc       string
		id         uuid.UUID
		resp       []byte
		statusCode int
	}{
		{"request successful", uuid.UUID{}, res, http.StatusOK},
		{"entity not found", uuid.UUID{}, []byte(""), http.StatusBadRequest},
		{"invalid id", uuid.Nil, []byte(""), http.StatusBadRequest},
		{"database connectivity error", uuid.UUID{}, []byte(""), http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodGet, http.NoBody, map[string]string{"id": tc.id.URN()})

		h.GetByID(w, r)

		body, resp := readBody(w, t)

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
		resp       []byte
		statusCode int
	}{
		{"entity updated successfully", uuid.UUID{}, bytes.NewReader(fields), fields, http.StatusOK},
		{"entity does not exist", uuid.Nil, bytes.NewReader([]byte("")), []byte(""), http.StatusNotFound},
		{"unable to read body", uuid.UUID{}, mockReader{}, []byte(""), http.StatusBadRequest},
		{"unmarshal error", uuid.UUID{}, bytes.NewReader([]byte("invalid Body")), []byte(""), http.StatusBadRequest},
		{"database connectivity error", uuid.UUID{}, bytes.NewReader([]byte("")), []byte(""), http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodPut, tc.body, map[string]string{"id": tc.id.URN()})

		h.Update(w, r)

		body, resp := readBody(w, t)

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if reflect.DeepEqual(body, tc.resp) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), string(tc.resp))
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

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodPut, http.NoBody, map[string]string{"id": tc.id.URN()})

		h.Delete(w, r)

		resp := w.Result()

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}
	}
}
