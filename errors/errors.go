package errors

import "net/http"

var (
	ErrNotFound = NewHttpAppError(http.StatusNotFound, "not found")
)

func ErrInvalidBodyJSON(err error) HttpAppError {
	return NewHttpAppError(http.StatusBadRequest, err.Error())
}
