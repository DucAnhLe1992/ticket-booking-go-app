package apperr

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the standard error envelope returned by services.
type ErrorResponse struct {
	Message     string            `json:"message"`
	FieldErrors map[string]string `json:"fieldErrors,omitempty"`
}

// HTTPError is an application error with an associated HTTP status code.
type HTTPError interface {
	error
	StatusCode() int
	Response() ErrorResponse
}

type baseHTTPError struct {
	msg  string
	code int
}

func (e baseHTTPError) Error() string           { return e.msg }
func (e baseHTTPError) StatusCode() int         { return e.code }
func (e baseHTTPError) Response() ErrorResponse { return ErrorResponse{Message: e.msg} }

func WriteJSONError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	resp := ErrorResponse{Message: http.StatusText(status)}

	if he, ok := err.(HTTPError); ok {
		status = he.StatusCode()
		resp = he.Response()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

// Common error constructors
func NewBadRequest(msg string) HTTPError { return baseHTTPError{msg: msg, code: http.StatusBadRequest} }
func NewUnauthorized(msg string) HTTPError {
	return baseHTTPError{msg: msg, code: http.StatusUnauthorized}
}
func NewForbidden(msg string) HTTPError { return baseHTTPError{msg: msg, code: http.StatusForbidden} }
func NewNotFound(msg string) HTTPError  { return baseHTTPError{msg: msg, code: http.StatusNotFound} }
func NewConflict(msg string) HTTPError  { return baseHTTPError{msg: msg, code: http.StatusConflict} }
