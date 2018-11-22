package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type NotFoundErrorTest struct {
	err          *NotFoundError
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

var NotFoundErrorTests = []NotFoundErrorTest{
	{
		err:        NewNotFoundError("Message 1"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 404: NOT FOUND. Message 1",
		getCode:    404,
		getMessage: "NOT FOUND. Message 1",
		getCause:   nil,
		json:       []byte(`{"error_code":404,"error_message":"NOT FOUND. Message 1"}`),
		rpcCode:    codes.NotFound,
		rpcMessage: "NOT FOUND. Message 1",
	},
	{
		err:        NewNotFoundError("Message 2", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 404: NOT FOUND. Message 2",
		getCode:    404,
		getMessage: "NOT FOUND. Message 2",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"error_code":404,"error_message":"NOT FOUND. Message 2"}`),
		rpcCode:    codes.NotFound,
		rpcMessage: "NOT FOUND. Message 2",
	},
	{
		err:        NewNotFoundError("Message 2", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 404: NOT FOUND. Message 2",
		getCode:    404,
		getMessage: "NOT FOUND. Message 2",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"error_code":404,"error_message":"NOT FOUND. Message 2"}`),
		rpcCode:    codes.NotFound,
		rpcMessage: "NOT FOUND. Message 2",
	},
}

func TestNotFoundErrorTimeout(t *testing.T) {
	for _, test := range NotFoundErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestNotFoundErrorTemporary(t *testing.T) {
	for _, test := range NotFoundErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestNotFoundErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range NotFoundErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestNotFoundErrorGetCode(t *testing.T) {
	for _, test := range NotFoundErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestNotFoundErrorGetMessage(t *testing.T) {
	for _, test := range NotFoundErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestNotFoundErrorGetCause(t *testing.T) {
	for _, test := range NotFoundErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestNotFoundErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range NotFoundErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestNotFoundErrorJson(t *testing.T) {
	for _, test := range NotFoundErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestNotFoundErrorGrpc(t *testing.T) {
	for _, test := range NotFoundErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
