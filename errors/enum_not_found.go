package errors

import (
	"fmt"
)

type EnumNotFound struct{}

func (e EnumNotFound) Error() string {
	return fmt.Sprintf("Invalid enum value")
}
