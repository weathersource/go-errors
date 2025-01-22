package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type DeadlineExceededErrorTest struct {
	err          *DeadlineExceededError
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

var DeadlineExceededErrorTests = []DeadlineExceededErrorTest{
	{
		err:        NewDeadlineExceededError("foo"),
		timeout:    true,
		temporary:  false,
		errorInfo:  "error 504: DEADLINE EXCEEDED. Server timeout. foo",
		getCode:    504,
		getMessage: "DEADLINE EXCEEDED. Server timeout. foo",
		getCause:   nil,
		json:       []byte(`{"errorCode":504,"errorMessage":"DEADLINE EXCEEDED. Server timeout."}`),
		rpcCode:    codes.DeadlineExceeded,
		rpcMessage: "DEADLINE EXCEEDED. Server timeout.",
	},
	{
		err:        NewDeadlineExceededError("foo", errors.New("causal error")),
		timeout:    true,
		temporary:  false,
		errorInfo:  "error 504: DEADLINE EXCEEDED. Server timeout. foo",
		getCode:    504,
		getMessage: "DEADLINE EXCEEDED. Server timeout. foo",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"errorCode":504,"errorMessage":"DEADLINE EXCEEDED. Server timeout."}`),
		rpcCode:    codes.DeadlineExceeded,
		rpcMessage: "DEADLINE EXCEEDED. Server timeout.",
	},
	{
		err:        NewDeadlineExceededError("foo", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    true,
		temporary:  false,
		errorInfo:  "error 504: DEADLINE EXCEEDED. Server timeout. foo",
		getCode:    504,
		getMessage: "DEADLINE EXCEEDED. Server timeout. foo",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"errorCode":504,"errorMessage":"DEADLINE EXCEEDED. Server timeout."}`),
		rpcCode:    codes.DeadlineExceeded,
		rpcMessage: "DEADLINE EXCEEDED. Server timeout.",
	},
}

func TestDeadlineExceededErrorTimeout(t *testing.T) {
	for _, test := range DeadlineExceededErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestDeadlineExceededErrorTemporary(t *testing.T) {
	for _, test := range DeadlineExceededErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestDeadlineExceededErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range DeadlineExceededErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestDeadlineExceededErrorGetCode(t *testing.T) {
	for _, test := range DeadlineExceededErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestDeadlineExceededErrorGetMessage(t *testing.T) {
	for _, test := range DeadlineExceededErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestDeadlineExceededErrorGetCause(t *testing.T) {
	for _, test := range DeadlineExceededErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestDeadlineExceededErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range DeadlineExceededErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestDeadlineExceededErrorJson(t *testing.T) {
	for _, test := range DeadlineExceededErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestDeadlineExceededErrorGrpc(t *testing.T) {
	for _, test := range DeadlineExceededErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
