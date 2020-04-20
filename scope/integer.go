package scope

import (
	"github.com/benpate/data/expression"
	"github.com/benpate/derp"
	"github.com/benpate/presto"
	"github.com/labstack/echo/v4"
)

// Integer generates a presto.ScoperFunc using the values provided.  Every context parameter will be compared with an "equals" comparison scope.
func Integer(values ...string) presto.ScopeFunc {

	return func(ctx echo.Context) (expression.Expression, *derp.Error) {

		token := ctx.Get("token")

		if token == "" {
			return nil, derp.New(500, "scope.DefaultToken", "Token cannot be empty")
		}

		return expression.New("token", expression.OperatorEqual, token), nil
	}
}
