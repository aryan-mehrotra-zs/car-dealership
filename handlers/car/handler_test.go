package car

import (
	"bytes"
	"encoding/json"
	goError "errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
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

// todo use NewRandom instead of MustParse
// todo remove goerror
// todo error should not be supressed

var car = models.Car{
	ID:              uuid.MustParse("8f443772-132b-4ae5-9f8f-9960649b3fb4"),
	Model:           "x",
	ManufactureYear: 2020,
	Brand:           "BMW",
	FuelType:        types.Petrol,
	Engine:          models.Engine{Displacement: 100, NCylinder: 2},
}

func TestHandler_Create(t *testing.T) {
	var (
		body = []byte(`{"id":"8f443772-132b-4ae5-9f8f-9960649b3fb4","model":"x","yearOfManufacture":2020,"brand":"BMW","fuelType":"petrol",
		"engine":{"displacement":100,"noOfCylinder":2,"range":0}}`)
	)

	cases := []struct {
		desc       string
		mockOutput *models.Car
		mockErr    error
		resp       *models.Car
		statusCode int
	}{
		{"success case", &car, nil, &car, http.StatusCreated},
		{"entity already exists", nil, errors.EntityAlreadyExists{}, nil, http.StatusOK},
		{"internal server error", nil, errors.DB{}, nil, http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, mockService, r, w := initializeTest(t, http.MethodPost, bytes.NewReader(body), nil, nil)

		mockService.EXPECT().Create(&car).Return(tc.mockOutput, tc.mockErr)

		h.Create(w, r)

		resp := w.Result()

		body, err := getResponseBody(resp)
		if err != nil {
			t.Errorf("error in reading body : %v", err)
		}

		output := getOutput(t, body)

		if resp.StatusCode != tc.statusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if !reflect.DeepEqual(output, tc.resp) {
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

	output := getOutput(t, body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", resp.StatusCode, http.StatusBadRequest)
	}

	if output != nil {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", output, resp.Body)
	}
}

func TestHandler_GetAll(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("unable to create uuid : %v", err)
	}

	withEngine := []models.Car{
		{
			ID:              id,
			Model:           "X",
			ManufactureYear: 2020,
			Brand:           "BMW",
			FuelType:        types.Petrol,
			Engine: models.Engine{
				ID:           id,
				Displacement: 100,
				NCylinder:    2,
			},
		},
	}

	withoutEngine := []models.Car{
		{
			ID:              id,
			Model:           "X",
			ManufactureYear: 2020,
			Brand:           "BMW",
			FuelType:        types.Petrol,
		},
	}

	cases := []struct {
		desc       string
		filter     filters.Car
		mockOutput []models.Car
		mockErr    error
		statusCode int
	}{
		{"get all cars  with engine", filters.Car{Brand: "BMW", Engine: true}, withEngine, nil, http.StatusOK},
		{"get all cars without engine", filters.Car{Brand: "BMW", Engine: false}, withoutEngine, nil, http.StatusOK},
		{"invalid parameter", filters.Car{Brand: "xyz"}, nil, errors.InvalidParam{Param: []string{"brand"}}, http.StatusBadRequest},
		{"internal server error", filters.Car{Brand: "BMW"}, nil, errors.DB{}, http.StatusInternalServerError},
	}

	for i, tc := range cases {
		params := map[string][]string{"brand": {tc.filter.Brand}, "engine": {strconv.FormatBool(tc.filter.Engine)}}

		h, mockService, r, w := initializeTest(t, http.MethodGet, http.NoBody, nil, params)

		mockService.EXPECT().GetAll(tc.filter).Return(tc.mockOutput, tc.mockErr)

		h.GetAll(w, r)

		resp := w.Result()

		body, err := getResponseBody(resp)
		if err != nil {
			t.Errorf("error in reading body")
		}

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if !reflect.DeepEqual(body, tc.mockOutput) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), (tc.mockOutput))
		}
	}
}

func TestHandler_GetByID(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("unable to create uuid : %v", err)
	}

	car.ID = id

	cases := []struct {
		desc       string
		mockOutput *models.Car
		mockErr    error
		statusCode int
	}{
		{"request successful", &car, nil, http.StatusOK},
		{"entity does not exist", nil, errors.EntityNotFound{}, http.StatusNotFound},
		{"internal server error", nil, errors.DB{}, http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, mockService, r, w := initializeTest(t, http.MethodGet, http.NoBody, map[string]string{"id": id.URN()}, nil)

		mockService.EXPECT().GetByID(id).Return(tc.mockOutput, tc.mockErr)

		h.GetByID(w, r)

		resp := w.Result()

		body, err := getResponseBody(resp)
		if err != nil {
			t.Errorf("error in reading body")
		}

		if tc.statusCode != resp.StatusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if reflect.DeepEqual(body, tc.mockOutput) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), tc.mockOutput)
		}
	}
}

