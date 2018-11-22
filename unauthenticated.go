package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// UnauthenticatedError indicates the request does not have valid
// authentication credentials for the operation.
//
// Example error Message:
//
//		UNAUTHENTICATED. Invalid authentication credentials.
//
// HTTP Mapping: 401 UNAUTHORIZED
//
// RPC Mapping: UNAUTHENTICATED
type UnauthenticatedError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
	cause   error
	stack   stack
	rpcCode codes.Code
}

// NewUnauthenticatedError returns a new UnauthenticatedError.
func NewUnauthenticatedError(Message string, cause ...error) *UnauthenticatedError {
	var c error
	if len(cause) > 0 {
		c = Errors(cause)
	}
	return &UnauthenticatedError{
		Code:    401,
		Message: "UNAUTHENTICATED. " + Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.Unauthenticated,
	}
}

// Error implements the error interface
func (e *UnauthenticatedError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *UnauthenticatedError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *UnauthenticatedError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *UnauthenticatedError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *UnauthenticatedError) GetMessage() string { return e.Message }

// GetCause returns any causal errors associated with this error.
func (e *UnauthenticatedError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *UnauthenticatedError) GetStack() stack { return e.stack }

// GRPCStatus impliments an interface required to return proper GRPC status codes
func (e *UnauthenticatedError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
