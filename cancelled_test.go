package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type CancelledErrorTest struct {
	err          *CancelledError
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

var CancelledErrorTests = []CancelledErrorTest{
	{
		err:        NewCancelledError("foo"),
		timeout:    true,
		temporary:  false,
		errorInfo:  "error 499: CANCELLED. Request cancelled by the client. foo",
		getCode:    499,
		getMessage: "CANCELLED. Request cancelled by the client. foo",
		getCause:   nil,
		json:       []byte(`{"error_code":499,"error_message":"CANCELLED. Request cancelled by the client."}`),
		rpcCode:    codes.Canceled,
		rpcMessage: "CANCELLED. Request cancelled by the client.",
	},
	{
		err:        NewCancelledError("foo", errors.New("causal error")),
		timeout:    true,
		temporary:  false,
		errorInfo:  "error 499: CANCELLED. Request cancelled by the client. foo",
		getCode:    499,
		getMessage: "CANCELLED. Request cancelled by the client. foo",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"error_code":499,"error_message":"CANCELLED. Request cancelled by the client."}`),
		rpcCode:    codes.Canceled,
		rpcMessage: "CANCELLED. Request cancelled by the client.",
	},
	{
		err:        NewCancelledError("foo", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    true,
		temporary:  false,
		errorInfo:  "error 499: CANCELLED. Request cancelled by the client. foo",
		getCode:    499,
		getMessage: "CANCELLED. Request cancelled by the client. foo",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"error_code":499,"error_message":"CANCELLED. Request cancelled by the client."}`),
		rpcCode:    codes.Canceled,
		rpcMessage: "CANCELLED. Request cancelled by the client.",
	},
}

func TestCancelledErrorTimeout(t *testing.T) {
	for _, test := range CancelledErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestCancelledErrorTemporary(t *testing.T) {
	for _, test := range CancelledErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestCancelledErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range CancelledErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestCancelledErrorGetCode(t *testing.T) {
	for _, test := range CancelledErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestCancelledErrorGetMessage(t *testing.T) {
	for _, test := range CancelledErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestCancelledErrorGetCause(t *testing.T) {
	for _, test := range CancelledErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestCancelledErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range CancelledErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestCancelledErrorJson(t *testing.T) {
	for _, test := range CancelledErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestCancelledErrorGrpc(t *testing.T) {
	for _, test := range CancelledErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
