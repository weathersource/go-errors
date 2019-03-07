package errors

import (
	"fmt"
	"strings"
	"sync"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Errors is a container for multiple errors and implements the error interface
type Errors struct {
	sync.RWMutex
	errs []error
}

// NewErrors returns an error that consists of multiple errors.
func NewErrors(errs ...error) *Errors {
	e := Errors{errs: errs}
	return &e
}

// Error implements the error interface
func (e Errors) Error() string {
	if len(e.errs) <= 0 {
		return ""
	} else if len(e.errs) == 1 {
		return e.errs[0].Error()
	}
	logWithNumber := make([]string, len(e.errs))
	for i, l := range e.errs {
		if l != nil {
			logWithNumber[i] = fmt.Sprintf("#%d: %s", i+1, l.Error())
		}
	}

	return fmt.Sprintf("MULTIPLE ERRORS.\n%s", strings.Join(logWithNumber, "\n"))
}

// Len returns the number of errors in e
func (e *Errors) Len() int {
	if e == nil {
		return 0
	}
	if len(e.errs) == 1 {
		err := e.peek()
		if err == nil {
			return 0
		}
	}
	return len(e.errs)
}

// Pop removes and returns the last error from Errors
func (e *Errors) Append(err error) {
	if err == nil {
		return
	}
	e.Lock()
	defer e.Unlock()
	e.errs = append(e.errs, err)
	return
}

// Pop removes and returns the last error from Errors
func (e *Errors) Pop() error {
	e.Lock()
	defer e.Unlock()
	if e.Len() > 0 {
		err := (e.errs)[len(e.errs)-1]
		e.errs = (e.errs)[:len(e.errs)-1]
		return err
	}
	return nil
}

// Shift removes and returns the first error from Errors
func (e *Errors) Shift() error {
	e.Lock()
	defer e.Unlock()
	if e.Len() > 0 {
		err := (e.errs)[0]
		e.errs = (e.errs)[1:]
		return err
	}
	return nil
}

// Timeout indicates if this error is the result of a timeout.
func (e *Errors) Timeout() bool {
	err := e.peekLocked()
	if err != nil {
		wxErr, ok := err.(interface{ Timeout() bool })
		if ok {
			return wxErr.Timeout()
		}
	}
	return false
}

// Temporary indicates if this error is potentially recoverable.
func (e *Errors) Temporary() bool {
	err := e.peekLocked()
	if err != nil {
		wxErr, ok := err.(interface{ Temporary() bool })
		if ok {
			return wxErr.Temporary()
		}
	}
	return false
}

// GetCode returns the HTTP status code associated with this error.
func (e *Errors) GetCode() int {

	if e == nil {
		return 200
	}

	e.RLock()
	defer e.RUnlock()
	if e.Len() == 0 {
		return 200
	} else if e.Len() == 1 {

		err := e.peek()

		wxErr, ok := err.(interface{ GetCode() int })
		if ok {
			return wxErr.GetCode()
		}

		// validate this is a gRPC error
		_, ok = err.(interface{ GRPCStatus() *status.Status })
		if ok {
			switch status.Code(err) {
			case codes.Aborted:
				return 409
			case codes.AlreadyExists:
				return 409
			case codes.Canceled:
				return 499
			case codes.DataLoss:
				return 500
			case codes.DeadlineExceeded:
				return 504
			case codes.FailedPrecondition:
				return 400
			case codes.Internal:
				return 500
			case codes.InvalidArgument:
				return 400
			case codes.NotFound:
				return 404
			case codes.OutOfRange:
				return 400
			case codes.PermissionDenied:
				return 403
			case codes.ResourceExhausted:
				return 429
			case codes.Unauthenticated:
				return 401
			case codes.Unavailable:
				return 503
			case codes.Unimplemented:
				return 501
			case codes.Unknown:
				return 500
			}
		}
	}

	return 500
}

// GetMessage returns the message associated with this error.
func (e *Errors) GetMessage() string {

	if e == nil {
		return ""
	}

	e.RLock()
	defer e.RUnlock()
	if e.Len() == 0 {
		return ""
	} else if e.Len() == 1 {
		err := e.peek()
		wxErr, ok := err.(interface{ GetMessage() string })
		if ok {
			return wxErr.GetMessage()
		}
		return err.Error()
	}
	return "MULTIPLE ERRORS."
}

// GetCause returns any causal errors associated with this error.
func (e *Errors) GetCause() error {

	if e == nil {
		return nil
	}

	e.RLock()
	defer e.RUnlock()
	if e.Len() == 0 {
		return nil
	} else if e.Len() == 1 {
		err := e.peek()
		wxErr, ok := err.(interface{ GetCause() error })
		if ok {
			return wxErr.GetCause()
		}
		return nil
	}
	return e
}

// GetStack returns the trace stack associated with this error.
func (e *Errors) GetStack() stack {

	var s stack

	if e == nil {
		return s
	}

	e.RLock()
	defer e.RUnlock()
	if e.Len() == 1 {
		err := e.peek()
		wxErr, ok := err.(interface{ GetStack() stack })
		if ok {
			return wxErr.GetStack()
		}
	}
	return s
}

// GRPCStatus implements an interface required to return proper GRPC status codes
func (e *Errors) GRPCStatus() *status.Status {

	if e == nil {
		return nil
	}

	e.RLock()
	defer e.RUnlock()
	if e.Len() == 0 {
		return nil
	} else if e.Len() == 1 {
		err := e.peek()
		grpcErr, ok := err.(interface{ GRPCStatus() *status.Status })
		if ok {
			return grpcErr.GRPCStatus()
		}
	}
	return status.New(codes.Unknown, e.Error())
}

// peek returns the first error in e, but leaves it in the slice
func (e *Errors) peek() error {
	if e == nil {
		return nil
	}
	if len(e.errs) > 0 {
		return e.errs[0]
	}
	return nil
}

// peek returns the first error in e, but leaves it in the slice with a read lock.
func (e *Errors) peekLocked() error {
	e.RLock()
	defer e.RUnlock()
	err := e.peek()
	return err
}
