package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// InternalError means some invariants expected by underlying
// system has been broken. If you see one of these errors,
// something is very broken.
//
// A litmus test that may help a service implementor in deciding between
// UnknownError and InternalError:
//
//  (a) Use UnknownError for generic server-side errors that may be recoverable
//      (UnknownError.Temporary() will return true).
//  (b) Use InternalError for generic server-side errors that are not recoverable
//      (InternalError.Temporary() will return false).
//
// Since the client cannot fix this server error, it is not useful to generate
// additional error details. To avoid leaking sensitive information under error
// conditions, only a generic error Message is marshalled to JSON or returned
// via GRPC status.
//
// Error Message:
//
//		INTERNAL ERROR.
//
// HTTP Mapping: 500 INTERNAL SERVER ERROR
//
// RPC Mapping: INTERNAL
type InternalError struct {
	Code       int    `json:"errorCode"`
	Message    string `json:"errorMessage"`
	logMessage string
	cause      error
	stack      stack
	rpcCode    codes.Code
}

// NewInternalError returns a new InternalError.
func NewInternalError(Message string, cause ...error) *InternalError {
	var c error
	if len(cause) > 0 {
		c = NewErrors(cause...)
	}
	return &InternalError{
		Code:       500,
		Message:    "INTERNAL ERROR.",
		logMessage: Message,
		cause:      c,
		stack:      getTrace(),
		rpcCode:    codes.Internal,
	}
}

// Error implements the error interface
func (e *InternalError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *InternalError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *InternalError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *InternalError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *InternalError) GetMessage() string { return e.Message + " " + e.logMessage }

// GetCause returns any causal errors associated with this error.
func (e *InternalError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *InternalError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *InternalError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}

// appends additional error causes to this error
func (e *InternalError) Append(errs ...error) *InternalError {

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
