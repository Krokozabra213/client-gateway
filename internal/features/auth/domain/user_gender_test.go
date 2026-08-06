package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUserGender(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    UserGender
		expectError bool
	}{
		{name: "male lowercase", input: "male", expected: GenderMale},
		{name: "female lowercase", input: "female", expected: GenderFemale},
		{name: "male uppercase", input: "MALE", expected: GenderMale},
		{name: "female mixed case", input: "FeMaLe", expected: GenderFemale},
		{name: "invalid gender", input: "unknown", expectError: true},
		{name: "empty string", input: "", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewUserGender(tt.input)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
