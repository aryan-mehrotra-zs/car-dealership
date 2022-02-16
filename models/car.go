package models

import (
	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/types"
)

type Car struct {
	ID              uuid.UUID  `json:"id"`
	Model           string     `json:"model"`
	ManufactureYear int        `json:"yearOfManufacture"`
	Brand           string     `json:"brand"`
	FuelType        types.Fuel `json:"fuelType"`
	Engine          Engine     `json:"engine"`
}
