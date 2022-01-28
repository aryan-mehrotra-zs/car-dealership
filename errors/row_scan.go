package errors

import (
	"fmt"
)

type RowScan struct{}

func (e RowScan) Error() string {
	return fmt.Sprintf("row scan error")
}
