package car

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
)

func initializeTest(method string, body io.Reader, pathParams map[string]string, queryParams url.Values) (handler, *http.Request, *httptest.ResponseRecorder) {
	h := New(mockService{})

	req := httptest.NewRequest(method, "http://car", body)
	r := mux.SetURLVars(req, pathParams)
	r.URL.RawQuery = queryParams.Encode()

	w := httptest.NewRecorder()

	return h, r, w
}

func getResponseBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil
}

func TestHandler_Create(t *testing.T) {
	res := []byte(`{"id":"8f443772-132b-4ae5-9f8f-9960649b3fb4","model":"x","yearOfManufacture":2020,"brand":"BMW","fuelType":0,"engine":{"displacement":200,"noOfCylinder":2,"range":0}}`)
	res2 := []byte(`{"id":"8f443772-132b-4ae5-9f8f-9960649b3fb4","model":"y","yearOfManufacture":2020,"brand":"BMW","fuelType":0,"engine":{"displacement":200,"noOfCylinder":2,"range":0}}`)
	res3 := []byte(`{"id":"8f443772-132b-4ae5-9f8f-9960649b3fb4","model":"z","yearOfManufacture":2020,"brand":"BMW","fuelType":0,"engine":{"displacement":200,"noOfCylinder":2,"range":0}}`)

	car := models.Car{
		ID:                uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Model:             "X",
		YearOfManufacture: 2020,
		Brand:             "BMW",
		FuelType:          0,
		Engine: models.Engine{
			Displacement: 100,
			NCylinder:    2,
			Range:        0,
		},
	}

	cases := []struct {
		desc       string
		body       []byte
		car        models.Car
		statusCode int
	}{
		{"created successfully", res, car, http.StatusCreated},
		{"entity already exists", res2, models.Car{}, http.StatusOK},
		//{"unmarshal error", []byte(`invalid body`), models.Car{}, http.StatusBadRequest},
		{"database error", res3, models.Car{}, http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodPost, bytes.NewReader(tc.body), nil, nil)

		h.Create(w, r)

		resp := w.Result()

		body, err := getResponseBody(resp)
		if err != nil {
			t.Errorf("error in reading body")
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if reflect.DeepEqual(body, resp.Body) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), tc.car)
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
		filter     filters.Car
		resp       []byte
		statusCode int
	}{
		{"get all cars of a brand with engine", filters.Car{Brand: "BMW", Engine: true}, withoutEngine, http.StatusOK},
		{"get all cars without engine", filters.Car{Brand: "BMW", Engine: false}, withEngine, http.StatusOK},
		{"get all cars from all brands without engine", filters.Car{Brand: "", Engine: true}, allWithoutEngine, http.StatusOK},
		{"get all cars from all brands with engine", filters.Car{Brand: "", Engine: true}, allWithEngine, http.StatusOK},
		{"invalid brand name", filters.Car{Brand: "xyz", Engine: true}, []byte(""), http.StatusBadRequest},
		{"database error", filters.Car{Brand: "", Engine: true}, []byte(""), http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodPost, http.NoBody, nil, map[string][]string{"brand": {tc.filter.Brand}, "engine": {strconv.FormatBool(tc.filter.Engine)}})

		h.GetAll(w, r)

		resp := w.Result()

		body, err := getResponseBody(resp)
		if err != nil {
			t.Errorf("error in reading body")
		}

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
		{"database error", uuid.UUID{}, []byte(""), http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodGet, http.NoBody, map[string]string{"id": tc.id.URN()}, nil)

		h.GetByID(w, r)

		resp := w.Result()

		body, err := getResponseBody(resp)
		if err != nil {
			t.Errorf("error in reading body")
		}

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
		{"database error", uuid.UUID{}, bytes.NewReader([]byte("")), []byte(""), http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodPut, tc.body, map[string]string{"id": tc.id.URN()}, nil)

		h.Update(w, r)

		resp := w.Result()

		body, err := getResponseBody(resp)
		if err != nil {
			t.Errorf("error in reading body")
		}

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
		{"database error", uuid.UUID{}, http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodPut, http.NoBody, map[string]string{"id": tc.id.URN()}, nil)

		h.Delete(w, r)

		resp := w.Result()

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}
	}
}
