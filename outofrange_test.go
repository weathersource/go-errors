package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type OutOfRangeErrorTest struct {
	err          *OutOfRangeError
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

var OutOfRangeErrorTests = []OutOfRangeErrorTest{
	{
		err:        NewOutOfRangeError("Message 1"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 400: OUT OF RANGE. Message 1",
		getCode:    400,
		getMessage: "OUT OF RANGE. Message 1",
		getCause:   nil,
		json:       []byte(`{"errorCode":400,"errorMessage":"OUT OF RANGE. Message 1"}`),
		rpcCode:    codes.OutOfRange,
		rpcMessage: "OUT OF RANGE. Message 1",
	},
	{
		err:        NewOutOfRangeError("Message 2", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 400: OUT OF RANGE. Message 2",
		getCode:    400,
		getMessage: "OUT OF RANGE. Message 2",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"errorCode":400,"errorMessage":"OUT OF RANGE. Message 2"}`),
		rpcCode:    codes.OutOfRange,
		rpcMessage: "OUT OF RANGE. Message 2",
	},
	{
		err:        NewOutOfRangeError("Message 2", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 400: OUT OF RANGE. Message 2",
		getCode:    400,
		getMessage: "OUT OF RANGE. Message 2",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"errorCode":400,"errorMessage":"OUT OF RANGE. Message 2"}`),
		rpcCode:    codes.OutOfRange,
		rpcMessage: "OUT OF RANGE. Message 2",
	},
}

func TestOutOfRangeErrorTimeout(t *testing.T) {
	for _, test := range OutOfRangeErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestOutOfRangeErrorTemporary(t *testing.T) {
	for _, test := range OutOfRangeErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestOutOfRangeErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range OutOfRangeErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestOutOfRangeErrorGetCode(t *testing.T) {
	for _, test := range OutOfRangeErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestOutOfRangeErrorGetMessage(t *testing.T) {
	for _, test := range OutOfRangeErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestOutOfRangeErrorGetCause(t *testing.T) {
	for _, test := range OutOfRangeErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestOutOfRangeErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range OutOfRangeErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestOutOfRangeErrorJson(t *testing.T) {
	for _, test := range OutOfRangeErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestOutOfRangeErrorGrpc(t *testing.T) {
	for _, test := range OutOfRangeErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
