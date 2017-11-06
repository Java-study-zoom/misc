package gerror

import (
	"fmt"
	"strings"
)

// Error is a generic type of error
type Error struct {
	code string // code is the type of the error.
	err  error  // err is the error message, human friendly.

}

func codeErrorf(code string, msg string, args ...interface{}) *Error {
	err := fmt.Errorf(msg, args...)
	return &Error{
		code: code,
		err:  err,
	}
}

// AltError replace the message of a Gerror, keep code the same.
func AltError(e *Error, msg string, args ...interface{}) *Error {
	return codeErrorf(e.code, msg, args...)
}

// NotFound returned a not found error.
func NotFound(msg string) *Error {
	return codeErrorf("not-found", "Not Found: %s", msg)
}

// Invalid returned an error of invalid argument.
func Invalid(msg string) *Error {
	return codeErrorf("invalid", "Invalid argument: %s", msg)

}

// Internal returned an internal error.
func Internal(msg string) *Error {
	return codeErrorf("interal", "Internal error: %s", msg)
}

// Unauthorized returned an error caused by unauthrozied request.
func Unauthorized(msg string) *Error {
	return codeErrorf("unauthrozied", "Unauthorized: %s", msg)
}

// Temperary returned an temperory error
func Temperary(msg string) *Error {
	return codeErrorf("temperary", "Temperary Unavailable: %s", msg)
}

// Other returned an errors is not defined.
func Other(msg string) *Error {
	return codeErrorf("other", "Other errors: %s", msg)
}

// GetError convert an error into a GError.
func GetError(err error) *Error {
	gErr, ok := err.(*Error)
	if ok {
		return gErr
	}
	msg := err.Error()
	errType := strings.SplitN(msg, ":", 2)[0]
	switch errType {
	case "Not Found":
		return NotFound(msg)
	case "Invalid argument":
		return Invalid(msg)
	case "Internal error":
		return Internal(msg)
	case "Unauthorized":
		return Unauthorized(msg)
	default:
		return Other(msg)
	}
}

func (e *Error) Error() string {
	return e.err.Error()
}
