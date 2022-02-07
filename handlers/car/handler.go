package car

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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

// nolint:revive // handler should not be exported
func New(service services.Car) handler {
	return handler{service: service}
}

// Create takes the clients request to create entity in database
func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	var car *models.Car

	car, err := getCar(r)
	if err != nil {
		setStatusCode(w, r.Method, car, err)

		return
	}

	car, err = h.service.Create(car)
	setStatusCode(w, r.Method, car, err)
}

// GetAll writes all the cars from the database based on the query parameter
func (h handler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var hasEngine bool

	param := strings.TrimSpace(query.Get("engine"))
	if strings.EqualFold(param, "true") {
		hasEngine = true
	}

	brand := strings.TrimSpace(query.Get("brand"))

	filter := filters.Car{Brand: brand, Engine: hasEngine}

	resp, err := h.service.GetAll(filter)
	setStatusCode(w, r.Method, resp, err)
}

// GetByID writes the response based on ID of the resp
func (h handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		setStatusCode(w, r.Method, nil, err)

		return
	}

	car, err := h.service.GetByID(id)
	log.Println(err)
	setStatusCode(w, r.Method, car, err)
}

// Update writes the updated resp entity in the database
func (h handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		setStatusCode(w, r.Method, nil, err)

		return
	}

	car, err := getCar(r)
	if err != nil {
		setStatusCode(w, r.Method, car, err)

		return
	}

	car.ID = id

	car, err = h.service.Update(car)
	setStatusCode(w, r.Method, car, err)
}

// Delete removes the resp from database based on ID
func (h handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		setStatusCode(w, r.Method, nil, err)

		return
	}

	err = h.service.Delete(id)
	setStatusCode(w, r.Method, nil, err)
}

// getID reads the id from path parameter of url
func getID(r *http.Request) (uuid.UUID, error) {
	param := mux.Vars(r)
	idParam := strings.TrimSpace(param["id"])

	if idParam == "" {
		return uuid.Nil, errors.MissingParam{Param: "id"}
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		return uuid.Nil, errors.InvalidParam{Param: []string{"id"}}
	}

	return id, nil
}

// getCar reads request body and returns car
func getCar(r *http.Request) (*models.Car, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	var car models.Car

	err = json.Unmarshal(body, &car)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	return &car, nil
}
