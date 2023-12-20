package utils

import (
	"regexp"
	"strings"
)

// IsValidEmail checks if the provided string matches the pattern for a valid email address.
// It ensures the email address has a standard structure with username, domain, and top-level domain.
func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

// IsValidPassword checks if the provided string meets password complexity requirements.
// It enforces the password to contain at least one lowercase letter, one uppercase letter,
// one numeric digit, one special character, and have a minimum length of 8 characters.
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper, hasLower, hasDigit, hasSpecial bool
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case strings.ContainsAny(string(char), "@$!%*?&"):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}
