package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUserRole(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    UserRole
		expectError bool
	}{
		{name: "regular lowercase", input: "regular", expected: RoleRegular},
		{name: "admin lowercase", input: "admin", expected: RoleAdmin},
		{name: "regular uppercase", input: "REGULAR", expected: RoleRegular},
		{name: "admin mixed case", input: "AdMiN", expected: RoleAdmin},
		{name: "invalid role", input: "superuser", expectError: true},
		{name: "empty string", input: "", expectError: true},
		{name: "spaces", input: "   ", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewUserRole(tt.input)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
