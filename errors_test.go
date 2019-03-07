package errors

import (
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func TestErrorsError(t *testing.T) {
	tests := []struct {
		err *Errors
		str string
	}{
		{
			err: NewErrors(),
			str: "",
		},
		{
			err: NewErrors(
				errors.New("foo"),
			),
			str: "foo",
		},
		{
			err: NewErrors(
				errors.New("foo"),
				errors.New("bar"),
			),
			str: "MULTIPLE ERRORS.\n#1: foo\n#2: bar",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.str, test.err.Error())
	}
}

func TestErrorsShift(t *testing.T) {
	tests := []struct {
		errs *Errors
		text string
		cnt  int
	}{
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			"foo",
			1,
		},
	}
	for _, test := range tests {
		err := test.errs.Shift()
		assert.Equal(t, test.text, err.Error())
		assert.Equal(t, test.cnt, test.errs.Len())
	}
	testNils := []struct {
		errs *Errors
		cnt  int
	}{
		{
			NewErrors(),
			0,
		},
	}
	for _, test := range testNils {
		err := test.errs.Shift()
		assert.Nil(t, err)
		assert.Equal(t, test.cnt, test.errs.Len())
	}
}

func TestErrorsAppend(t *testing.T) {
	tests := []struct {
		errs *Errors
		err  error
		cnt  int
	}{
		{
			NewErrors(),
			nil,
			0,
		},
		{
			NewErrors(errors.New("foo")),
			errors.New("bar"),
			2,
		},
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			errors.New("baz"),
			3,
		},
	}
	for _, test := range tests {
		test.errs.Append(test.err)
		assert.Equal(t, test.cnt, test.errs.Len())
	}
}

func TestErrorsPop(t *testing.T) {
	tests := []struct {
		errs *Errors
		text string
		cnt  int
	}{
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			"bar",
			1,
		},
	}
	for _, test := range tests {
		err := test.errs.Pop()
		assert.Equal(t, test.text, err.Error())
		assert.Equal(t, test.cnt, test.errs.Len())
	}
	testNils := []struct {
		errs *Errors
		cnt  int
	}{
		{
			NewErrors(),
			0,
		},
	}
	for _, test := range testNils {
		err := test.errs.Pop()
		assert.Nil(t, err)
		assert.Equal(t, test.cnt, test.errs.Len())
	}
}

func TestErrorsLen(t *testing.T) {
	tests := []struct {
		errs *Errors
		cnt  int
	}{
		{
			nil,
			0,
		},
		{
			NewErrors(),
			0,
		},
		{
			NewErrors(nil),
			0,
		},
		{
			NewErrors(errors.New("foo")),
			1,
		},
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			2,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.cnt, test.errs.Len())
	}
}

func TestErrorsTimeout(t *testing.T) {
	tests := []struct {
		errs    *Errors
		timeout bool
	}{
		{
			NewErrors(NewCancelledError("foo")),
			true,
		},
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			false,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.timeout, test.errs.Timeout())
	}
}

func TestErrorsTemporary(t *testing.T) {
	tests := []struct {
		errs      *Errors
		temporary bool
	}{
		{
			NewErrors(NewUnavailableError("foo")),
			true,
		},
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			false,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.temporary, test.errs.Temporary())
	}
}

