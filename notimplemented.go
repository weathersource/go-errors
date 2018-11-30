package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// NotImplementedError indicates operation is not implemented or not
// supported/enabled in this service.
//
// Example error Message:
//
//		NOT IMPLEMENTED. Method 'xxx' not implemented.
//
// HTTP Mapping: 501 NOT IMPLEMENTED
//
// RPC Mapping: NOT_IMPLEMENTED
type NotImplementedError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
	cause   error
	stack   stack
	rpcCode codes.Code
}

// NewNotImplementedError returns a new NotImplementedError.
func NewNotImplementedError(Message string, cause ...error) *NotImplementedError {
	var c error
	if len(cause) > 0 {
		c = Errors(cause)
	}
	return &NotImplementedError{
		Code:    501,
		Message: "NOT IMPLEMENTED. " + Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.Unimplemented,
	}
}

// Error implements the error interface
func (e *NotImplementedError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *NotImplementedError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *NotImplementedError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *NotImplementedError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *NotImplementedError) GetMessage() string { return e.Message }

// GetCause returns any causal errors associated with this error.
func (e *NotImplementedError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *NotImplementedError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *NotImplementedError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
