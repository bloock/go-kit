package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotFound = NewHttpAppError(http.StatusNotFound, "not found")
)

func ErrInvalidBodyJSON(err error) HttpAppError {
	return NewHttpAppError(http.StatusBadRequest, err.Error())
}

func ErrUnexpected(err error) HttpAppError {
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
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	default:
		return ErrUnexpected(err)
	}
}