func TestErrorsGetCode(t *testing.T) {
	tests := []struct {
		errs *Errors
		code int
	}{
		{
			nil,
			200,
		},
		{
			NewErrors(),
			200,
		},
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			500,
		},
		{
			NewErrors(NewUnavailableError("foo")),
			503,
		},
		{
			NewErrors(status.Error(codes.Aborted, "foo bar")),
			409,
		},
		{
			NewErrors(status.Error(codes.AlreadyExists, "foo bar")),
			409,
		},
		{
			NewErrors(status.Error(codes.Canceled, "foo bar")),
			499,
		},
		{
			NewErrors(status.Error(codes.DataLoss, "foo bar")),
			500,
		},
		{
			NewErrors(status.Error(codes.DeadlineExceeded, "foo bar")),
			504,
		},
		{
			NewErrors(status.Error(codes.FailedPrecondition, "foo bar")),
			400,
		},
		{
			NewErrors(status.Error(codes.Internal, "foo bar")),
			500,
		},
		{
			NewErrors(status.Error(codes.InvalidArgument, "foo bar")),
			400,
		},
		{
			NewErrors(status.Error(codes.NotFound, "foo bar")),
			404,
		},
		{
			NewErrors(status.Error(codes.OutOfRange, "foo bar")),
			400,
		},
		{
			NewErrors(status.Error(codes.PermissionDenied, "foo bar")),
			403,
		},
		{
			NewErrors(status.Error(codes.ResourceExhausted, "foo bar")),
			429,
		},
		{
			NewErrors(status.Error(codes.Unauthenticated, "foo bar")),
			401,
		},
		{
			NewErrors(status.Error(codes.Unavailable, "foo bar")),
			503,
		},
		{
			NewErrors(status.Error(codes.Unimplemented, "foo bar")),
			501,
		},
		{
			NewErrors(status.Error(codes.Unknown, "foo bar")),
			500,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.code, test.errs.GetCode())
	}
}

func TestErrorsGetMessage(t *testing.T) {
	tests := []struct {
		errs    *Errors
		message string
	}{
		{
			nil,
			"",
		},
		{
			NewErrors(),
			"",
		},
		{
			NewErrors(errors.New("foo")),
			"foo",
		},
		{
			NewErrors(NewInternalError("foo")),
			"INTERNAL ERROR. foo",
		},
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			"MULTIPLE ERRORS.",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.message, test.errs.GetMessage())
	}
}

func TestErrorsGetCause(t *testing.T) {
	tests := []struct {
		errs  *Errors
		isnil bool
	}{
		{
			nil,
			true,
		},
		{
			NewErrors(),
			true,
		},
		{
			NewErrors(errors.New("foo")),
			true,
		},
		{
			NewErrors(NewInternalError("foo", errors.New(""))),
			false,
		},
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			false,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.isnil, test.errs.GetCause() == nil)
	}
}

func TestErrorsGetStack(t *testing.T) {
	tests := []struct {
		errs   *Errors
		isZero bool
	}{
		{
			nil,
			true,
		},
		{
			NewErrors(),
			true,
		},
		{
			NewErrors(errors.New("foo")),
			true,
		},
		{
			NewErrors(NewInternalError("foo", errors.New(""))),
			false,
		},
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			true,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.isZero, 0 == len(test.errs.GetStack()))
	}
}

func TestErrorsGRPCStatus(t *testing.T) {
	tests := []struct {
		errs    *Errors
		message string
		code    codes.Code
	}{
		{
			NewErrors(errors.New("foo")),
			"foo",
			codes.Unknown,
		},
		{
			NewErrors(NewNotFoundError("foo")),
			"NOT FOUND. foo",
			codes.NotFound,
		},
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			"MULTIPLE ERRORS.\n#1: foo\n#2: bar",
			codes.Unknown,
		},
	}
	for _, test := range tests {
		status := test.errs.GRPCStatus()
		assert.NotNil(t, status)
		assert.Equal(t, test.message, status.Message())
		assert.Equal(t, test.code, status.Code())
	}
	testNils := []struct {
		errs *Errors
	}{
		{nil},
		{NewErrors()},
	}
	for _, test := range testNils {
		assert.Nil(t, test.errs.GRPCStatus())
	}
}
func TestErrorsPeekLocked(t *testing.T) {
	tests := []struct {
		errs *Errors
		err  error
	}{
		{
			NewErrors(),
			nil,
		},
		{
			NewErrors(errors.New("foo")),
			errors.New("foo"),
		},
		{
			NewErrors(errors.New("foo"), errors.New("bar")),
			errors.New("foo"),
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.err, test.errs.peekLocked())
	}
}