func Test_GetByIDInvalidID(t *testing.T) {
	h, _, r, w := initializeTest(t, http.MethodGet, http.NoBody, nil, nil)
	h.GetByID(w, r)

	resp := w.Result()

	body, err := getResponseBody(resp)
	if err != nil {
		t.Errorf("error in reading body : %v", err)
	}

	output := getOutput(t, body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", resp.StatusCode, http.StatusBadRequest)
	}

	if output != nil {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", string(body), resp.Body)
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
		{"internal server error", errors.DB{}, nil, http.StatusInternalServerError},
	}

	for i, tc := range cases {
		param := map[string]string{"id": "8f443772-132b-4ae5-9f8f-9960649b3fb4"}

		h, mockService, r, w := initializeTest(t, http.MethodPut, bytes.NewReader(body), param, nil)

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

func Test_UpdateInvalidID(t *testing.T) {
	param := map[string]string{"id": "8f443772-132b-4ae5-9f8f-9960649b3fb4"}
	h, _, r, w := initializeTest(t, http.MethodPut, nil, param, nil)
	h.Update(w, r)

	resp := w.Result()

	body, err := getResponseBody(resp)
	if err != nil {
		t.Errorf("error in reading body")
	}

	output := getOutput(t, body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", resp.StatusCode, http.StatusBadRequest)
	}

	if output != nil {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", string(body), resp.Body)
	}
}

func Test_UpdateInvalidBody(t *testing.T) {
	h, _, r, w := initializeTest(t, http.MethodPut, nil, nil, nil)
	h.Update(w, r)

	resp := w.Result()

	body, err := getResponseBody(resp)
	if err != nil {
		t.Errorf("error in reading body")
	}

	output := getOutput(t, body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", resp.StatusCode, http.StatusBadRequest)
	}

	if output != nil {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", string(body), resp.Body)
	}
}

func TestHandler_Delete(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("unable to create uuid : %v", err)
	}

	cases := []struct {
		desc       string
		mockErr    error
		statusCode int
	}{
		{"delete successful", nil, http.StatusNoContent},
		{"entity does not exist", errors.EntityNotFound{}, http.StatusNotFound},
		{"internal server error", errors.DB{}, http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, mockService, r, w := initializeTest(t, http.MethodDelete, http.NoBody, map[string]string{"id": id.URN()}, nil)

		mockService.EXPECT().Delete(id).Return(tc.mockErr)

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

func Test_DeleteInvalidID(t *testing.T) {
	h, _, r, w := initializeTest(t, http.MethodDelete, http.NoBody, nil, nil)
	h.Delete(w, r)

	resp := w.Result()

	body, err := getResponseBody(resp)
	if err != nil {
		t.Errorf("error in reading body")
	}

	output := getOutput(t, body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", resp.StatusCode, http.StatusBadRequest)
	}

	if output != nil {
		t.Errorf("\n[TEST] Failed. Desc : invalid body\nGot %v\nExpected %v", string(body), resp.Body)
	}
}

func Test_getCar(t *testing.T) {
	var (
		body = bytes.NewReader([]byte(`{"id":"8f443772-132b-4ae5-9f8f-9960649b3fb4","model":"x","yearOfManufacture":2020,"brand":"BMW","fuelType":"petrol",
		"engine":{"displacement":100,"noOfCylinder":2,"range":0}}`))

		invalidBody = bytes.NewReader([]byte("invalid body"))
	)

	cases := []struct {
		desc   string
		body   io.Reader
		output *models.Car
		err    error
	}{
		{"success", body, &car, nil},
		{"bind error", mockReader{}, nil, errors.InvalidParam{Param: []string{"body"}}},
		{"unmarshal error", invalidBody, nil, errors.InvalidParam{Param: []string{"body"}}},
	}

	for i, tc := range cases {
		_, _, r, _ := initializeTest(t, "", tc.body, nil, nil)

		output, err := getCar(r)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("\n[TEST %v] Failed. Desc %v: \nGot %v\nExpected %v", i, tc.desc, err, tc.err)
		}

		if !reflect.DeepEqual(output, tc.output) {
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

func Test_writeResponseBodyMarshalError(t *testing.T) {
	data := complex(1, 1)

	w := httptest.NewRecorder()
	expectedStatusCode := http.StatusInternalServerError

	writeResponseBody(w, http.StatusOK, data)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		t.Errorf("\n[TEST] Failed. Desc : Marshal Error \nGot %v\nExpected %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

func Test_writeResponseBodyWriteError(t *testing.T) {
	data := []byte(`{"id":"8f443772-132b-4ae5-9f8f-9960649b3fb4","model":"x","yearOfManufacture":2020,"brand":"BMW","fuelType":"petrol",
		"engine":{"displacement":100,"noOfCylinder":2,"range":0}}`)

	w := mockResponseWriter{}

	var b bytes.Buffer
	log.SetOutput(&b)

	writeResponseBody(w, http.StatusOK, data)

	if !strings.Contains(b.String(), "error in writing response") {
		t.Errorf("\n[TEST] Failed. Desc : Write Error \nGot %v\nExpected 'error in writing response' in logs", b.String())
	}
}

type mockReader struct{}

func (m mockReader) Read(p []byte) (n int, err error) {
	return 0, goError.New("bind error")
}

type mockResponseWriter struct{}

func (m mockResponseWriter) Header() http.Header {
	header := make(map[string][]string)

	return header
}

func (m mockResponseWriter) Write([]byte) (int, error) {
	return 0, goError.New("error")
}

func (m mockResponseWriter) WriteHeader(statusCode int) {

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
