package errors

import (
	"context"
	"testing"

	assert "github.com/stretchr/testify/assert"
	status "google.golang.org/grpc/status"
)

func TestNewPassthroughError(t *testing.T) {
	var tests = []struct {
		msg string
		err error
		exp error
	}{
		{
			msg: "foo",
			err: context.Canceled,
			exp: NewCanceledError("foo", context.Canceled),
		},
		{
			msg: "foo",
			err: NewCanceledError("bar"),
			exp: NewCanceledError("foo", NewCanceledError("bar")),
		},
		{
			msg: "foo",
			err: context.DeadlineExceeded,
			exp: NewDeadlineExceededError("foo", context.DeadlineExceeded),
		},
		{
			msg: "foo",
			err: NewDeadlineExceededError("bar"),
			exp: NewDeadlineExceededError("foo", NewDeadlineExceededError("bar")),
		},
		{
			msg: "foo",
			err: NewUnavailableError("bar"),
			exp: NewUnavailableError("foo", NewUnavailableError("bar")),
		},
		{
			msg: "foo",
			err: NewUnknownError("bar"),
			exp: NewUnknownError("foo", NewUnknownError("bar")),
		},
		{
			msg: "foo",
			err: NewInternalError("bar"),
			exp: NewInternalError("foo", NewInternalError("bar")),
		},
	}

	for _, test := range tests {

		SetVerbosity(Info)

		res := NewPassthroughError(test.msg, test.err)
		resError, _ := res.(interface{ Error() string })
		resTimeout, _ := res.(interface{ Timeout() bool })
		resTemporary, _ := res.(interface{ Temporary() bool })
		resGetCode, _ := res.(interface{ GetCode() int })
		resGetMessage, _ := res.(interface{ GetMessage() string })
		resGetCause, _ := res.(interface{ GetCause() error })
		resGetStack, _ := res.(interface{ GetStack() stack })
		resGRPCStatus, _ := res.(interface{ GRPCStatus() *(status.Status) })

		expError, _ := test.exp.(interface{ Error() string })
		expTimeout, _ := test.exp.(interface{ Timeout() bool })
		expTemporary, _ := test.exp.(interface{ Temporary() bool })
		expGetCode, _ := test.exp.(interface{ GetCode() int })
		expGetMessage, _ := test.exp.(interface{ GetMessage() string })
		expGetCause, _ := test.exp.(interface{ GetCause() error })
		expGRPCStatus, _ := test.exp.(interface{ GRPCStatus() *(status.Status) })

		assert.Equal(t, resError.Error(), expError.Error())
		assert.Equal(t, resTimeout.Timeout(), expTimeout.Timeout())
		assert.Equal(t, resTemporary.Temporary(), expTemporary.Temporary())
		assert.Equal(t, resGetCode.GetCode(), expGetCode.GetCode())
		assert.Equal(t, resGetMessage.GetMessage(), expGetMessage.GetMessage())

		if resGetCause.GetCause() == nil {
			assert.Nil(t, expGetCause.GetCause())
		} else {
			assert.Equal(t, resGetCause.GetCause().Error(), expGetCause.GetCause().Error())
		}

		// trace output prevents testing for string match in different contexts
		assert.NotNil(t, resGetStack.GetStack())

		r := resGRPCStatus.GRPCStatus()
		e := expGRPCStatus.GRPCStatus()
		assert.Equal(t, r.Code(), e.Code())
		assert.Equal(t, r.Message(), e.Message())
	}
}
