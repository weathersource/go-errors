package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type NotImplementedErrorTest struct {
	err          *NotImplementedError
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

var NotImplementedErrorTests = []NotImplementedErrorTest{
	{
		err:        NewNotImplementedError("Message 1"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 501: NOT IMPLEMENTED. Message 1",
		getCode:    501,
		getMessage: "NOT IMPLEMENTED. Message 1",
		getCause:   nil,
		json:       []byte(`{"errorCode":501,"errorMessage":"NOT IMPLEMENTED. Message 1"}`),
		rpcCode:    codes.Unimplemented,
		rpcMessage: "NOT IMPLEMENTED. Message 1",
	},
	{
		err:        NewNotImplementedError("Message 2", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 501: NOT IMPLEMENTED. Message 2",
		getCode:    501,
		getMessage: "NOT IMPLEMENTED. Message 2",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"errorCode":501,"errorMessage":"NOT IMPLEMENTED. Message 2"}`),
		rpcCode:    codes.Unimplemented,
		rpcMessage: "NOT IMPLEMENTED. Message 2",
	},
	{
		err:        NewNotImplementedError("Message 2", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 501: NOT IMPLEMENTED. Message 2",
		getCode:    501,
		getMessage: "NOT IMPLEMENTED. Message 2",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"errorCode":501,"errorMessage":"NOT IMPLEMENTED. Message 2"}`),
		rpcCode:    codes.Unimplemented,
		rpcMessage: "NOT IMPLEMENTED. Message 2",
	},
}

func TestNotImplementedErrorTimeout(t *testing.T) {
	for _, test := range NotImplementedErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestNotImplementedErrorTemporary(t *testing.T) {
	for _, test := range NotImplementedErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestNotImplementedErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range NotImplementedErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestNotImplementedErrorGetCode(t *testing.T) {
	for _, test := range NotImplementedErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestNotImplementedErrorGetMessage(t *testing.T) {
	for _, test := range NotImplementedErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestNotImplementedErrorGetCause(t *testing.T) {
	for _, test := range NotImplementedErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestNotImplementedErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range NotImplementedErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestNotImplementedErrorJson(t *testing.T) {
	for _, test := range NotImplementedErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestNotImplementedErrorGrpc(t *testing.T) {
	for _, test := range NotImplementedErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
