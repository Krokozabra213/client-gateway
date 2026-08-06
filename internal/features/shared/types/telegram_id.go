package types

import (
	"errors"

	commontypes "github.com/Krokozabra213/common/types"
)

type TelegramID struct {
	commontypes.PositiveInt[int64]
}

func NewTelegramID(id int64) (TelegramID, error) {
	p, err := NewPositiveInt(id)
	if err != nil {
		return TelegramID{}, errors.New("TelegramID must be positive")
	}

	return TelegramID{p}, nil
}
