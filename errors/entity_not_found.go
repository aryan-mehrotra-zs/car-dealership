package errors

type EntityNotFound struct {
}

func (e EntityNotFound) Error() string {
	return "Entity not found"
}
