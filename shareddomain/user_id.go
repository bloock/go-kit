package shareddomain

import (
	"errors"
)

type UserID struct {
	id UUID
}

var (
	ErrInvalidUserID = errors.New("the user ID should have UUID format")
)

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

type UUID struct {
	id string
}
