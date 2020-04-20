package scope

import (
	"testing"

	"github.com/benpate/data/expression"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {

	ctx := getTestContext()
	ctx.SetParamNames("myId", "yourId", "hisId", "herId")
	ctx.SetParamValues("1", "2", "3", "4")

	scope := String("myId", "yourId", "hisId", "herId")
	result, err := scope(ctx)

	expected := expression.New("myId", "=", "1").And("yourId", "=", "2").And("hisId", "=", "3").And("herId", "=", "4")

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestStringFail(t *testing.T) {

	ctx := getTestContext()

	ctx.SetParamNames("empty")
	ctx.SetParamValues("")

	scope := String("empty")
	result, err := scope(ctx)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	t.Log(err)
}
