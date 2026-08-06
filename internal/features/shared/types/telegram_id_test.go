package types

import (
	"testing"
)

func TestNewTelegramID(t *testing.T) {
	tests := []struct {
		name        string
		input       int64
		expected    int64
		expectError bool
	}{
		{
			name:     "valid positive id",
			input:    1,
			expected: 1,
		},
		{
			name:     "valid large id",
			input:    1234567890,
			expected: 1234567890,
		},
		{
			name:        "zero id",
			input:       0,
			expectError: true,
		},
		{
			name:        "negative id",
			input:       -42,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewTelegramID(tt.input)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Value() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result.Value())
			}
		})
	}
}
