package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// DeadlineExceededError means operation expired before completion.
// For operations that change the state of the system, this error may be
// returned even if the operation has completed successfully. For
// example, a successful response from a server could have been delayed
// long enough for the deadline to expire.
//
// Since the client cannot fix this server error, it is not useful to generate
// additional error details. To avoid leaking sensitive information under error
// conditions, only a generic error Message is marshalled to JSON or returned
// via GRPC status.
//
// Error Message:
//
//		DEADLINE EXCEEDED. Server timeout.
//
// HTTP Mapping: 504 GATEWAY TIMEOUT
//
// RPC Mapping: DEADLINE_EXCEEDED
type DeadlineExceededError struct {
	Code       int    `json:"errorCode"`
	Message    string `json:"errorMessage"`
	logMessage string
	cause      error
	stack      stack
	rpcCode    codes.Code
}

// NewDeadlineExceededError returns a new DeadlineExceededError.
func NewDeadlineExceededError(Message string, cause ...error) *DeadlineExceededError {
	var c error
	if len(cause) > 0 {
		c = NewErrors(cause...)
	}
	return &DeadlineExceededError{
		Code:       504,
		Message:    "DEADLINE EXCEEDED. Server timeout.",
		logMessage: Message,
		cause:      c,
		stack:      getTrace(),
		rpcCode:    codes.DeadlineExceeded,
	}
}

// Error implements the error interface
func (e *DeadlineExceededError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *DeadlineExceededError) Timeout() bool { return true }

// Temporary indicates if this error is potentially recoverable.
func (e *DeadlineExceededError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *DeadlineExceededError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *DeadlineExceededError) GetMessage() string { return e.Message + " " + e.logMessage }

// GetCause returns any causal errors associated with this error.
func (e *DeadlineExceededError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *DeadlineExceededError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *DeadlineExceededError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}

// appends additional error causes to this error
func (e *DeadlineExceededError) Append(errs ...error) *DeadlineExceededError {

	if e.cause == nil {
		e.cause = NewErrors(errs...)
	} else {
		c, ok := e.cause.(*Errors)
		if ok {
			c.Append(errs...)
			e.cause = c
		}
	}
	return e
}
