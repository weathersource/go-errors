package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// FailedPreconditionError indicates operation was rejected because the
// system is not in a state required for the operation's execution.
// For example, directory to be deleted may be non-empty, an rmdir
// operation is applied to a non-directory, etc.
//
// Service implementorse should use FailedPreconditionError if the client
// should not retry until the system state has been explicitly fixed
// (E.g., if an "rmdir" fails because the directory is non-empty,
// FailedPreconditionError should be returned since the client should not
// retry unless they have first fixed up the directory by deleting files from
// it). If the client performs conditional REST Get/Update/Delete on a resource
// and the resource on the server does not match the condition,
// FailedPreconditionError should be used (e.g., conflicting read-modify-write
// on the same resource).
//
// Example error Message:
//
//		FAILED PRECONDITION. Resource xxx is a non-empty directory, so it cannot be deleted.
//
// HTTP Mapping: 400 BAD REQUEST
//
// RPC Mapping: FAILED_PRECONDITION
type FailedPreconditionError struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
	cause   error
	stack   stack
	rpcCode codes.Code
}

// NewFailedPreconditionError returns a new FailedPreconditionError.
func NewFailedPreconditionError(Message string, cause ...error) *FailedPreconditionError {
	var c error
	if len(cause) > 0 {
		c = NewErrors(cause...)
	}
	return &FailedPreconditionError{
		Code:    400,
		Message: "FAILED PRECONDITION. " + Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.FailedPrecondition,
	}
}

// Error implements the error interface
func (e *FailedPreconditionError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *FailedPreconditionError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *FailedPreconditionError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *FailedPreconditionError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *FailedPreconditionError) GetMessage() string { return e.Message }

// GetCause returns any causal errors associated with this error.
func (e *FailedPreconditionError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *FailedPreconditionError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *FailedPreconditionError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
