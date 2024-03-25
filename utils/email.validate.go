package utils

import (
	"regexp"
)

// IsValidEmail checks if the provided email address is in a valid format.
func IsValidEmail(email string) bool {
	// Regex pattern for email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
