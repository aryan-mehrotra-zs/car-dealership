package models

import "github.com/google/uuid"

type Engine struct {
	ID           uuid.UUID `json:"-id"`
	Displacement int       `json:"displacement"`
	NCylinder    int       `json:"noOfCylinder"`
	Range        int       `json:"range"
`
}
