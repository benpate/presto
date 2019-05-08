package presto

import (
	"testing"

	"github.com/benpate/data"
	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestScope(t *testing.T) {

	type TestStruct struct {
		Integer int
		String  string
	}

	WithScopes(func(ctx echo.Context) (data.Expression, *derp.Error) {
		return data.Expression{{"integer", "=", 42}}, nil
	})

	collection := Collection{}

	e := echo.New()

	ctx := e.NewContext(nil, nil)

	expression, err := collection.getScope(ctx)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// Test global scopes
	record := TestStruct{
		Integer: 42,
		String:  "hoo-ah",
	}

	assert.True(t, expression.Match(record))

	record.Integer = 43
	assert.False(t, expression.Match(record))

	// Test scopes attached to the collection
	collection.WithScopes(func(ctx echo.Context) (data.Expression, *derp.Error) {
		return data.Expression{{"string", "=", "hey-oh"}}, nil
	})

	expression, err = collection.getScope(ctx)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	assert.False(t, expression.Match(record))

	record.String = "hey-oh"
	assert.False(t, expression.Match(record))

	record.Integer = 42
	assert.True(t, expression.Match(record))
}
