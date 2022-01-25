package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/amehrotra/car-dealership/drivers"
	handlers "github.com/amehrotra/car-dealership/handlers/car"
	services "github.com/amehrotra/car-dealership/services/car"
	"github.com/amehrotra/car-dealership/stores/car"
	"github.com/amehrotra/car-dealership/stores/engine"
)

func main() {
	db, err := drivers.ConnectToSQL()
	if err != nil {
		return
	}

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	carStore := car.New(db)
	engineStore := engine.New(db)
	service := services.New(engineStore, carStore)
	handler := handlers.New(service)

	r := mux.NewRouter()
	r.HandleFunc("/car", handler.Create).Methods(http.MethodPost)
	r.HandleFunc("/car", handler.GetAll).Queries("brand", "{brand}", "engine", "{engine}").Methods(http.MethodGet)
	r.HandleFunc("/car/{id}", handler.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/car/{id}", handler.Update).Methods(http.MethodPut)
	r.HandleFunc("/car/{id}", handler.Delete).Methods(http.MethodDelete)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())

}
