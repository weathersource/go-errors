package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// ResourceExhaustedError indicates some resource has been exhausted, perhaps
// a per-user quota, or perhaps the entire file system is out of space.
//
// Example error Message:
//
//		RESOURCE EXHAUSTED. Quota limit 'xxx' exceeded.
//
// HTTP Mapping: 429 TOO MANY REQUESTS
//
// RPC Mapping: RESOURCE_EXHAUSTED
type ResourceExhaustedError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
	cause   error
	stack   stack
	rpcCode codes.Code
}

// NewResourceExhaustedError returns a new ResourceExhaustedError.
func NewResourceExhaustedError(Message string, cause ...error) *ResourceExhaustedError {
	var c error
	if len(cause) > 0 {
		c = Errors(cause)
	}
	return &ResourceExhaustedError{
		Code:    429,
		Message: "RESOURCE EXHAUSTED. " + Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.ResourceExhausted,
	}
}

// Error implements the error interface
func (e *ResourceExhaustedError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *ResourceExhaustedError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *ResourceExhaustedError) Temporary() bool { return true }

// GetCode returns the HTTP status code associated with this error.
func (e *ResourceExhaustedError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *ResourceExhaustedError) GetMessage() string { return e.Message }

// GetCause returns any causal errors associated with this error.
func (e *ResourceExhaustedError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *ResourceExhaustedError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *ResourceExhaustedError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
