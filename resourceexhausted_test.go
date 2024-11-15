package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type ResourceExhaustedErrorTest struct {
	err          *ResourceExhaustedError
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

var ResourceExhaustedErrorTests = []ResourceExhaustedErrorTest{
	{
		err:        NewResourceExhaustedError("Message 1"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 429: RESOURCE EXHAUSTED. Message 1",
		getCode:    429,
		getMessage: "RESOURCE EXHAUSTED. Message 1",
		getCause:   nil,
		json:       []byte(`{"errorCode":429,"errorMessage":"RESOURCE EXHAUSTED. Message 1"}`),
		rpcCode:    codes.ResourceExhausted,
		rpcMessage: "RESOURCE EXHAUSTED. Message 1",
	},
	{
		err:        NewResourceExhaustedError("Message 2", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 429: RESOURCE EXHAUSTED. Message 2",
		getCode:    429,
		getMessage: "RESOURCE EXHAUSTED. Message 2",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"errorCode":429,"errorMessage":"RESOURCE EXHAUSTED. Message 2"}`),
		rpcCode:    codes.ResourceExhausted,
		rpcMessage: "RESOURCE EXHAUSTED. Message 2",
	},
	{
		err:        NewResourceExhaustedError("Message 2", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 429: RESOURCE EXHAUSTED. Message 2",
		getCode:    429,
		getMessage: "RESOURCE EXHAUSTED. Message 2",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"errorCode":429,"errorMessage":"RESOURCE EXHAUSTED. Message 2"}`),
		rpcCode:    codes.ResourceExhausted,
		rpcMessage: "RESOURCE EXHAUSTED. Message 2",
	},
}

func TestResourceExhaustedErrorTimeout(t *testing.T) {
	for _, test := range ResourceExhaustedErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestResourceExhaustedErrorTemporary(t *testing.T) {
	for _, test := range ResourceExhaustedErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestResourceExhaustedErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range ResourceExhaustedErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestResourceExhaustedErrorGetCode(t *testing.T) {
	for _, test := range ResourceExhaustedErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestResourceExhaustedErrorGetMessage(t *testing.T) {
	for _, test := range ResourceExhaustedErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestResourceExhaustedErrorGetCause(t *testing.T) {
	for _, test := range ResourceExhaustedErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestResourceExhaustedErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range ResourceExhaustedErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestResourceExhaustedErrorJson(t *testing.T) {
	for _, test := range ResourceExhaustedErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestResourceExhaustedErrorGrpc(t *testing.T) {
	for _, test := range ResourceExhaustedErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}

func TestResourceExhaustedErrorAppend(t *testing.T) {

	e1       := NewResourceExhaustedError("Message 1")
	e1append := e1.Append(errors.New("foo"))
	e1alt    := NewResourceExhaustedError("Message 1", errors.New("foo"))
	assert.Equal(t, e1alt.GetCause().Error(), e1append.GetCause().Error())

	e2       :=NewResourceExhaustedError("Message 2", errors.New("foo"))
	e2append := e2.Append(errors.New("bar"))
	e2alt    := NewResourceExhaustedError("Message 2", errors.New("foo"), errors.New("bar"))
	assert.Equal(t, e2alt.GetCause().Error(), e2append.GetCause().Error())
}
