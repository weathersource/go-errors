package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// CancelledError indicates the operation was canceled (typically by the
// caller).
//
// Example error Message:
//
//		CANCELLED.
//
// HTTP Mapping: 499 CLIENT CLOSED REQUEST
//
// RPC Mapping: CANCELLED
type CancelledError struct {
	Code       int    `json:"errorCode"`
	Message    string `json:"errorMessage"`
	logMessage string
	cause      error
	stack      stack
	rpcCode    codes.Code
}

// NewCancelledError returns a new CancelledError.
func NewCancelledError(Message string, cause ...error) *CancelledError {
	var c error
	if len(cause) > 0 {
		c = NewErrors(cause...)
	}
	return &CancelledError{
		Code:       499,
		Message:    "CANCELLED. Request cancelled by the client.",
		logMessage: Message,
		cause:      c,
		stack:      getTrace(),
		rpcCode:    codes.Canceled,
	}
}

// Error implements the error interface
func (e *CancelledError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *CancelledError) Timeout() bool { return true }

// Temporary indicates if this error is potentially recoverable.
func (e *CancelledError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *CancelledError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *CancelledError) GetMessage() string { return e.Message + " " + e.logMessage }

// GetCause returns any causal errors associated with this error.
func (e *CancelledError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *CancelledError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *CancelledError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
