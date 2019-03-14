package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type UnknownErrorTest struct {
	err          *UnknownError
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

var UnknownErrorTests = []UnknownErrorTest{
	{
		err:        NewUnknownError("foo"),
		timeout:    false,
		temporary:  true,
		errorInfo:  "error 500: UNKNOWN ERROR. foo",
		getCode:    500,
		getMessage: "UNKNOWN ERROR. foo",
		getCause:   nil,
		json:       []byte(`{"errorCode":500,"errorMessage":"UNKNOWN ERROR. foo"}`),
		rpcCode:    codes.Unknown,
		rpcMessage: "UNKNOWN ERROR. foo",
	},
	{
		err:        NewUnknownError("foo", errors.New("causal error")),
		timeout:    false,
		temporary:  true,
		errorInfo:  "error 500: UNKNOWN ERROR. foo",
		getCode:    500,
		getMessage: "UNKNOWN ERROR. foo",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"errorCode":500,"errorMessage":"UNKNOWN ERROR. foo"}`),
		rpcCode:    codes.Unknown,
		rpcMessage: "UNKNOWN ERROR. foo",
	},
	{
		err:        NewUnknownError("foo", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  true,
		errorInfo:  "error 500: UNKNOWN ERROR. foo",
		getCode:    500,
		getMessage: "UNKNOWN ERROR. foo",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"errorCode":500,"errorMessage":"UNKNOWN ERROR. foo"}`),
		rpcCode:    codes.Unknown,
		rpcMessage: "UNKNOWN ERROR. foo",
	},
}

func TestUnknownErrorTimeout(t *testing.T) {
	for _, test := range UnknownErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestUnknownErrorTemporary(t *testing.T) {
	for _, test := range UnknownErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestUnknownErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range UnknownErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestUnknownErrorGetCode(t *testing.T) {
	for _, test := range UnknownErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestUnknownErrorGetMessage(t *testing.T) {
	for _, test := range UnknownErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestUnknownErrorGetCause(t *testing.T) {
	for _, test := range UnknownErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestUnknownErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range UnknownErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestUnknownErrorJson(t *testing.T) {
	for _, test := range UnknownErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestUnknownErrorGrpc(t *testing.T) {
	for _, test := range UnknownErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
