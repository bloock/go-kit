package errors

import "fmt"

type AppError struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func NewAppError(message string, requestID string) *AppError {
	return &AppError{Message: message, RequestID: requestID}
}

func (e AppError) Error() string {
	return fmt.Sprintf("errors: %s. RequestID: %s", e.Error(), e.RequestID)
}
