package common

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return NewValidationError("email is required")
	}

	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	if !matched {
		return NewValidationError("email format is invalid")
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return NewValidationError("password must be at least 8 characters long")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit {
		return NewValidationError("password must contain uppercase, lowercase, and digit")
	}

	return nil
}

func ValidateNonEmpty(fieldName, value string) error {
	if strings.TrimSpace(value) == "" {
		return NewValidationError(fmt.Sprintf("%s is required", fieldName))
	}
	return nil
}

func ValidateNonEmptyMultiple(fields map[string]string) error {
	for name, value := range fields {
		if err := ValidateNonEmpty(name, value); err != nil {
			return err
		}
	}
	return nil
}

func ValidateMinLength(fieldName string, value string, minLen int) error {
	if len(strings.TrimSpace(value)) < minLen {
		return NewValidationError(fmt.Sprintf("%s must be at least %d characters", fieldName, minLen))
	}
	return nil
}

func ValidateMaxLength(fieldName string, value string, maxLen int) error {
	if len(value) > maxLen {
		return NewValidationError(fmt.Sprintf("%s must not exceed %d characters", fieldName, maxLen))
	}
	return nil
}
