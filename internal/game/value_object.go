package game

import (
	"errors"
)

var ErrInvalidIdentity = errors.New("identity must container only lower-case letters, and '-'")

// Identity is a validated identifier containing only lower-case letters and '-'.
type Identity struct {
	value string
}

func (i Identity) String() string { return i.value }

// NewIdentity returns an Identity, or an error if value is not a valid identifier.
func NewIdentity(value string) (Identity, error) {
	if !isValidIdentity(value) {
		return Identity{}, ErrInvalidIdentity
	}
	return Identity{value: value}, nil
}

func isValidIdentity(value string) bool {
	for _, v := range value {
		if (v < 'a' || v > 'z') && v != '-' {
			return false
		}
	}
	return value != ""
}
