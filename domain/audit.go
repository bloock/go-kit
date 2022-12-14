package domain

import (
	"time"
)

type Audit struct {
	userID    UserID
	updateAT  time.Time
	createdAT time.Time
}

func (a Audit) UserId() string {
	return a.userID.ID()
}

func (a Audit) UpdateAt() time.Time {
	return a.updateAT
}

func (a Audit) CreateAt() time.Time {
	return a.createdAT
}

func NewAudit(user_id string, update_at time.Time, created_at time.Time) (Audit, error) {
	userUUID, err := NewUserIDStr(user_id)
	if err != nil {
		return Audit{}, err
	}
	return Audit{
		userID:    userUUID,
		updateAT:  update_at,
		createdAT: created_at,
	}, nil
}
