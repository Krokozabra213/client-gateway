package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPositiveInt64(t *testing.T) {
	tests := []struct {
		name        string
		input       int64
		expected    int64
		expectError bool
	}{
		{
			name:     "valid positive",
			input:    10,
			expected: 10,
		},
		{
			name:     "valid minimal positive",
			input:    1,
			expected: 1,
		},
		{
			name:        "zero",
			input:       0,
			expectError: true,
		},
		{
			name:        "negative",
			input:       -5,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewPositiveInt[int64](tt.input)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.Value())
		})
	}
}
