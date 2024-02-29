package domain

import (
	"errors"
	httpError "github.com/bloock/go-kit/errors"
)

type SubscriptionStripeStatus int32

var (
	ErrInvalidSubscriptionStripeStatus = httpError.ErrInvalidBodyJSON(errors.New("invalid subscription stripe status"))
)

const (
	Trialing SubscriptionStripeStatus = iota
	Active
	Incomplete
	IncompleteExpired
	PastDue
	Canceled
	Unpaid
	Paused
)

func NewSubscriptionStripeStatus(s string) (SubscriptionStripeStatus, error) {
	switch s {
	case "trialing":
		return Trialing, nil
	case "active":
		return Active, nil
	case "incomplete":
		return Incomplete, nil
	case "incomplete_expired":
		return IncompleteExpired, nil
	case "past_due":
		return PastDue, nil
	case "canceled":
		return Canceled, nil
	case "unpaid":
		return Unpaid, nil
	case "paused":
		return Paused, nil
	default:
		return -1, ErrInvalidSubscriptionStripeStatus
	}
}

func (s SubscriptionStripeStatus) String() string {
	switch s {
	case Trialing:
		return "trialing"
	case Active:
		return "active"
	case Incomplete:
		return "incomplete"
	case IncompleteExpired:
		return "incomplete_expired"
	case PastDue:
		return "past_due"
	case Canceled:
		return "canceled"
	case Unpaid:
		return "unpaid"
	case Paused:
		return "paused"
	default:
		return ""
	}
}
