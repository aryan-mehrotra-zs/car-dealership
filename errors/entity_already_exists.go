package errors

import "fmt"

type EntityAlreadyExists struct {
	Entity string
}

func (e EntityAlreadyExists) Error() string {
	return fmt.Sprintf("Entity  %v Already Exists", e.Entity)
}
