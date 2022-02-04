package car

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/services"
	"github.com/amehrotra/car-dealership/types"
)

func initializeTest(t *testing.T, method string, body io.Reader, pParam map[string]string, qParam url.Values) (handler, *services.MockCar,
	*http.Request, *httptest.ResponseRecorder) {
	ctrl := gomock.NewController(t)

	mockService := services.NewMockCar(ctrl)
	h := New(mockService)

	req := httptest.NewRequest(method, "http://cars", body)
	r := mux.SetURLVars(req, pParam)
	r.URL.RawQuery = qParam.Encode()

	w := httptest.NewRecorder()

	return h, mockService, r, w
}

func getResponseBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil
}

func getOutput(t *testing.T, respBody []byte) *models.Car {
	var output *models.Car

	if len(respBody) != 0 {
		output = &models.Car{}
		if err := json.Unmarshal(respBody, output); err != nil {
			t.Error(err)
		}
	}

	return output
}

func TestHandler_Create(t *testing.T) {
	var (
		body = []byte(`{"id":"8f443772-132b-4ae5-9f8f-9960649b3fb4","model":"x","yearOfManufacture":2020,"brand":"BMW","fuelType":"petrol",
		"engine":{"displacement":100,"noOfCylinder":2,"range":0}}`)

		car = models.Car{
			ID:              uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
			Model:           "x",
			ManufactureYear: 2020,
			Brand:           "BMW",
			FuelType:        types.Petrol,
			Engine:          models.Engine{Displacement: 100, NCylinder: 2},
		}
	)

	cases := []struct {
		desc       string
		body       io.Reader
		mockOutput *models.Car
		mockErr    error
		resp       *models.Car
		statusCode int
	}{
		{"success case", bytes.NewReader(body), &car, nil, &car, http.StatusCreated},
		{"entity already exists", bytes.NewReader(body), nil, errors.EntityAlreadyExists{}, nil, http.StatusOK},
		{"invalid body", mockReader{}, nil, errors.InvalidParam{}, nil, http.StatusBadRequest},
	}

	for i, tc := range cases {
		h, mockService, r, w := initializeTest(t, http.MethodPost, tc.body, nil, nil)

		mockService.EXPECT().Create(&car).Return(tc.mockOutput, tc.mockErr)

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
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), tc.resp)
		}
	}
}

func Test_CreateInvalidBody(t *testing.T) {
	h, _, r, w := initializeTest(t, http.MethodPost, mockReader{}, nil, nil)
	h.Create(w, r)

	resp := w.Result()

	body, err := getResponseBody(resp)
	if err != nil {
		t.Errorf("error in reading body")
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", resp.StatusCode, http.StatusBadRequest)
	}

	if reflect.DeepEqual(body, resp.Body) {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", string(body), resp.Body)
	}
}

