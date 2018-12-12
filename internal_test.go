package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type InternalErrorTest struct {
	err          *InternalError
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

var InternalErrorTests = []InternalErrorTest{
	{
		err:        NewInternalError("foo"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 500: INTERNAL ERROR. foo",
		getCode:    500,
		getMessage: "INTERNAL ERROR. foo",
		getCause:   nil,
		json:       []byte(`{"errorCode":500,"errorMessage":"INTERNAL ERROR."}`),
		rpcCode:    codes.Internal,
		rpcMessage: "INTERNAL ERROR.",
	},
	{
		err:        NewInternalError("foo", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 500: INTERNAL ERROR. foo",
		getCode:    500,
		getMessage: "INTERNAL ERROR. foo",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"errorCode":500,"errorMessage":"INTERNAL ERROR."}`),
		rpcCode:    codes.Internal,
		rpcMessage: "INTERNAL ERROR.",
	},
	{
		err:        NewInternalError("foo", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 500: INTERNAL ERROR. foo",
		getCode:    500,
		getMessage: "INTERNAL ERROR. foo",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"errorCode":500,"errorMessage":"INTERNAL ERROR."}`),
		rpcCode:    codes.Internal,
		rpcMessage: "INTERNAL ERROR.",
	},
}

func TestInternalErrorTimeout(t *testing.T) {
	for _, test := range InternalErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestInternalErrorTemporary(t *testing.T) {
	for _, test := range InternalErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestInternalErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range InternalErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestInternalErrorGetCode(t *testing.T) {
	for _, test := range InternalErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestInternalErrorGetMessage(t *testing.T) {
	for _, test := range InternalErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestInternalErrorGetCause(t *testing.T) {
	for _, test := range InternalErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestInternalErrorGetStack(t *testing.T) {
	for _, test := range InternalErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestInternalErrorJson(t *testing.T) {
	for _, test := range InternalErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestInternalErrorGrpc(t *testing.T) {
	for _, test := range InternalErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
