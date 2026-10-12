package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestApplicationErrors(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		want string
	}{
		{name: "invalid credentials", err: ErrInvalidCredentials, want: "invalid credentials"},
		{name: "inactive user", err: ErrInactiveUser, want: "inactive user"},
		{name: "duplicate email", err: ErrEmailAlreadyExists, want: "email already exists"},
		{name: "missing user", err: ErrUserNotFound, want: "user not found"},
		{name: "admin access denied", err: ErrUnauthorizedAdminAccess, want: "admin user not found"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := testCase.err.Error(); got != testCase.want {
				t.Errorf("Error() = %q, want %q", got, testCase.want)
			}
			wrapped := fmt.Errorf("wrapped: %w", testCase.err)
			if !errors.Is(wrapped, testCase.err) {
				t.Errorf("errors.Is(%v, %v) = false, want true", wrapped, testCase.err)
			}
		})
	}
}
