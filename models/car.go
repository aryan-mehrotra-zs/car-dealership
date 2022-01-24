package models

import (
	"github.com/google/uuid"

	_type "github.com/amehrotra/car-dealership/type"
)

type Car struct {
	ID                uuid.UUID  `json:"id"`
	Model             string     `json:"model"`
	YearOfManufacture int        `json:"yearOfManufacture"`
	Brand             string     `json:"brand"`
	FuelType          _type.Fuel `json:"fuelType"`
	Engine            Engine     `json:"engine"`
}