func TestHandler_GetAll(t *testing.T) {
	withEngine := []models.Car{
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

	withoutEngine := []models.Car{
		{
			ID:              uuid.Nil,
			Model:           "X",
			ManufactureYear: 2020,
			Brand:           "BMW",
			FuelType:        types.Petrol,
			Engine:          models.Engine{},
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
		h, mockService, r, w := initializeTest(t, http.MethodGet, http.NoBody, nil,
			map[string][]string{"brand": {tc.filter.Brand}, "engine": {strconv.FormatBool(tc.filter.Engine)}})

		mockService.EXPECT().GetAll(tc.filter).Return(tc.output, nil)

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
		ID:              uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Model:           "X",
		ManufactureYear: 2020,
		Brand:           "BMW",
		FuelType:        types.Petrol,
		Engine: models.Engine{
			Displacement: 100,
			NCylinder:    2,
		},
	}

	cases := []struct {
		desc       string
		id         uuid.UUID
		output     models.Car
		expErr     error
		statusCode int
	}{
		{"request successful", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), car, nil,
			http.StatusOK},
		{"entity does not exist", uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"), models.Car{},
			errors.EntityAlreadyExists{}, http.StatusOK},
	}

	for i, tc := range cases {
		h, mockService, r, w := initializeTest(t, http.MethodGet, http.NoBody, map[string]string{"id": tc.id.URN()}, nil)

		mockService.EXPECT().GetByID(tc.id).Return(&tc.output, tc.expErr)

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
	body := []byte(`{"id":"8f443772-132b-4ae5-9f8f-9960649b3fb4","model":"X","yearOfManufacture":2020,"brand":"BMW"
		,"fuelType":"petrol","engine":{"displacement":200,"noOfCylinder":2,"range":0}}`)

	car := models.Car{
		ID:              uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
		Model:           "X",
		ManufactureYear: 2020,
		Brand:           "BMW",
		FuelType:        types.Petrol,
		Engine: models.Engine{
			Displacement: 200,
			NCylinder:    2,
		},
	}

	cases := []struct {
		desc       string
		mockErr    error
		resp       *models.Car
		statusCode int
	}{
		{"entity updated successfully", nil, &car, http.StatusOK},
		{"entity not found", errors.EntityNotFound{}, nil, http.StatusNotFound},
	}

	for i, tc := range cases {
		h, mockService, r, w := initializeTest(t, http.MethodPut, bytes.NewReader(body),
			map[string]string{"id": "8f443772-132b-4ae5-9f8f-9960649b3fb4"}, nil)

		mockService.EXPECT().Update(&car).Return(tc.resp, tc.mockErr)

		h.Update(w, r)

		resp := w.Result()

		respBody, err := getResponseBody(resp)
		if err != nil {
			t.Errorf("error in reading body")
		}

		output := getOutput(t, respBody)

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if !reflect.DeepEqual(output, tc.resp) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, output, tc.resp)
		}
	}
}

func TestHandler_Delete(t *testing.T) {
	cases := []struct {
		desc       string
		id         uuid.UUID
		mockErr    error
		statusCode int
	}{
		{"delete successful", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), nil, http.StatusNoContent},
		{"entity does not exist", uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"), errors.EntityNotFound{}, http.StatusNotFound},
		{"invalid id", uuid.Nil, errors.InvalidParam{}, http.StatusBadRequest},
	}

	for i, tc := range cases {
		h, mockService, r, w := initializeTest(t, http.MethodDelete, http.NoBody, map[string]string{"id": tc.id.URN()}, nil)

		mockService.EXPECT().Delete(tc.id).Return(tc.mockErr)

		h.Delete(w, r)

		resp := w.Result()

		if err := resp.Body.Close(); err != nil {
			log.Println("error in closing response body")
		}

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}
	}
}

type mockReader struct{}

func (m mockReader) Read(p []byte) (n int, err error) {
	return 0, errors.BindError{}
}

func Test_getCar(t *testing.T) {
	cases := []struct {
		desc   string
		body   io.Reader
		output *models.Car
		err    error
	}{
		{"bind error", mockReader{}, nil, errors.InvalidParam{Param: []string{"body"}}},
		{"unmarshal error", bytes.NewReader([]byte("invalid body")), nil, errors.InvalidParam{Param: []string{"body"}}},
	}

	for i, tc := range cases {
		_, _, r, _ := initializeTest(t, "", tc.body, nil, nil)

		output, err := getCar(r)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("\n[TEST %v] Failed. Desc %v: \nGot %v\nExpected %v", i, tc.desc, err, tc.err)
		}

		if output != tc.output {
			t.Errorf("\n[TEST %v] Failed. Desc %v: \nGot %v\nExpected %v", i, tc.desc, output, tc.output)
		}
	}
}

func Test_getID(t *testing.T) {
	cases := []struct {
		desc   string
		id     string
		output uuid.UUID
		err    error
	}{
		{"id parsed success", "8f443772-132b-4ae5-9f8f-9960649b3fb4", uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"), nil},
		{"empty string", "", uuid.Nil, errors.MissingParam{Param: "id"}},
		{"invalid id", "12223", uuid.Nil, errors.InvalidParam{Param: []string{"id"}}},
	}

	for i, tc := range cases {
		_, _, r, _ := initializeTest(t, "", nil, map[string]string{"id": tc.id}, nil)

		output, err := getID(r)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("\n[TEST %v] Failed. Desc %v: \nGot %v\nExpected %v", i, tc.desc, err, tc.err)
		}

		if output != tc.output {
			t.Errorf("\n[TEST %v] Failed. Desc %v: \nGot %v\nExpected %v", i, tc.desc, output, tc.output)
		}
	}
}
