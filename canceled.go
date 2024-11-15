package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// CanceledError indicates the operation was canceled (typically by the
// caller).
//
// Example error Message:
//
//		CANCELED.
//
// HTTP Mapping: 499 CLIENT CLOSED REQUEST
//
// RPC Mapping: CANCELED
type CanceledError struct {
	Code       int    `json:"errorCode"`
	Message    string `json:"errorMessage"`
	logMessage string
	cause      error
	stack      stack
	rpcCode    codes.Code
}

// NewCanceledError returns a new CanceledError.
func NewCanceledError(Message string, cause ...error) *CanceledError {
	var c error
	if len(cause) > 0 {
		c = NewErrors(cause...)
	}
	return &CanceledError{
		Code:       499,
		Message:    "CANCELED. Request canceled by the client.",
		logMessage: Message,
		cause:      c,
		stack:      getTrace(),
		rpcCode:    codes.Canceled,
	}
}

// Error implements the error interface
func (e *CanceledError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *CanceledError) Timeout() bool { return true }

// Temporary indicates if this error is potentially recoverable.
func (e *CanceledError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *CanceledError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *CanceledError) GetMessage() string { return e.Message + " " + e.logMessage }

// GetCause returns any causal errors associated with this error.
func (e *CanceledError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *CanceledError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *CanceledError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}

// appends additional error causes to this error
func (e *CanceledError) Append(errs ...error) *CanceledError {

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
