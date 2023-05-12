package key_entity

import (
	"net/http"

	"github.com/bloock/go-kit/errors"
)

type KeyProtection int

var ErrInvalidKeyProtection = errors.NewHttpAppError(http.StatusBadRequest, "key protection not supported")

const (
	SoftwareProtected KeyProtection = iota
	//HSMProtected
)

func NewKeyProtection(_type int) (KeyProtection, error) {
	switch _type {
	case 1:
		return SoftwareProtected, nil
	/*case 2:
		return HSMProtected, nil*/
	default:
		return -1, ErrInvalidKeyProtection
	}
}

func (k KeyProtection) Int() int {
	switch k {
	case SoftwareProtected:
		return 1
		/*case HSMProtected:
		return 2*/
	default:
		return 0
	}
}
