package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// PermissionDeniedError indicates the caller does not have permission to
// execute the specified operation. It must not be used for rejections
// caused by exhausting some resource (use ResourceExhaustedError
// instead for those errors). It must not be used if the caller cannot be
// identified (use UnauthenticatedError instead for those errors).
//
// Example error Message:
//
//		PERMISSION DENIED. Permission 'xxx' denied on file 'yyy'.
//
// HTTP Mapping: 403 FORBIDDEN
//
// RPC Mapping: PERMISSION_DENIED
type PermissionDeniedError struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
	cause   error
	stack   stack
	rpcCode codes.Code
}

// NewPermissionDeniedError returns a new PermissionDeniedError.
func NewPermissionDeniedError(Message string, cause ...error) *PermissionDeniedError {
	var c error
	if len(cause) > 0 {
		c = NewErrors(cause...)
	}
	return &PermissionDeniedError{
		Code:    403,
		Message: "PERMISSION DENIED. " + Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.PermissionDenied,
	}
}

// Error implements the error interface
func (e *PermissionDeniedError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *PermissionDeniedError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *PermissionDeniedError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *PermissionDeniedError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *PermissionDeniedError) GetMessage() string { return e.Message }

// GetCause returns any causal errors associated with this error.
func (e *PermissionDeniedError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *PermissionDeniedError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *PermissionDeniedError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}

// appends additional error causes to this error
func (e *PermissionDeniedError) Append(errs ...error) *PermissionDeniedError {

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
