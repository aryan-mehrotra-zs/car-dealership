package errors

import "fmt"

type InvalidParam struct {
	Param string
}

func (e InvalidParam) Error() string {
	return fmt.Sprintf("parameter %s is invalid", e.Param)
}
