package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type PermissionDeniedErrorTest struct {
	err          *PermissionDeniedError
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

var PermissionDeniedErrorTests = []PermissionDeniedErrorTest{
	{
		err:        NewPermissionDeniedError("Message 1"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 403: PERMISSION DENIED. Message 1",
		getCode:    403,
		getMessage: "PERMISSION DENIED. Message 1",
		getCause:   nil,
		json:       []byte(`{"error_code":403,"error_message":"PERMISSION DENIED. Message 1"}`),
		rpcCode:    codes.PermissionDenied,
		rpcMessage: "PERMISSION DENIED. Message 1",
	},
	{
		err:        NewPermissionDeniedError("Message 2", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 403: PERMISSION DENIED. Message 2",
		getCode:    403,
		getMessage: "PERMISSION DENIED. Message 2",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"error_code":403,"error_message":"PERMISSION DENIED. Message 2"}`),
		rpcCode:    codes.PermissionDenied,
		rpcMessage: "PERMISSION DENIED. Message 2",
	},
	{
		err:        NewPermissionDeniedError("Message 2", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 403: PERMISSION DENIED. Message 2",
		getCode:    403,
		getMessage: "PERMISSION DENIED. Message 2",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"error_code":403,"error_message":"PERMISSION DENIED. Message 2"}`),
		rpcCode:    codes.PermissionDenied,
		rpcMessage: "PERMISSION DENIED. Message 2",
	},
}

func TestPermissionDeniedErrorTimeout(t *testing.T) {
	for _, test := range PermissionDeniedErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestPermissionDeniedErrorTemporary(t *testing.T) {
	for _, test := range PermissionDeniedErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestPermissionDeniedErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range PermissionDeniedErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestPermissionDeniedErrorGetCode(t *testing.T) {
	for _, test := range PermissionDeniedErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestPermissionDeniedErrorGetMessage(t *testing.T) {
	for _, test := range PermissionDeniedErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestPermissionDeniedErrorGetCause(t *testing.T) {
	for _, test := range PermissionDeniedErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestPermissionDeniedErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range PermissionDeniedErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestPermissionDeniedErrorJson(t *testing.T) {
	for _, test := range PermissionDeniedErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestPermissionDeniedErrorGrpc(t *testing.T) {
	for _, test := range PermissionDeniedErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
