package types

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/amehrotra/car-dealership/errors"
)

type Fuel int

const (
	Diesel Fuel = iota + 1
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
	case "diesel:":
		*f = Diesel
	case "electric":
		*f = Electric
	default:
		return errors.InvalidParam{Param: []string{"fuelType"}}
	}

	return nil
}

func (f Fuel) Value() (driver.Value, error) {
	return int64(f), nil
}

func (f *Fuel) Scan(value interface{}) error {
	switch value {
	case "petrol":
		*f = Fuel(1)
		return nil
	case "diesel":
		*f = Fuel(2)
		return nil
	case "ev":
		*f = Fuel(3)
		return nil
	}

	return nil
}
