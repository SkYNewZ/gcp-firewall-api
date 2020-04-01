package models

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/api/googleapi"
)

// ApplicationError describe a global application error
type ApplicationError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// JSON return ApplicationError as JSON
func (e *ApplicationError) Error() string {
	res, _ := json.Marshal(e)
	return string(res)
}

// NewGoogleApplicationError describe a http error response from Google
func NewGoogleApplicationError(err *googleapi.Error) *ApplicationError {
	return &ApplicationError{
		Code:    err.Code,
		Message: fmt.Sprintf("Google error: %s", err.Message),
	}
}

// NewBadTokenError describe a http error when decoding JWT
func NewBadTokenError(message ...string) *ApplicationError {
	e := NewBadRequestError()
	e.Message = "Invalid Bearer token. Please make sure you are using 'gcloud auth print-identity-token'"

	if len(message) > 0 {
		e.Message = fmt.Sprintf("%s: %s", e.Message, message[0])
	}

	return e
}

// NewForbiddenError describe a http error response 403 Forbidden
func NewForbiddenError(message ...string) *ApplicationError {
	e := &ApplicationError{
		Code:    http.StatusForbidden,
		Message: http.StatusText(http.StatusForbidden),
	}

	if len(message) > 0 {
		e.Message = message[0]
	}

	return e
}

// NewInternalError describe a http error response 500 Internal Server Error
func NewInternalError() *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}
}

// NewMethodNotAllowedError describe a http error response 405 Method Not Allowed
func NewMethodNotAllowedError() *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusMethodNotAllowed,
		Message: http.StatusText(http.StatusMethodNotAllowed),
	}
}

// NewNotFoundError describe a http error response 404 Not Found
func NewNotFoundError() *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}
}

// NewBadRequestError describe a http error response 400 Bad Request
func NewBadRequestError(message ...string) *ApplicationError {
	e := &ApplicationError{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
	}

	if len(message) > 0 {
		e.Message = message[0]
	}

	return e
}
