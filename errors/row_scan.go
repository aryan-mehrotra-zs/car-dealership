package errors

type RowScan struct{}

func (e RowScan) Error() string {
	return "row scan error"
}
