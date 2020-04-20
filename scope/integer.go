package scope

import (
	"strconv"

	"github.com/benpate/data/expression"
	"github.com/benpate/derp"
	"github.com/benpate/presto"
	"github.com/labstack/echo/v4"
)

// Integer generates a presto.ScoperFunc using the values provided.  Every context parameter will be compared with an "equals" comparison scope.
func Integer(params ...string) presto.ScopeFunc {

	return func(ctx echo.Context) (expression.Expression, *derp.Error) {

		result := expression.And()

		for _, param := range params {

			value := ctx.Param(param)

			integer, err := strconv.Atoi(value)

			if err != nil {
				return nil, derp.New(500, "scope.Integer", "Invalid parameter", param, value, err)
			}

			result = result.And(param, expression.OperatorEqual, integer)
		}

		return result, nil
	}
}
