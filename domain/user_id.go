package domain

import (
	"github.com/bloock/go-kit/errors"
	"net/http"
)

var (
	ErrInvalidUserID = errors.NewHttpAppError(http.StatusBadRequest, "the user ID should have UUID format")
)

type UserID struct {
	id UUID
}

func NewUserIDStr(id string) (UserID, error) {
	uid, err := NewUUID(id)
	if err != nil {
		return UserID{}, ErrInvalidUserID
	}
	return UserID{id: uid}, nil
}

func NewUserID() UserID {
	uuid := GenUUID()
	return UserID{
		id: UUID{id: uuid},
	}
}

func (i UserID) ID() string {
	return i.id.ID()
}
