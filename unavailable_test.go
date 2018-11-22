package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type UnavailableErrorTest struct {
	err          *UnavailableError
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

var UnavailableErrorTests = []UnavailableErrorTest{
	{
		err:        NewUnavailableError("foo"),
		timeout:    false,
		temporary:  true,
		errorInfo:  "error 503: UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance. foo",
		getCode:    503,
		getMessage: "UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance. foo",
		getCause:   nil,
		json:       []byte(`{"error_code":503,"error_message":"UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance."}`),
		rpcCode:    codes.Unavailable,
		rpcMessage: "UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance.",
	},
	{
		err:        NewUnavailableError("foo", errors.New("causal error")),
		timeout:    false,
		temporary:  true,
		errorInfo:  "error 503: UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance. foo",
		getCode:    503,
		getMessage: "UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance. foo",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"error_code":503,"error_message":"UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance."}`),
		rpcCode:    codes.Unavailable,
		rpcMessage: "UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance.",
	},
	{
		err:        NewUnavailableError("foo", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  true,
		errorInfo:  "error 503: UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance. foo",
		getCode:    503,
		getMessage: "UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance. foo",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"error_code":503,"error_message":"UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance."}`),
		rpcCode:    codes.Unavailable,
		rpcMessage: "UNAVAILABLE. Unable to handle the request due to a temporary overloading or maintenance.",
	},
}

func TestUnavailableErrorTimeout(t *testing.T) {
	for _, test := range UnavailableErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestUnavailableErrorTemporary(t *testing.T) {
	for _, test := range UnavailableErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestUnavailableErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range UnavailableErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestUnavailableErrorGetCode(t *testing.T) {
	for _, test := range UnavailableErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestUnavailableErrorGetMessage(t *testing.T) {
	for _, test := range UnavailableErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestUnavailableErrorGetCause(t *testing.T) {
	for _, test := range UnavailableErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestUnavailableErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range UnavailableErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestUnavailableErrorJson(t *testing.T) {
	for _, test := range UnavailableErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestUnavailableErrorGrpc(t *testing.T) {
	for _, test := range UnavailableErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
