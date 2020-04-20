package scope

import (
	"testing"

	"github.com/benpate/data/expression"
	"github.com/stretchr/testify/assert"
)

func TestInteger(t *testing.T) {

	ctx := getTestContext()
	ctx.SetParamNames("first", "second")
	ctx.SetParamValues("1", "2")

	scope := Integer("first", "second")

	result, err := scope(ctx)

	expected := expression.New("first", "=", 1).And("second", "=", 2)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestIntegerFail(t *testing.T) {

	for _, value := range []string{"", "bad-value"} {

		ctx := getTestContext()
		ctx.SetParamNames("error")
		ctx.SetParamValues(value)

		scope := Integer("error")

		result, err := scope(ctx)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	}
}
