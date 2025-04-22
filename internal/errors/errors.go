package errors

import (
	"errors"
	"fmt"
)

// Reference: errors
func New(message string) error {
	return errors.New(message)
}

// Wrap wraps the error with a message
func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// Reference: errors
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Reference: errors
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// Reference: errors
func As(err error, target any) bool {
	return errors.As(err, target)
}
