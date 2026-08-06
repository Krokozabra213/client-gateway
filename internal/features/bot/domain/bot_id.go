package domain

import "github.com/google/uuid"

type BotID uuid.UUID

func NewBotID(id string) (BotID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return BotID(uuid.Nil), err
	}

	return FromUuid(parsed), err
}

func FromUuid(id uuid.UUID) BotID {
	return BotID(id)
}

func (b BotID) String() string {
	return uuid.UUID(b).String()
}
