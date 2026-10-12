package validation

import "testing"

func TestIsValidPassword(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		want     bool
	}{
		{name: "contains all required character classes", password: "Secret123!", want: true},
		{name: "missing uppercase", password: "secret123!", want: false},
		{name: "missing lowercase", password: "SECRET123!", want: false},
		{name: "missing digit", password: "SecretPass!", want: false},
		{name: "missing special character", password: "Secret1234", want: false},
		{name: "special character from allowed set", password: "Secret123?", want: true},
		{name: "special character outside allowed set", password: "Secret123#", want: false},
		{name: "empty password", password: "", want: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := IsValidPassword(testCase.password); got != testCase.want {
				t.Errorf("IsValidPassword(%q) = %t, want %t", testCase.password, got, testCase.want)
			}
		})
	}
}
