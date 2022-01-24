package errors

type MissingParam struct {
}

func (e MissingParam) Error() string {
	return "Missing Param"
}
