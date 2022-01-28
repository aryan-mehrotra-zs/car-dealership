package types

import (
	"database/sql/driver"
)

type Fuel int

const (
	Diesel Fuel = iota + 1
	Petrol
	Electric
)

func (f Fuel) Value() (driver.Value, error) {
	return int64(f), nil
}

func (f *Fuel) Scan(value interface{}) error {
	switch value {
	case "petrol":
		*f = Fuel(1)
		return nil
	}
	return nil
}
