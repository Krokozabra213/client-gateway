package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewBotID(t *testing.T) {
	validUUID := uuid.New().String()
	invalidUUID := "not-a-uuid"

	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{name: "valid uuid", input: validUUID},
		{name: "invalid uuid", input: invalidUUID, expectError: true},
		{name: "empty string", input: "", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBotID(tt.input)
			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, BotID(uuid.Nil), result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.input, uuid.UUID(result).String())
		})
	}
}

func TestFromUuid(t *testing.T) {
	id := uuid.New()
	userID := FromUuid(id)

	require.Equal(t, id, uuid.UUID(userID))
}
