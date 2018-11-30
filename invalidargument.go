package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// InvalidArgumentError indicates client specified an invalid argument.
// Note that this differs from FailedPreconditionError. It indicates
// arguments that are problematic regardless of the state of the system
// (e.g., a malformed file name).
//
// Example error Message:
//
//		INVALID ARGUMENT. Request field x.y.z is xxx, expected one of [yyy, zzz].
//
// HTTP Mapping: 400 BAD REQUEST
//
// RPC Mapping: INVALID_ARGUMENT
type InvalidArgumentError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
	cause   error
	stack   stack
	rpcCode codes.Code
}

// NewInvalidArgumentError returns a new InvalidArgumentError.
func NewInvalidArgumentError(Message string, cause ...error) *InvalidArgumentError {
	var c error
	if len(cause) > 0 {
		c = Errors(cause)
	}
	return &InvalidArgumentError{
		Code:    400,
		Message: "INVALID ARGUMENT. " + Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.InvalidArgument,
	}
}

// Error implements the error interface
func (e *InvalidArgumentError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *InvalidArgumentError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *InvalidArgumentError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *InvalidArgumentError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *InvalidArgumentError) GetMessage() string { return e.Message }

// GetCause returns any causal errors associated with this error.
func (e *InvalidArgumentError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *InvalidArgumentError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *InvalidArgumentError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
