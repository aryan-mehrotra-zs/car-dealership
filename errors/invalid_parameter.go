package errors

type InvalidParam struct {
}

func (e InvalidParam) Error() string {
	return "Invalid Param"
}
