package errors

import (
	"encoding/json"
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
)

type DataLossErrorTest struct {
	err          *DataLossError
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

var DataLossErrorTests = []DataLossErrorTest{
	{
		err:        NewDataLossError("foo"),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 500: DATA LOSS. Unrecoverable data loss or data corruption. foo",
		getCode:    500,
		getMessage: "DATA LOSS. Unrecoverable data loss or data corruption. foo",
		getCause:   nil,
		json:       []byte(`{"error_code":500,"error_message":"DATA LOSS. Unrecoverable data loss or data corruption."}`),
		rpcCode:    codes.DataLoss,
		rpcMessage: "DATA LOSS. Unrecoverable data loss or data corruption.",
	},
	{
		err:        NewDataLossError("foo", errors.New("causal error")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 500: DATA LOSS. Unrecoverable data loss or data corruption. foo",
		getCode:    500,
		getMessage: "DATA LOSS. Unrecoverable data loss or data corruption. foo",
		getCause:   errors.New("causal error"),
		json:       []byte(`{"error_code":500,"error_message":"DATA LOSS. Unrecoverable data loss or data corruption."}`),
		rpcCode:    codes.DataLoss,
		rpcMessage: "DATA LOSS. Unrecoverable data loss or data corruption.",
	},
	{
		err:        NewDataLossError("foo", errors.New("causal error"), errors.New("causal error 2")),
		timeout:    false,
		temporary:  false,
		errorInfo:  "error 500: DATA LOSS. Unrecoverable data loss or data corruption. foo",
		getCode:    500,
		getMessage: "DATA LOSS. Unrecoverable data loss or data corruption. foo",
		getCause:   NewErrors(errors.New("causal error"), errors.New("causal error 2")),
		json:       []byte(`{"error_code":500,"error_message":"DATA LOSS. Unrecoverable data loss or data corruption."}`),
		rpcCode:    codes.DataLoss,
		rpcMessage: "DATA LOSS. Unrecoverable data loss or data corruption.",
	},
}

func TestDataLossErrorTimeout(t *testing.T) {
	for _, test := range DataLossErrorTests {
		assert.Equal(t, test.timeout, test.err.Timeout())
	}
}

func TestDataLossErrorTemporary(t *testing.T) {
	for _, test := range DataLossErrorTests {
		assert.Equal(t, test.temporary, test.err.Temporary())
	}
}

func TestDataLossErrorError(t *testing.T) {
	// note, all other verbosity states tested in error_test.go
	SetVerbosity(Info)
	for _, test := range DataLossErrorTests {
		assert.Equal(t, test.errorInfo, test.err.Error())
	}
}

func TestDataLossErrorGetCode(t *testing.T) {
	for _, test := range DataLossErrorTests {
		assert.Equal(t, test.getCode, test.err.GetCode())
	}
}

func TestDataLossErrorGetMessage(t *testing.T) {
	for _, test := range DataLossErrorTests {
		assert.Equal(t, test.getMessage, test.err.GetMessage())
	}
}

func TestDataLossErrorGetCause(t *testing.T) {
	for _, test := range DataLossErrorTests {
		if test.getCause == nil {
			assert.Nil(t, test.err.GetCause())
		} else {
			assert.Equal(t, test.getCause.Error(), test.err.GetCause().Error())
		}
	}
}

func TestDataLossErrorGetStack(t *testing.T) {
	// trace output prevents testing for string match in different contexts
	for _, test := range DataLossErrorTests {
		assert.NotNil(t, test.err.GetStack())
	}
}

func TestDataLossErrorJson(t *testing.T) {
	for _, test := range DataLossErrorTests {
		json, _ := json.Marshal(test.err)
		assert.Equal(t, string(json[:]), string((test.json)[:]))
	}
}

func TestDataLossErrorGrpc(t *testing.T) {
	for _, test := range DataLossErrorTests {
		s := test.err.GRPCStatus()
		assert.Equal(t, test.rpcCode, s.Code())
		assert.Equal(t, test.rpcMessage, s.Message())
	}
}
