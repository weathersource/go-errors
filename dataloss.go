package errors

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// DataLossError indicates unrecoverable data loss or corruption.
//
// Since the client cannot fix this server error, it is not useful to generate
// additional error details. To avoid leaking sensitive information under error
// conditions, only a generic error Message is marshalled to JSON or returned
// via GRPC status.
//
// Error Message:
//
//		DATA LOSS. Unrecoverable data loss or data corruption.
//
// HTTP Mapping: 500 INTERNAL SERVER ERROR
//
// RPC Mapping: DATA_LOSS
type DataLossError struct {
	Code       int    `json:"error_code"`
	Message    string `json:"error_message"`
	logMessage string
	cause      error
	stack      stack
	rpcCode    codes.Code
}

// NewDataLossError returns a new DataLossError.
func NewDataLossError(Message string, cause ...error) *DataLossError {
	var c error
	if len(cause) > 0 {
		c = Errors(cause)
	}
	return &DataLossError{
		Code:       500,
		Message:    "DATA LOSS. Unrecoverable data loss or data corruption.",
		logMessage: Message,
		cause:      c,
		stack:      getTrace(),
		rpcCode:    codes.DataLoss,
	}
}

// Error implements the error interface
func (e *DataLossError) Error() string { return errorStr(e) }

// Timeout indicates if this error is the result of a timeout.
func (e *DataLossError) Timeout() bool { return false }

// Temporary indicates if this error is potentially recoverable.
func (e *DataLossError) Temporary() bool { return false }

// GetCode returns the HTTP status code associated with this error.
func (e *DataLossError) GetCode() int { return e.Code }

// GetMessage returns the message associated with this error.
func (e *DataLossError) GetMessage() string { return e.Message + " " + e.logMessage }

// GetCause returns any causal errors associated with this error.
func (e *DataLossError) GetCause() error { return e.cause }

// GetStack returns the trace stack associated with this error.
func (e *DataLossError) GetStack() stack { return e.stack }

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *DataLossError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}
