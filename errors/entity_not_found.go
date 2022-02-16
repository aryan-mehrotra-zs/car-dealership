package errors

import "fmt"

type EntityNotFound struct {
	Entity string
	ID     string
}

func (e EntityNotFound) Error() string {
	return fmt.Sprintf("entity %s with id %s not found", e.Entity, e.ID)
}
