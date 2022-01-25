package models

import "github.com/google/uuid"

type Engine struct {
	ID           uuid.UUID `json:"-"`
	Displacement int       `json:"displacement"`
	NoOfCylinder int       `json:"noOfCylinder"`
	Range        int       `json:"range"`
}
