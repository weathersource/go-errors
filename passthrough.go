package errors

import (
	"context"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// NewPassthroughError handles an error from an external dependency. If the error is a
// timeout, canceled, unavailable, unknown, or temporary error, it is passed through as
// the appropriate correlating type from this package. Otherwise, an internal error
// with the provided message is returned.
func NewPassthroughError(msg string, err error) error {

	// test err against target interfaces
	sErr, sOk := err.(interface{ GRPCStatus() *(status.Status) })
	tiErr, tiOk := err.(interface{ Timeout() bool })
	teErr, teOk := err.(interface{ Temporary() bool })

	// CanceledError
	switch {
	case err == context.Canceled:
		fallthrough
	case sOk && codes.Canceled == sErr.GRPCStatus().Code():
		return NewCanceledError(msg, err)
	}

	// DeadlineExceededError
	switch {
	case err == context.DeadlineExceeded:
		fallthrough
	case sOk && codes.DeadlineExceeded == sErr.GRPCStatus().Code():
		fallthrough
	case tiOk && tiErr.Timeout():
		return NewDeadlineExceededError(msg, err)
	}

	// UnavailableError
	if sOk && codes.Unavailable == sErr.GRPCStatus().Code() {
		return NewUnavailableError(msg, err)
	}

	// UnknownError
	switch {
	case sOk && codes.Unknown == sErr.GRPCStatus().Code():
		fallthrough
	case teOk && teErr.Temporary():
		return NewUnknownError(msg, err)
	}

	// InternalError
	return NewInternalError(msg, err)
}
