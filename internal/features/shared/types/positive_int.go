package types

import commontypes "github.com/Krokozabra213/common/types"

// PositiveInt is a type alias for the common PositiveInt type.
type PositiveInt[T commontypes.Integer] = commontypes.PositiveInt[T]

// NewPositiveInt is re-exported from common/types.
func NewPositiveInt[T commontypes.Integer](value T) (PositiveInt[T], error) {
	return commontypes.NewPositiveInt[T](value)
}
