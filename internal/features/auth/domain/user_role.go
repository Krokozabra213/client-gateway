package domain

import (
	"errors"
	"strings"
)

type UserRole string

const (
	RoleRegular UserRole = "regular"
	RoleAdmin   UserRole = "admin"
)

func NewUserRole(role string) (UserRole, error) {
	r := UserRole(strings.ToLower(role))

	switch r {
	case RoleRegular, RoleAdmin:
		return r, nil
	default:
		return "", errors.New("invalid user role")
	}
}
