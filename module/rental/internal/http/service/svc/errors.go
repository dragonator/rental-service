package svc

import "net/http"

// Predefined errors.
var (
	ErrInvalidQueryParameters = &Error{StatusCode: http.StatusBadRequest, Message: "invalid query parameters"}
)

// Error represets a server error.
type Error struct {
	StatusCode int
	Message    string
}

// Error return the error message of the server error.
func (e *Error) Error() string {
	return e.Message
}
