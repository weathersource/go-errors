package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type FailedPreconditionErrorTest struct {
	err          *FailedPreconditionError
	timeout      bool
	temporary    bool
	errorInfo    string
	errorVerbose string
	errorDebug   string
	errorTrace   string
	getCode      int
	getMessage   string
	getCause     error
	json         []byte
	rpcCode      codes.Code
	rpcMessage   string
}

var FailedPreconditionErrorTests = []FailedPreconditionErrorTest{
	{
		err:        NewFailedPreconditionError("Message 1"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 400: FAILED PRECONDITION. Message 1",
		getCode:    400,
		getMessage: "FAILED PRECONDITION. Message 1",
		getCause:   nil,
		json:       []byte(`{"errorCode":400,"errorMessage":"FAILED PRECONDITION. Message 1"}`),
		rpcCode:    codes.FailedPrecondition,
		rpcMessage: "FAILED PRECONDITION. Message 1",
	},
	{
		err:        NewFailedPreconditionError("Message 2", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 400: FAILED PRECONDITION. Message 2",
		getCode:    400,
		getMessage: "FAILED PRECONDITION. Message 2",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"errorCode":400,"errorMessage":"FAILED PRECONDITION. Message 2"}`),
		rpcCode:    codes.FailedPrecondition,
		rpcMessage: "FAILED PRECONDITION. Message 2",
	},
	{
		err:        NewFailedPreconditionError("Message 2", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 400: FAILED PRECONDITION. Message 2",
		getCode:    400,
		getMessage: "FAILED PRECONDITION. Message 2",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"errorCode":400,"errorMessage":"FAILED PRECONDITION. Message 2"}`),
		rpcCode:    codes.FailedPrecondition,
		rpcMessage: "FAILED PRECONDITION. Message 2",
	},
}

func TestFailedPreconditionErrorTimeout(t *testing.T) {
	for _, test := range FailedPreconditionErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestFailedPreconditionErrorTemporary(t *testing.T) {
	for _, test := range FailedPreconditionErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestFailedPreconditionErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range FailedPreconditionErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestFailedPreconditionErrorGetCode(t *testing.T) {
	for _, test := range FailedPreconditionErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestFailedPreconditionErrorGetMessage(t *testing.T) {
	for _, test := range FailedPreconditionErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestFailedPreconditionErrorGetCause(t *testing.T) {
	for _, test := range FailedPreconditionErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestFailedPreconditionErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range FailedPreconditionErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestFailedPreconditionErrorJson(t *testing.T) {
	for _, test := range FailedPreconditionErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestFailedPreconditionErrorGrpc(t *testing.T) {
	for _, test := range FailedPreconditionErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}

func TestFailedPreconditionErrorAppend(t *testing.T) {

	e1       := NewFailedPreconditionError("Message 1")
	e1append := e1.Append(errors.New("foo"))
	e1alt    := NewFailedPreconditionError("Message 1", errors.New("foo"))
	assert.Equal(t, e1alt.GetCause().Error(), e1append.GetCause().Error())

	e2       :=NewFailedPreconditionError("Message 2", errors.New("foo"))
	e2append := e2.Append(errors.New("bar"))
	e2alt    := NewFailedPreconditionError("Message 2", errors.New("foo"), errors.New("bar"))
	assert.Equal(t, e2alt.GetCause().Error(), e2append.GetCause().Error())
}
