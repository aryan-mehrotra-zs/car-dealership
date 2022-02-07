package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/amehrotra/car-dealership/drivers"
	handlers "github.com/amehrotra/car-dealership/handlers/car"
	"github.com/amehrotra/car-dealership/middlewares"
	services "github.com/amehrotra/car-dealership/services/car"
	"github.com/amehrotra/car-dealership/stores/car"
	"github.com/amehrotra/car-dealership/stores/engine"
)

func main() {
	// database connection
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

	// dependency injection
	carStore := car.New(db)
	engineStore := engine.New(db)
	service := services.New(engineStore, carStore)
	handler := handlers.New(service)

	r := mux.NewRouter()
	r.HandleFunc("/car", handler.Create).Methods(http.MethodPost)
	r.HandleFunc("/car", handler.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/car/{id}", handler.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/car/{id}", handler.Update).Methods(http.MethodPut)
	r.HandleFunc("/car/{id}", handler.Delete).Methods(http.MethodDelete)

	// authentication middleware
	r.Use(middlewares.AuthMiddleware)

	// setup server variables
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
	}

	// start server
	log.Println(srv.ListenAndServe())
}
