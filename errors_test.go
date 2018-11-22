package errors

import (
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

type errorsTest struct {
	err Errors
	str string
}

var errorsTests = []errorsTest{
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
		str: "Errors:\n#1: foo\n#2: bar",
	},
}

func TestErrorsError(t *testing.T) {
	for _, test := range errorsTests {
		assert.Equal(t, test.str, test.err.Error())
	}
}
