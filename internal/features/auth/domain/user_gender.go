package domain

import (
	"errors"
	"strings"
)

type UserGender string

const (
	GenderMale   UserGender = "male"
	GenderFemale UserGender = "female"
)

func NewUserGender(gender string) (UserGender, error) {
	g := UserGender(strings.ToLower(gender))

	switch g {
	case GenderMale, GenderFemale:
		return g, nil
	default:
		return "", errors.New("invalid user gender")
	}
}
