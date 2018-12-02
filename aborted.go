package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// AbortedError indicates the operation was aborted, typically due to a
// concurrency issue like sequencer check failures, transaction aborts,
// etc.
//
// A litmus test that may help a service implementor in deciding
// between ResourceExhaustedError, UnavailableError, and AbortedError:
//
//  (a) Use ResourceExhaustedError for client errors like exceeding allowed
//      rate limits. The client may retry the failing call after they have
//      resolved the causal issue.
//  (b) Use UnavailableError for server errors like inability to accomodate
//      current load or planned server maintenance. The client may retry the
//      failing call.
//  (c) Use AbortedError if the client should retry at a higher-level
//      (e.g., restarting a read-modify-write sequence).
//
// Example error Message:
//
//		ABORTED. Couldn’t acquire lock on resource ‘xxx’.
//
// HTTP Mapping: 409 CONFLICT
//
// RPC Mapping: ABORTED
type AbortedError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
	cause   error
	stack   stack
	rpcCode codes.Code
}

// NewAbortedError returns a new AbortedError.
func NewAbortedError(Message string, cause ...error) *AbortedError {
	var c error
	if len(cause) > 0 {
		c = Errors(cause)
	}
	return &AbortedError{
		Code:    409,
		Message: "ABORTED. " + Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.Aborted,
	}
}

// Error implements the error interface
func (e *AbortedError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *AbortedError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *AbortedError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *AbortedError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *AbortedError) GetMessage() string { return e.Message }

// GetCause returns any causal errors associated with this error.
func (e *AbortedError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *AbortedError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *AbortedError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
