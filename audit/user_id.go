package audit

import (
	"errors"

	"github.com/google/uuid"
)

type UserID struct {
	id uuid.UUID
}

var (
	ErrInvalidUserID = errors.New("The user ID should have UUID format")
)

func NewUserID(id string) (UserID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return UserID{}, ErrInvalidUserID
	}
	return UserID{id: uuid}, err
}

func (i UserID) ID() string {
	return i.id.String()
}
