package identity_entity

import (
	"github.com/bloock/go-kit/errors"
	"net/http"
)

type PublishIntervalMinutes int

var ErrInvalidPublishIntervalMinutes = errors.NewHttpAppError(http.StatusBadRequest, "publish interval minutes not supported")

const (
	PublishIntervalMinutes1 PublishIntervalMinutes = iota
	PublishIntervalMinutes5
	PublishIntervalMinutes15
	PublishIntervalMinutes60
)

func NewPublishIntervalMinutes(_type int) (PublishIntervalMinutes, error) {
	switch _type {
	case 1:
		return PublishIntervalMinutes1, nil
	case 5:
		return PublishIntervalMinutes5, nil
	case 15:
		return PublishIntervalMinutes15, nil
	case 60:
		return PublishIntervalMinutes60, nil
	default:
		return -1, ErrInvalidPublishIntervalMinutes
	}
}

func (p PublishIntervalMinutes) Int() int {
	switch p {
	case PublishIntervalMinutes1:
		return 1
	case PublishIntervalMinutes5:
		return 5
	case PublishIntervalMinutes15:
		return 15
	case PublishIntervalMinutes60:
		return 60
	default:
		return 0
	}
}
