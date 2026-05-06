package utils

import (
	"errors"
	"net/http"
)

// Custom error types
var (
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
	ErrValidation     = errors.New("validation error")
	ErrInternalServer = errors.New("internal server error")
)

// ErrorWithCode wraps an error with HTTP status code
type ErrorWithCode struct {
	Err     error
	Code    int
	Message string
}

// GetHTTPStatus returns appropriate HTTP status code for an error
func GetHTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	case ErrValidation:
		return http.StatusBadRequest
	case ErrInternalServer:
		return http.StatusInternalServerError
	default:
		// Check if it's a custom ErrorWithCode
		var ewc *ErrorWithCode
		if errors.As(err, &ewc) {
			return ewc.Code
		}
		return http.StatusInternalServerError
	}
}

// GetErrorMessage returns a user-friendly error message
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	switch err {
	case ErrUnauthorized:
		return "Authentication required or invalid credentials"
	case ErrForbidden:
		return "You do not have permission to access this resource"
	case ErrNotFound:
		return "Resource not found"
	case ErrConflict:
		return "Resource already exists"
	case ErrValidation:
		return "Invalid input provided"
	case ErrInternalServer:
		return "An internal server error occurred"
	default:
		// Check if it's a custom ErrorWithCode
		var ewc *ErrorWithCode
		if errors.As(err, &ewc) {
			return ewc.Message
		}
		return err.Error()
	}
}

// NewErrorWithCode creates a new error with custom code and message
func NewErrorWithCode(err error, code int, message string) *ErrorWithCode {
	return &ErrorWithCode{
		Err:     err,
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface
func (e *ErrorWithCode) Error() string {
	return e.Message
}

// Unwrap implements the Unwrap interface for errors.Is and errors.As
func (e *ErrorWithCode) Unwrap() error {
	return e.Err
}
