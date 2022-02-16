package errors

import "fmt"

type EntityAlreadyExists struct {
	Entity string
}

func (e EntityAlreadyExists) Error() string {
	return fmt.Sprintf("entity  %v already exists", e.Entity)
}
