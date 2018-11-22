package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type InvalidArgumentErrorTest struct {
	err          *InvalidArgumentError
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

var InvalidArgumentErrorTests = []InvalidArgumentErrorTest{
	{
		err:        NewInvalidArgumentError("Message 1"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 400: INVALID ARGUMENT. Message 1",
		getCode:    400,
		getMessage: "INVALID ARGUMENT. Message 1",
		getCause:   nil,
		json:       []byte(`{"error_code":400,"error_message":"INVALID ARGUMENT. Message 1"}`),
		rpcCode:    codes.InvalidArgument,
		rpcMessage: "INVALID ARGUMENT. Message 1",
	},
	{
		err:        NewInvalidArgumentError("Message 2", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 400: INVALID ARGUMENT. Message 2",
		getCode:    400,
		getMessage: "INVALID ARGUMENT. Message 2",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"error_code":400,"error_message":"INVALID ARGUMENT. Message 2"}`),
		rpcCode:    codes.InvalidArgument,
		rpcMessage: "INVALID ARGUMENT. Message 2",
	},
	{
		err:        NewInvalidArgumentError("Message 2", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 400: INVALID ARGUMENT. Message 2",
		getCode:    400,
		getMessage: "INVALID ARGUMENT. Message 2",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"error_code":400,"error_message":"INVALID ARGUMENT. Message 2"}`),
		rpcCode:    codes.InvalidArgument,
		rpcMessage: "INVALID ARGUMENT. Message 2",
	},
}

func TestInvalidArgumentErrorTimeout(t *testing.T) {
	for _, test := range InvalidArgumentErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestInvalidArgumentErrorTemporary(t *testing.T) {
	for _, test := range InvalidArgumentErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestInvalidArgumentErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range InvalidArgumentErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestInvalidArgumentErrorGetCode(t *testing.T) {
	for _, test := range InvalidArgumentErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestInvalidArgumentErrorGetMessage(t *testing.T) {
	for _, test := range InvalidArgumentErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestInvalidArgumentErrorGetCause(t *testing.T) {
	for _, test := range InvalidArgumentErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestInvalidArgumentErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range InvalidArgumentErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestInvalidArgumentErrorJson(t *testing.T) {
	for _, test := range InvalidArgumentErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestInvalidArgumentErrorGrpc(t *testing.T) {
	for _, test := range InvalidArgumentErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
