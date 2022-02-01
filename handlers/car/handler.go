package car

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/services"
)

type handler struct {
	service services.Car
}

func New(service services.Car) handler {
	return handler{service: service}
}

// Create takes the clients request to create entity in database
func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	var car models.Car

	car, err := getCar(r)
	if err != nil {
		setStatusCode(w, err, r.Method, car)

		return
	}

	car, err = h.service.Create(&car)
	setStatusCode(w, err, r.Method, car)
}

// GetAll writes all the cars from the database based on the query parameter
func (h handler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var engine bool
	if query.Get("engine") == "true" {
		engine = true
	}

	filter := filters.Car{
		Brand:  query.Get("brand"),
		Engine: engine,
	}

	resp, err := h.service.GetAll(filter)
	setStatusCode(w, err, r.Method, resp)
}

// GetByID writes the car based on ID of the car
func (h handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		setStatusCode(w, errors.InvalidParam{}, r.Method, nil)

		return
	}

	var car models.Car

	car, err = h.service.GetByID(id)

	setStatusCode(w, err, r.Method, car)
}

// Update writes the updated car entity in the database
func (h handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		setStatusCode(w, errors.InvalidParam{}, r.Method, nil)

		return
	}

	var car models.Car
	car.ID = id

	car, err = getCar(r)
	if err != nil {
		setStatusCode(w, err, r.Method, car)

		return
	}

	car, err = h.service.Update(car)
	setStatusCode(w, err, r.Method, car)
}

// Delete removes the car from database based on ID
func (h handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		setStatusCode(w, errors.InvalidParam{}, r.Method, nil)

		return
	}

	err = h.service.Delete(id)
	setStatusCode(w, err, r.Method, nil)
}

// getCar return the car by reading and unmarshal the body
func getCar(r *http.Request) (models.Car, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return models.Car{}, err
	}

	var car models.Car

	err = json.Unmarshal(body, &car)
	if err != nil {
		return models.Car{}, err
	}

	return car, nil
}

// setStatusCode writes the status code based on the error type
func setStatusCode(w http.ResponseWriter, err error, method string, data interface{}) {
	switch err.(type) {
	case errors.EntityAlreadyExists:
		w.WriteHeader(http.StatusOK)
	case errors.MissingParam, errors.InvalidParam:
		w.WriteHeader(http.StatusBadRequest)
	case nil:
		writeSuccessResponse(method, w, data)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// writeSuccessResponse based on the method type it calls function writeResponseBody
func writeSuccessResponse(method string, w http.ResponseWriter, data interface{}) {
	switch method {
	case http.MethodPost:
		if data != nil {
			writeResponseBody(w, http.StatusCreated, data)

			return
		}

		writeResponseBody(w, http.StatusOK, data)
	case http.MethodGet:
		writeResponseBody(w, http.StatusOK, data)
	case http.MethodPut:
		writeResponseBody(w, http.StatusOK, data)
	case http.MethodDelete:
		writeResponseBody(w, http.StatusNoContent, data)
	}
}

// writeResponseBody marshals the data and writes the body which is sent to client
func writeResponseBody(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(statusCode)

	_, err = w.Write(resp)
	if err != nil {
		log.Println("error in writing response")

		return
	}
}

// getID reads the id from path parameter of url
func getID(r *http.Request) (uuid.UUID, error) {
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := uuid.Parse(idParam)
	if err != nil {
		return uuid.Nil, errors.InvalidParam{}
	}

	return id, nil
}
