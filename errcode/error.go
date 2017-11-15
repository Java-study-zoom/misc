package errcode

import (
	"fmt"
)

// Error is a generic error with a string error code.
type Error struct {
	Code  string // code is the type of the error.
	error        // err is the error message, human friendly.
}

// Common general error codes
const (
	NotFound     = "not-found"
	InvalidArg   = "invalid-arg"
	Internal     = "internal"
	Unauthorized = "unauthorized"
)

// Add creates a new error with code as the error code.
func Add(code string, err error) *Error {
	return &Error{
		Code:  code,
		error: err,
	}
}

// Of returns the code of the error. For errors that
// do not have a code, it returns empty string.
func Of(err error) string {
	if codedErr, ok := err.(*Error); ok {
		return codedErr.Code
	}
	return ""
}

// IsNotFound checks if it is a not-found error.
func IsNotFound(err error) bool {
	return Of(err) == NotFound
}

// IsInvalidArg checks if it is an invalid argument error.
func IsInvalidArg(err error) bool {
	return Of(err) == InvalidArg
}

// IsInternal checks if it is an internal error.
func IsInternal(err error) bool {
	return Of(err) == Internal
}

// IsUnauthorized checks if it is an unauthorized error.
func IsUnauthorized(err error) bool {
	return Of(err) == Unauthorized
}

// Errorf creates an Error with the given error code.
func Errorf(code string, f string, args ...interface{}) *Error {
	return Add(code, fmt.Errorf(f, args...))
}

// AltErrorf replaces the message of e, but keeps the error code.
func AltErrorf(err error, msg string, args ...interface{}) error {
	cerr, ok := err.(*Error)
	if !ok {
		return fmt.Errorf(msg, args...)
	}

	return Errorf(cerr.Code, msg, args...)
}

// NotFoundf creates a new not-found error.
func NotFoundf(f string, args ...interface{}) *Error {
	return Errorf(NotFound, f, args...)
}

// InvalidArgf creates a new invalid arugment error.
func InvalidArgf(f string, args ...interface{}) *Error {
	return Errorf(InvalidArg, f, args...)
}

// Internalf creates a new internal error.
func Internalf(f string, args ...interface{}) *Error {
	return Errorf(Internal, f, args...)
}

// Unauthorizedf returns an error caused by unauthrozied request.
func Unauthorizedf(f string, args ...interface{}) *Error {
	return Errorf(Unauthorized, f, args...)
}
