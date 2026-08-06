package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLastName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Username
		expectError bool
	}{
		{name: "valid minimal length", input: strings.Repeat("a", UsernameMinLength), expected: Username(strings.Repeat("a", UsernameMinLength))},
		{name: "valid max length", input: strings.Repeat("a", UsernameMaxLength), expected: Username(strings.Repeat("a", UsernameMaxLength))},
		{name: "trim spaces", input: "  Bob  ", expected: "Bob"},
		{name: "too short", input: strings.Repeat("a", UsernameMinLength-1), expectError: true},
		{name: "too long", input: strings.Repeat("a", UsernameMaxLength+1), expectError: true},
		{name: "empty string", input: "", expectError: true},
		{name: "only spaces", input: "     ", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewUsername(tt.input)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
