package contract

// ErrorResponse is a generic error response object.
type ErrorResponse struct {
	Error  string   `json:"error,omitempty"`
	Errors []string `json:"errors,omitempty"`
}
