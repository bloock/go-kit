package domain

import (
	"fmt"
	httpError "github.com/bloock/go-kit/errors"
	"net/http"

	"github.com/google/uuid"
)

var (
	ErrIDEmpty   = httpError.NewHttpAppError(http.StatusBadRequest, "the provided ID is empty")
	ErrIDInvalid = httpError.NewHttpAppError(http.StatusBadRequest, "the provided ID should have UUID format")
	ErrIDNull    = httpError.NewHttpAppError(http.StatusBadRequest, "the provided ID is null UUID")
)

type UUID struct {
	id string
}

func NewUUID(id string) (UUID, error) {
	if err := validateID(id); err != nil {
		return UUID{}, fmt.Errorf("NewUUID: %w", err)
	}
	return UUID{id: id}, nil
}

func GenUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return id.String()
}

func (u UUID) ID() string {
	return u.id
}

func validateID(id string) error {
	if id == "" {
		return ErrIDEmpty
	}
	if _, err := uuid.Parse(id); err != nil {
		return ErrIDInvalid
	}
	return nil
}
