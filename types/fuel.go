package types

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/amehrotra/car-dealership/errors"
)

type Fuel int

const (
	Diesel Fuel = iota
	Petrol
	Electric
)

func (f Fuel) MarshalJSON() ([]byte, error) {
	var s string
	switch f {
	case Petrol:
		s = "petrol"
	case Diesel:
		s = "diesel"
	case Electric:
		s = "electric"
	default:
		return nil, errors.InvalidParam{Param: []string{"fuelType"}}
	}

	return json.Marshal(s)
}

func (f *Fuel) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case "petrol":
		*f = Petrol
	case "diesel":
		*f = Diesel
	case "electric":
		*f = Electric
	default:
		return errors.InvalidParam{Param: []string{"fuelType"}}
	}

	return nil
}

func (f Fuel) Value() (driver.Value, error) {
	switch f {
	case Diesel:
		return "diesel", nil
	case Petrol:
		return "petrol", nil
	case Electric:
		return "electric", nil
	}

	return nil, errors.InvalidParam{Param: []string{"fuelType"}}
}

func (f *Fuel) Scan(value interface{}) error {
	fuel, ok := value.([]byte)
	if !ok {
		return errors.InvalidParam{Param: []string{"fuelType"}}
	}

	switch string(fuel) {
	case "diesel":
		*f = Diesel
	case "petrol":
		*f = Petrol
	case "electric":
		*f = Electric
	default:
		return errors.InvalidParam{Param: []string{"fuelType"}}
	}

	return nil
}
