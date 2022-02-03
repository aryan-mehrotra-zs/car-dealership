package errors

type InvalidJson struct {
}

func (e InvalidJson) Error() string {
	return "invalid json"
}
