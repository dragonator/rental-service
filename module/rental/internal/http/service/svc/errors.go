package svc

// Error represets a server error.
type Error struct {
	StatusCode int
	Message    string
}

// Error return the error message of the server error.
func (e *Error) Error() string {
	return e.Message
}
