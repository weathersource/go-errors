package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type CanceledErrorTest struct {
	err          *CanceledError
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

var CanceledErrorTests = []CanceledErrorTest{
	{
		err:        NewCanceledError("foo"),
		timeout:    true,
		temporary:  false,
		errorInfo:  "error 499: CANCELED. Request canceled by the client. foo",
		getCode:    499,
		getMessage: "CANCELED. Request canceled by the client. foo",
		getCause:   nil,
		json:       []byte(`{"errorCode":499,"errorMessage":"CANCELED. Request canceled by the client."}`),
		rpcCode:    codes.Canceled,
		rpcMessage: "CANCELED. Request canceled by the client.",
	},
	{
		err:        NewCanceledError("foo", errors.New("causal error")),
		timeout:    true,
		temporary:  false,
		errorInfo:  "error 499: CANCELED. Request canceled by the client. foo",
		getCode:    499,
		getMessage: "CANCELED. Request canceled by the client. foo",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"errorCode":499,"errorMessage":"CANCELED. Request canceled by the client."}`),
		rpcCode:    codes.Canceled,
		rpcMessage: "CANCELED. Request canceled by the client.",
	},
	{
		err:        NewCanceledError("foo", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    true,
		temporary:  false,
		errorInfo:  "error 499: CANCELED. Request canceled by the client. foo",
		getCode:    499,
		getMessage: "CANCELED. Request canceled by the client. foo",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"errorCode":499,"errorMessage":"CANCELED. Request canceled by the client."}`),
		rpcCode:    codes.Canceled,
		rpcMessage: "CANCELED. Request canceled by the client.",
	},
}

func TestCanceledErrorTimeout(t *testing.T) {
	for _, test := range CanceledErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestCanceledErrorTemporary(t *testing.T) {
	for _, test := range CanceledErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestCanceledErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range CanceledErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestCanceledErrorGetCode(t *testing.T) {
	for _, test := range CanceledErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestCanceledErrorGetMessage(t *testing.T) {
	for _, test := range CanceledErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestCanceledErrorGetCause(t *testing.T) {
	for _, test := range CanceledErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestCanceledErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range CanceledErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestCanceledErrorJson(t *testing.T) {
	for _, test := range CanceledErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestCanceledErrorGrpc(t *testing.T) {
	for _, test := range CanceledErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}

func TestCanceledErrorAppend(t *testing.T) {

	e1       := NewCanceledError("Message 1")
	e1append := e1.Append(errors.New("foo"))
	e1alt    := NewCanceledError("Message 1", errors.New("foo"))
	assert.Equal(t, e1alt.GetCause().Error(), e1append.GetCause().Error())

	e2       :=NewCanceledError("Message 2", errors.New("foo"))
	e2append := e2.Append(errors.New("bar"))
	e2alt    := NewCanceledError("Message 2", errors.New("foo"), errors.New("bar"))
	assert.Equal(t, e2alt.GetCause().Error(), e2append.GetCause().Error())
}
