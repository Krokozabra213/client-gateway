package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNonEmptyString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "valid string",
			input:       "hello",
			expected:    "hello",
			expectError: false,
		},
		{
			name:        "string with spaces",
			input:       "   hello   ",
			expected:    "hello",
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "only spaces",
			input:       "     ",
			expectError: true,
		},
		{
			name:        "string with tabs and newlines",
			input:       "\n\t hello \t",
			expected:    "hello",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewNonEmptyString(tt.input)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.Value())
		})
	}
}
