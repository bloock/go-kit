package errors

import (
	"context"
	"fmt"
	"github.com/bloock/go-kit/observability"
)

type HttpAppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHttpAppError(code int, message string) HttpAppError {
	return HttpAppError{Code: code, Message: message}
}

func (e HttpAppError) Error() string {
	return fmt.Sprintf("errors: %s. code: %d.", e.Message, e.Code)
}

func (e HttpAppError) LogErrorTraceability(ctx context.Context, logger observability.Logger) HttpAppError {
	logger.Info(ctx).Err(e).Msg("")
	return e
}
