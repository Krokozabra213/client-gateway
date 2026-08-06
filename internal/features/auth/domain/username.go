package domain

import (
	"errors"
	"strings"
)

type Username string

const (
	UsernameMinLength int = 1
	UsernameMaxLength int = 50
)

func NewUsername(name string) (Username, error) {
	n := Username(strings.TrimSpace(name))

	if len(n) < UsernameMinLength || len(n) > UsernameMaxLength {
		return "", errors.New("invalid username length")
	}

	return n, nil
}
