package helpers

import (
	"fmt"
	"net/http"
)

// Error type
type Error struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Trace    error  `json:"trace"`
	HTTPCode int    `json:"-"`
}

// Pretty print of the error
func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Message)
}

// ErrorTrace is tracing the error
func (e Error) ErrorTrace() error {
	return e.Trace
}

// ErrorWithCode is creating an error with code
func ErrorWithCode(code string, message string, trace error) Error {
	return Error{Code: code, Message: message, Trace: trace}
}

// NewError is creating an error with code and HTTP code
func NewError(httpCode int, code string, message string, trace error) Error {
	return Error{Code: code, Message: message, HTTPCode: httpCode, Trace: trace}
}

// NoCtxError defines an error that can be contextualized with a trace
type NoCtxError func(error) Error

// declare is used to declare const errors
func declare(code string, httpCode int, message string) NoCtxError {
	return func(err error) Error {
		return Error{Code: code, HTTPCode: httpCode, Message: message, Trace: err}
	}
}

var (
	// ErrorInvalidInput occurs when input is invalid
	ErrorInvalidInput = declare("invalid_input", http.StatusBadRequest, "Failed to bind the body data")
	// ErrorInternal occurs when there is internal server error
	ErrorInternal = declare("internal_error", http.StatusBadRequest, "Internal error")
	// ValidationError is for body validation
	ValidationError = declare("validation_error", http.StatusBadRequest, "Validation error")

	// ErrorTokenGenAccess occurs when access token failed to be generated
	ErrorTokenGenAccess = declare("token_generation_failed", http.StatusInternalServerError, "Could not generate the access token")
	// ErrorTokenGenRefresh occurs when refresh token failed to be generated
	ErrorTokenGenRefresh = declare("token_generation_failed", http.StatusInternalServerError, "Could not generate the refresh token")
	// ErrorTokenRefreshInvalid occurs when the refresh token is invalid
	ErrorTokenRefreshInvalid = declare("refresh_token_invalid", http.StatusBadRequest, "Refresh token invalid")

	// ErrorResourceNotFound occurs
	ErrorResourceNotFound = declare("resource_not_found", http.StatusNotFound, "Resource does not exist")

	// ErrorUserUnauthorized occurs when user doesn't have enough permissions to access au resource
	ErrorUserUnauthorized = Error{Code: "user_unauthorized", HTTPCode: http.StatusUnauthorized, Message: "Insufficient permissions to access this resource"}
	// ErrorUserUpdate occurs when user failed to be updated
	ErrorUserUpdate = declare("update_user_failed", http.StatusInternalServerError, "Could not update the user")
	// ErrorUserNotExist occurs when user doesn't exist
	ErrorUserNotExist = declare("user_does_not_exist", http.StatusNotFound, "User does not exist")
	// ErrorUserWrongPassword occurs when provided password is incorrect
	ErrorUserWrongPassword = declare("incorrect_password", http.StatusUnauthorized, "Password is not correct")
	// ErrorUserNotActivated occurs when user needs to activate its account via email
	ErrorUserNotActivated = declare("user_needs_activation", http.StatusNotFound, "User needs to be activated via email")

	// ErrorInvalidToken occurs when token is invalid (for example expired)
	ErrorInvalidToken = declare("invalid_token", http.StatusBadRequest, "The given token is invalid")

	// ErrorFileOpening occurs when file opening fails
	ErrorFileOpening = declare("file_opening_error", http.StatusNotAcceptable, "File opening error")
	// ErrorFileParsing occurs when file parsing fails
	ErrorFileParsing = declare("file_parsing_error", http.StatusNotAcceptable, "File parsing error")
)
