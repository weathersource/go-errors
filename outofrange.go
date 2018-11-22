package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// OutOfRangeError means operation was attempted past the valid range.
// E.g., seeking or reading past end of file.
//
// Unlike InvalidArgumentError, this error indicates a problem that may
// be fixed if the system state changes. For example, a 32-bit file
// system will generate InvalidArgument if asked to read at an
// offset that is not in the range [0,2^32-1], but it will generate
// OutOfRangeError if asked to read from an offset past the current
// file size.
//
// There is a fair bit of overlap between FailedPreconditionError and
// OutOfRangeError. We recommend using OutOfRangeError (the more specific
// error) when it applies so that callers who are iterating through
// a space can easily look for an OutOfRangeError to detect when
// they are done.
//
// Example error Message:
//
//		OUT OF RANGE. Parameter 'age' is out of range [0, 125].
//
// HTTP Mapping: 400 BAD REQUEST
//
// RPC Mapping: OUT_OF_RANGE
type OutOfRangeError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
	cause   error
	stack   stack
	rpcCode codes.Code
}

// NewOutOfRangeError returns a new OutOfRangeError.
func NewOutOfRangeError(Message string, cause ...error) *OutOfRangeError {
	var c error
	if len(cause) > 0 {
		c = Errors(cause)
	}
	return &OutOfRangeError{
		Code:    400,
		Message: "OUT OF RANGE. " + Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.OutOfRange,
	}
}

// Error implements the error interface
func (e *OutOfRangeError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *OutOfRangeError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *OutOfRangeError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *OutOfRangeError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *OutOfRangeError) GetMessage() string { return e.Message }

// GetCause returns any causal errors associated with this error.
func (e *OutOfRangeError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *OutOfRangeError) GetStack() stack { return e.stack }

// GRPCStatus impliments an interface required to return proper GRPC status codes
func (e *OutOfRangeError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
