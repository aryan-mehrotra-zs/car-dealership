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

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	var car models.Car

	car, err := getCar(r)
	if err != nil {
		setStatusCode(w, err, r.Method, car)

		return
	}

	car, err = h.service.Create(car)
	setStatusCode(w, err, r.Method, car)
}

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

func (h handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	if id == uuid.Nil {
		setStatusCode(w, errors.InvalidParam{}, r.Method, nil)

		return
	}

	var car models.Car

	car, err := h.service.GetByID(id)

	setStatusCode(w, err, r.Method, car)
}

func (h handler) Update(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	if id == uuid.Nil {
		setStatusCode(w, errors.InvalidParam{}, r.Method, nil)

		return
	}

	var car models.Car
	car, err := getCar(r)
	if err != nil {
		setStatusCode(w, err, r.Method, car)

		return
	}

	car, err = h.service.Update(car)
	setStatusCode(w, err, r.Method, car)
}

func (h handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	if id == uuid.Nil {
		setStatusCode(w, errors.InvalidParam{}, r.Method, nil)

		return
	}

	err := h.service.Delete(id)
	setStatusCode(w, err, r.Method, nil)
}

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

func getID(r *http.Request) uuid.UUID {
	param := mux.Vars(r)
	idParam := param["id"]

	// How to handle panics
	// MustParse is like Parse but panics if the string cannot be parsed. It simplifies safe initialization of global variables holding compiled UUIDs.
	id := uuid.MustParse(idParam)

	return id
}
