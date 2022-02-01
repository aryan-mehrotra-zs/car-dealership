package errors

type EnumNotFound struct{}

func (e EnumNotFound) Error() string {
	return "Invalid enum value"
}
