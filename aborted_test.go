package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type AbortedErrorTest struct {
	err          *AbortedError
	appendErr    *AbortedError
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

var AbortedErrorTests = []AbortedErrorTest{
	{
		err:        NewAbortedError("Message 1"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 409: ABORTED. Message 1",
		getCode:    409,
		getMessage: "ABORTED. Message 1",
		getCause:   nil,
		json:       []byte(`{"errorCode":409,"errorMessage":"ABORTED. Message 1"}`),
		rpcCode:    codes.Aborted,
		rpcMessage: "ABORTED. Message 1",
	},
	{
		err:        NewAbortedError("Message 2", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 409: ABORTED. Message 2",
		getCode:    409,
		getMessage: "ABORTED. Message 2",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"errorCode":409,"errorMessage":"ABORTED. Message 2"}`),
		rpcCode:    codes.Aborted,
		rpcMessage: "ABORTED. Message 2",
	},
	{
		err:        NewAbortedError("Message 3", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 409: ABORTED. Message 3",
		getCode:    409,
		getMessage: "ABORTED. Message 3",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"errorCode":409,"errorMessage":"ABORTED. Message 3"}`),
		rpcCode:    codes.Aborted,
		rpcMessage: "ABORTED. Message 3",
	},
}

func TestAbortedErrorTimeout(t *testing.T) {
	for _, test := range AbortedErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestAbortedErrorTemporary(t *testing.T) {
	for _, test := range AbortedErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestAbortedErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range AbortedErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestAbortedErrorGetCode(t *testing.T) {
	for _, test := range AbortedErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestAbortedErrorGetMessage(t *testing.T) {
	for _, test := range AbortedErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestAbortedErrorGetCause(t *testing.T) {
	for _, test := range AbortedErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestAbortedErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range AbortedErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestAbortedErrorJson(t *testing.T) {
	for _, test := range AbortedErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestAbortedErrorGrpc(t *testing.T) {
	for _, test := range AbortedErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}

func TestAbortedErrorAppend(t *testing.T) {

	e1       := NewAbortedError("Message 1")
	e1append := e1.Append(errors.New("foo"))
	e1alt    := NewAbortedError("Message 1", errors.New("foo"))
	assert.Equal(t, e1alt.GetCause().Error(), e1append.GetCause().Error())

	e2       :=NewAbortedError("Message 2", errors.New("foo"))
	e2append := e2.Append(errors.New("bar"))
	e2alt    := NewAbortedError("Message 2", errors.New("foo"), errors.New("bar"))
	assert.Equal(t, e2alt.GetCause().Error(), e2append.GetCause().Error())
}
