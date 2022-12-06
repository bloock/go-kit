package shareddomain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrIDEmpty   = errors.New("the provided ID is empty")
	ErrIDInvalid = errors.New("the provided ID should have UUID format")
	ErrIDNull    = errors.New("the provided ID is null UUID")
)

func NewUUID(id string) (UUID, error) {
	if err := validateID(id); err != nil {
		return UUID{}, fmt.Errorf("NewUUID: %w", err)
	}
	return UUID{id: id}, nil
}

func GenUUID() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return uuid.String()
}

func (u UUID) ID() string {
	return u.id
}

func validateID(id string) error {
	if id == "" {
		return ErrIDEmpty
	}
	if id == "00000000-0000-0000-0000-000000000000" {
		return ErrIDNull
	}
	if _, err := uuid.Parse(id); err != nil {
		return ErrIDInvalid
	}
	return nil
}
