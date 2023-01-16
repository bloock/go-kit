package domain

import (
	"fmt"
	"regexp"
)

const lengthSHA256 = "64"

type Hash struct {
	hash string
}

func NewHash(hash string) (*Hash, error) {
	matched, _ := regexp.MatchString("^[a-f0-9]{"+lengthSHA256+"}$", hash)
	if !matched {
		return nil, fmt.Errorf("invalid sha 256 hash: %s", hash)
	}

	return &Hash{hash: hash}, nil
}

func (h Hash) Hash() string {
	return h.hash
}
