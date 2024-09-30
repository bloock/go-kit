package errors

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"net/http"
	"strings"
)

var (
	ErrNotFound = NewHttpAppError(http.StatusNotFound, "not found")
)

func ErrInvalidBodyJSON(err error) HttpAppError {
	return NewHttpAppError(http.StatusBadRequest, err.Error())
}

func ErrUnexpected(err error) HttpAppError {
	sentry.CaptureException(err)
	return NewHttpAppError(http.StatusInternalServerError, err.Error())
}

func ErrUnauthorized(err error) HttpAppError {
	return NewHttpAppError(http.StatusUnauthorized, err.Error())
}

func ErrForbidden(err error) HttpAppError {
	return NewHttpAppError(http.StatusForbidden, err.Error())
}

func ErrComponentNotFound(component string) HttpAppError {
	return NewHttpAppError(http.StatusNotFound, fmt.Sprintf("%s not found", component))
}

func WrapSqlRepositoryError(err error) HttpAppError {
	switch {
	case strings.Contains(err.Error(), "no rows in result set"):
		return ErrNotFound
	default:
		return ErrUnexpected(err)
	}
}
