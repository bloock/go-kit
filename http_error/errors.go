package httperror

import "net/http"

var (
	ErrNotFound = NewAppError(http.StatusNotFound, "not found")
)

func ErrInvalidBodyJSON(err error) *appError {
	return NewAppError(http.StatusBadRequest, err.Error())
}
