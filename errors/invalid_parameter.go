package errors

import (
	"fmt"
	"strings"
)

type InvalidParam struct {
	Param []string
}

func (e InvalidParam) Error() string {
	switch len(e.Param) {
	case 0:
		return "invalid parameters were sent for this request"
	case 1:
		return fmt.Sprintf("parameter %s is invalid", e.Param[0])
	default:
		return fmt.Sprintf("parameters %s are invalid", strings.Join(e.Param, ", "))
	}
}
