package car

import (
	"net/http"

	"github.com/amehrotra/car-dealership/services"
)

type handler struct {
	service services.Car
}

func New(car services.Car) handler {
	return handler{service: car}
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h handler) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (h handler) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (h handler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h handler) Delete(w http.ResponseWriter, r *http.Request) {

}
