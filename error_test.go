package errors

import (
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type FooError struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
	cause   error
	stack   stack
	rpcCode codes.Code
}

func NewFooError(Message string, cause ...error) *FooError {
	var c error
	if len(cause) > 0 {
		c = cause[0]
	}
	return &FooError{
		Code:    999,
		Message: Message,
		cause:   c,
		stack:   getTrace(),
		rpcCode: codes.OK,
	}
}
func (e *FooError) Error() string      { return errorStr(e) }
func (e *FooError) Timeout() bool      { return false }
func (e *FooError) Temporary() bool    { return false }
func (e *FooError) GetCode() int       { return e.Code }
func (e *FooError) GetMessage() string { return e.Message }
func (e *FooError) GetCause() error    { return e.cause }
func (e *FooError) GetStack() stack    { return e.stack }
func (e *FooError) GRPCStatus() *status.Status {
	return status.New(e.rpcCode, e.Message)
}

type test struct {
	foo             *FooError
	expectedInfo    string
	expectedVerbose string
	expectedDebug   string
}

var tests = []test{
	{
		foo:             NewFooError("error Message 1"),
		expectedInfo:    "error 999: error Message 1",
		expectedVerbose: "error 999: error Message 1",
	},
	{
		foo:             NewFooError("error Message 2", errors.New("causal error")),
		expectedInfo:    "error 999: error Message 2",
		expectedVerbose: "error 999: error Message 2\ncause: causal error",
	},
}

func TestVerbosity(t *testing.T) {

	SetVerbosity(Info)
	for _, test := range tests {
		assert.Equal(t, test.expectedInfo, test.foo.Error())
	}

	SetVerbosity(Verbose)
	for _, test := range tests {
		assert.Equal(t, test.expectedVerbose, test.foo.Error())
	}

	// trace portion of this output prevents testing for string match in different contexts
	SetVerbosity(Debug)
	for _, test := range tests {
		assert.NotNil(t, test.foo.Error())
	}

	// trace portion of this output prevents testing for string match in different contexts
	SetVerbosity(Trace)
	for _, test := range tests {
		assert.NotNil(t, test.foo.Error())
	}
}
