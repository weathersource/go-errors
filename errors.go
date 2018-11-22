package errors

import (
	"fmt"
	"strings"
)

// Errors is a container for multiple errors and implements the error interface
type Errors []error

// New returns an error that consists of multiple errors.
func NewErrors(errs ...error) Errors {
	return Errors(errs)
}

// Error implements the error interface
func (e Errors) Error() string {
	lenE := len(e)
	if lenE <= 0 {
		return ""
	} else if lenE == 1 {
		return e[0].Error()
	}
	logWithNumber := make([]string, lenE)
	for i, l := range e {
		if l != nil {
			logWithNumber[i] = fmt.Sprintf("#%d: %s", i+1, l.Error())
		}
	}

	return fmt.Sprintf("Errors:\n%s", strings.Join(logWithNumber, "\n"))
}
