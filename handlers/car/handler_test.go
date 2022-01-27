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
	"github.com/amehrotra/car-dealership/types"
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
	withEngine := []models.Car{
		{
			ID:                uuid.Nil,
			Model:             "X",
			YearOfManufacture: 2020,
			Brand:             "BMW",
			FuelType:          types.Petrol,
			Engine: models.Engine{
				ID:           uuid.Nil,
				Displacement: 100,
				NCylinder:    2,
				Range:        0,
			},
		},
	}

	withoutEngine := []models.Car{
		{
			ID:                uuid.Nil,
			Model:             "X",
			YearOfManufacture: 2020,
			Brand:             "BMW",
			FuelType:          types.Petrol,
			Engine:            models.Engine{},
		},
	}

	cases := []struct {
		desc       string
		filter     filters.Car
		output     []models.Car
		statusCode int
	}{
		{"get all cars  with engine", filters.Car{Brand: "BMW", Engine: true}, withEngine, http.StatusOK},
		{"get all cars without engine", filters.Car{Brand: "BMW", Engine: false}, withoutEngine, http.StatusOK},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodGet, http.NoBody, nil, map[string][]string{"brand": {tc.filter.Brand}, "engine": {strconv.FormatBool(tc.filter.Engine)}})

		h.GetAll(w, r)

		resp := w.Result()

		body, err := getResponseBody(resp)
		if err != nil {
			t.Errorf("error in reading body")
		}

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if reflect.DeepEqual(body, tc.output) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), (tc.output))
		}
	}
}

func TestHandler_GetByID(t *testing.T) {
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
		id         uuid.UUID
		output     models.Car
		statusCode int
	}{
		{"request successful", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), car, http.StatusOK},
		{"invalid id", uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), models.Car{}, http.StatusBadRequest},
		{"database error", uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"), models.Car{}, http.StatusInternalServerError},
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

		if reflect.DeepEqual(body, tc.output) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), tc.output)
		}
	}
}

func TestHandler_Update(t *testing.T) {
	fields := []byte(`{"id":"8f443772-132b-4ae5-9f8f-9960649b3fb4","model":"x","yearOfManufacture":2020,"brand":"BMW"
		,"fuelType":0,"engine":{"displacement":200,"noOfCylinder":2,"range":0}}`)

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
		id         uuid.UUID
		body       []byte
		resp       models.Car
		statusCode int
	}{
		{"entity updated successfully", uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"), fields, car, http.StatusOK},
		{"entity not found", uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"), []byte(""), models.Car{}, http.StatusNotFound},
		//{"unable to read body", uuid.UUID{}, mockReader{}, []byte(""), http.StatusBadRequest},
		{"database error", uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), []byte(""), models.Car{}, http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodPut, bytes.NewReader(tc.body), map[string]string{"id": tc.id.URN()}, nil)

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
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), tc.resp)
		}
	}
}

func TestHandler_Delete(t *testing.T) {
	cases := []struct {
		desc       string
		id         uuid.UUID
		statusCode int
	}{
		{"delete successful", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), http.StatusNoContent},
		//{"entity does not exist", uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"), http.StatusNotFound},
		{"invalid id", uuid.Nil, http.StatusBadRequest},
		{"database error", uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodDelete, http.NoBody, map[string]string{"id": tc.id.URN()}, nil)

		h.Delete(w, r)

		resp := w.Result()

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}
	}
}
