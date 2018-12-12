package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// NotFoundError means some requested entity (e.g., file or directory) was
// not found.
//
// Example error Message:
//
//		NOT FOUND. Resource 'xxx' not found.
//
// HTTP Mapping: 404 NOT FOUND
//
// RPC Mapping: NOT_FOUND
type NotFoundError struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
	cause   error
	stack   stack
	rpcCode codes.Code
}

// NewNotFoundError returns a new NotFoundError.
func NewNotFoundError(Message string, cause ...error) *NotFoundError {
	var c error
	if len(cause) > 0 {
		c = Errors(cause)
	}
	return &NotFoundError{
		Code:    404,
		Message: "NOT FOUND. " + Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.NotFound,
	}
}

// Error implements the error interface
func (e *NotFoundError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *NotFoundError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *NotFoundError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *NotFoundError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *NotFoundError) GetMessage() string { return e.Message }

// GetCause returns any causal errors associated with this error.
func (e *NotFoundError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *NotFoundError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *NotFoundError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
