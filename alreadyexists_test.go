package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type AlreadyExistsErrorTest struct {
	err          *AlreadyExistsError
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

var AlreadyExistsErrorTests = []AlreadyExistsErrorTest{
	{
		err:        NewAlreadyExistsError("Message 1"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 409: ALREADY EXISTS. Message 1",
		getCode:    409,
		getMessage: "ALREADY EXISTS. Message 1",
		getCause:   nil,
		json:       []byte(`{"error_code":409,"error_message":"ALREADY EXISTS. Message 1"}`),
		rpcCode:    codes.AlreadyExists,
		rpcMessage: "ALREADY EXISTS. Message 1",
	},
	{
		err:        NewAlreadyExistsError("Message 2", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 409: ALREADY EXISTS. Message 2",
		getCode:    409,
		getMessage: "ALREADY EXISTS. Message 2",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"error_code":409,"error_message":"ALREADY EXISTS. Message 2"}`),
		rpcCode:    codes.AlreadyExists,
		rpcMessage: "ALREADY EXISTS. Message 2",
	},
	{
		err:        NewAlreadyExistsError("Message 2", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 409: ALREADY EXISTS. Message 2",
		getCode:    409,
		getMessage: "ALREADY EXISTS. Message 2",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"error_code":409,"error_message":"ALREADY EXISTS. Message 2"}`),
		rpcCode:    codes.AlreadyExists,
		rpcMessage: "ALREADY EXISTS. Message 2",
	},
}

func TestAlreadyExistsErrorTimeout(t *testing.T) {
	for _, test := range AlreadyExistsErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestAlreadyExistsErrorTemporary(t *testing.T) {
	for _, test := range AlreadyExistsErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestAlreadyExistsErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range AlreadyExistsErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestAlreadyExistsErrorGetCode(t *testing.T) {
	for _, test := range AlreadyExistsErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestAlreadyExistsErrorGetMessage(t *testing.T) {
	for _, test := range AlreadyExistsErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestAlreadyExistsErrorGetCause(t *testing.T) {
	for _, test := range AlreadyExistsErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestAlreadyExistsErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range AlreadyExistsErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestAlreadyExistsErrorJson(t *testing.T) {
	for _, test := range AlreadyExistsErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestAlreadyExistsErrorGrpc(t *testing.T) {
	for _, test := range AlreadyExistsErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
