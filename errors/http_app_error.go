package errors

import (
	"fmt"
)

type HttpAppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHttpAppError(code int, message string) HttpAppError {
	return HttpAppError{Code: code, Message: message}
}

func (e HttpAppError) Error() string {
	return fmt.Sprintf("errors: %s.", e.Message)
}
