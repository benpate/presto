package scope

import (
	"testing"

	"github.com/benpate/data/expression"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestObjectID(t *testing.T) {

	ID1, _ := primitive.ObjectIDFromHex("123456789012345678901234")
	ID2, _ := primitive.ObjectIDFromHex("000000000000000000000000")

	ctx := getTestContext()
	ctx.SetParamNames("first", "second")
	ctx.SetParamValues(ID1.Hex(), ID2.Hex())

	scope := ObjectID("first", "second")

	result, err := scope(ctx)

	expected := expression.New("first", "=", ID1).And("second", "=", ID2)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestObjectIDFail(t *testing.T) {

	for _, value := range []string{"", "bad-value", "0000000000"} {

		ctx := getTestContext()
		ctx.SetParamNames("error")
		ctx.SetParamValues(value)

		scope := ObjectID("error")

		result, err := scope(ctx)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	}
}
