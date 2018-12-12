package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// AlreadyExistsError means an attempt to create an entity failed because one
// already exists.
//
// Example error Message:
//
//		ALREADY EXISTS. Resource 'xxx' already exists.
//
// HTTP Mapping: 409 CONFLICT
//
// RPC Mapping: ALREADY_EXISTS
type AlreadyExistsError struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
	cause   error
	stack   stack
	rpcCode codes.Code
}

// NewAlreadyExistsError returns a new AlreadyExistsError.
func NewAlreadyExistsError(Message string, cause ...error) *AlreadyExistsError {
	var c error
	if len(cause) > 0 {
		c = Errors(cause)
	}
	return &AlreadyExistsError{
		Code:    409,
		Message: "ALREADY EXISTS. " + Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.AlreadyExists,
	}
}

// Error implements the error interface
func (e *AlreadyExistsError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *AlreadyExistsError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *AlreadyExistsError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *AlreadyExistsError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *AlreadyExistsError) GetMessage() string { return e.Message }

// GetCause returns any causal errors associated with this error.
func (e *AlreadyExistsError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *AlreadyExistsError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *AlreadyExistsError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
