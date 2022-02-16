package car

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/amehrotra/car-dealership/errors"
)

// setStatusCode writes the status code based on the error type
func setStatusCode(w http.ResponseWriter, method string, data interface{}, err error) {
	switch err.(type) {
	case errors.EntityAlreadyExists:
		w.WriteHeader(http.StatusOK)
	case errors.MissingParam, errors.InvalidParam:
		w.WriteHeader(http.StatusBadRequest)
	case errors.EntityNotFound:
		w.WriteHeader(http.StatusNotFound)
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
		writeResponseBody(w, http.StatusCreated, data)
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
