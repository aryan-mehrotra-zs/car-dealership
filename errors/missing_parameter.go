package errors

import "fmt"

type MissingParam struct {
	Param string
}

func (e MissingParam) Error() string {
	return fmt.Sprintf("parameter %s is missing", e.Param)
}
