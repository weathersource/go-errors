package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// UnknownError is an unknown server error. An example of where this error may
// be returned is if a Status value received from another address space belongs
// to an error-space that is not known in this address space.
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
//		UNKNOWN ERROR.
//
// HTTP Mapping: 500 INTERNAL SERVER ERROR
//
// RPC Mapping: UNKNOWN
type UnknownError struct {
	Code       int    `json:"errorCode"`
	Message    string `json:"errorMessage"`
	logMessage string
	cause      error
	stack      stack
	rpcCode    codes.Code
}

// NewUnknownError returns a new UnknownError.
func NewUnknownError(Message string, cause ...error) *UnknownError {
	var c error
	if len(cause) > 0 {
		c = NewErrors(cause...)
	}
	return &UnknownError{
		Code:    500,
		Message: "UNKNOWN ERROR. " + Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.Unknown,
	}
}

// Error implements the error interface
func (e *UnknownError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *UnknownError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *UnknownError) Temporary() bool { return true }

// GetCode returns the HTTP status code associated with this error.
func (e *UnknownError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *UnknownError) GetMessage() string { return e.Message }

// GetCause returns any causal errors associated with this error.
func (e *UnknownError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *UnknownError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *UnknownError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}

// appends additional error causes to this error
func (e *UnknownError) Append(errs ...error) *UnknownError {

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
