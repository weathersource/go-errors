package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type UnauthenticatedErrorTest struct {
	err          *UnauthenticatedError
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

var UnauthenticatedErrorTests = []UnauthenticatedErrorTest{
	{
		err:        NewUnauthenticatedError("Message 1"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 401: UNAUTHENTICATED. Message 1",
		getCode:    401,
		getMessage: "UNAUTHENTICATED. Message 1",
		getCause:   nil,
		json:       []byte(`{"errorCode":401,"errorMessage":"UNAUTHENTICATED. Message 1"}`),
		rpcCode:    codes.Unauthenticated,
		rpcMessage: "UNAUTHENTICATED. Message 1",
	},
	{
		err:        NewUnauthenticatedError("Message 2", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 401: UNAUTHENTICATED. Message 2",
		getCode:    401,
		getMessage: "UNAUTHENTICATED. Message 2",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"errorCode":401,"errorMessage":"UNAUTHENTICATED. Message 2"}`),
		rpcCode:    codes.Unauthenticated,
		rpcMessage: "UNAUTHENTICATED. Message 2",
	},
	{
		err:        NewUnauthenticatedError("Message 2", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 401: UNAUTHENTICATED. Message 2",
		getCode:    401,
		getMessage: "UNAUTHENTICATED. Message 2",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"errorCode":401,"errorMessage":"UNAUTHENTICATED. Message 2"}`),
		rpcCode:    codes.Unauthenticated,
		rpcMessage: "UNAUTHENTICATED. Message 2",
	},
}

func TestUnauthenticatedErrorTimeout(t *testing.T) {
	for _, test := range UnauthenticatedErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestUnauthenticatedErrorTemporary(t *testing.T) {
	for _, test := range UnauthenticatedErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestUnauthenticatedErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range UnauthenticatedErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestUnauthenticatedErrorGetCode(t *testing.T) {
	for _, test := range UnauthenticatedErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestUnauthenticatedErrorGetMessage(t *testing.T) {
	for _, test := range UnauthenticatedErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestUnauthenticatedErrorGetCause(t *testing.T) {
	for _, test := range UnauthenticatedErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestUnauthenticatedErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range UnauthenticatedErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestUnauthenticatedErrorJson(t *testing.T) {
	for _, test := range UnauthenticatedErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestUnauthenticatedErrorGrpc(t *testing.T) {
	for _, test := range UnauthenticatedErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
