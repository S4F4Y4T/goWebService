package apperror

import (
	"errors"
	"fmt"
)

// Category is a type for classifying errors
type Category string

const (
	NotFound        Category = "NOT_FOUND"
	BadRequest      Category = "BAD_REQUEST"
	Conflict        Category = "CONFLICT"
	Internal        Category = "INTERNAL"
	Unauthorized    Category = "UNAUTHORIZED"
	Forbidden       Category = "FORBIDDEN"
)

// AppError represents a structured error in the application
type AppError struct {
	Category Category
	Message  string
	Err      error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Category, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Category, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Helper functions to create specific errors
func New(cat Category, msg string, err error) error {
	return &AppError{
		Category: cat,
		Message:  msg,
		Err:      err,
	}
}

// Is Category checks if the error belongs to a specific category
func IsCategory(err error, cat Category) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Category == cat
	}
	return false
}
