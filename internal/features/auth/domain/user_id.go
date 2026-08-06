package domain

import "github.com/google/uuid"

type UserID uuid.UUID

func NewUserId(id string) (UserID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return UserID(uuid.Nil), err
	}

	return FromUuid(parsed), err
}

func FromUuid(id uuid.UUID) UserID {
	return UserID(id)
}

func (b UserID) String() string {
	return uuid.UUID(b).String()
}
