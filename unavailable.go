package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// UnavailableError indicates the service is currently unavailable.
// This is a most likely a transient condition and may be corrected
// by retrying with a backoff.
//
// A litmus test that may help a service implementor in deciding
// between FailedPreconditionError, AbortedError, and UnavailableError:
//
//  (a) Use UnavailableError if the client can retry just the failing call.
//  (b) Use AbortedError if the client should retry at a higher-level
//      (e.g., restarting a read-modify-write sequence).
//  (c) Use FailedPreconditionError if the client should not retry until
//      the system state has been explicitly fixed. E.g., if an "rmdir"
//      fails because the directory is non-empty, FailedPreconditionError
//      should be returned since the client should not retry unless
//      they have first fixed up the directory by deleting files from it.
//  (d) Use FailedPreconditionError if the client performs conditional
//      REST Get/Update/Delete on a resource and the resource on the
//      server does not match the condition. E.g., conflicting
//      read-modify-write on the same resource.
//
// Since the client cannot fix this server error, it is not useful to generate
// additional error details. To avoid leaking sensitive information under error
// conditions, only a generic error Message is marshalled to JSON or returned
// via GRPC status.
//
// Error Message:
//
//		UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance.
//
// HTTP Mapping: 503 SERVICE UNAVAILABLE
//
// RPC Mapping: UNAVAILABLE
type UnavailableError struct {
	Code       int    `json:"error_code"`
	Message    string `json:"error_message"`
	logMessage string
	cause      error
	stack      stack
	rpcCode    codes.Code
}

// NewUnavailableError returns a new UnavailableError.
func NewUnavailableError(Message string, cause ...error) *UnavailableError {
	var c error
	if len(cause) > 0 {
		c = Errors(cause)
	}
	return &UnavailableError{
		Code:       503,
		Message:    "UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance.",
		logMessage: Message,
		cause:      c,
		stack:      getTrace(),
		rpcCode:    codes.Unavailable,
	}
}

// Error implements the error interface
func (e *UnavailableError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *UnavailableError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *UnavailableError) Temporary() bool { return true }

// GetCode returns the HTTP status code associated with this error.
func (e *UnavailableError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *UnavailableError) GetMessage() string { return e.Message + " " + e.logMessage }

// GetCause returns any causal errors associated with this error.
func (e *UnavailableError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *UnavailableError) GetStack() stack { return e.stack }

// GRPCStatus impliments an interface required to return proper GRPC status codes
func (e *UnavailableError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
