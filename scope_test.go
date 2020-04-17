package presto

/*
func TestScope(t *testing.T) {

	type TestStruct struct {
		Integer int
		String  string
	}

	e := echo.New()
	UseRouter(nil)

	fortyTwo := func(ctx echo.Context) (expression.Expression, *derp.Error) {
		return expression.New("integer", "=", 42), nil
	}

	UseScopes(fortyTwo)

	collection := Collection{}

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
	collection.UseScopes(func(ctx echo.Context) (expression.Expression, *derp.Error) {
		return expression.Expression("string", "=", "hey-oh"), nil
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
*/
