package helpers

import "fmt"

// Error type
type Error struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Trace    error  `json:"trace"`
	HttpCode int    `json:"-"`
}

// Pretty print of the error
func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Message)
}

// Tracing the error
func (e Error) ErrorTrace() error {
	return e.Trace
}

// Creating an error with code
func ErrorWithCode(code string, message string, trace error) Error {
	return Error{Code: code, Message: message, Trace: trace}
}

// Creating an error with code and HTTP code
func NewError(httpCode int, code string, message string, trace error) Error {
	return Error{Code: code, Message: message, HttpCode: httpCode, Trace: trace}
}
