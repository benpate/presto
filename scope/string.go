package scope

import (
	"github.com/benpate/data/expression"
	"github.com/benpate/derp"
	"github.com/benpate/presto"
)

// String generates a presto.ScoperFunc using the values provided.  Every context parameter will be compared with an "equals" comparison scope.
func String(params ...string) presto.ScopeFunc {

	return func(ctx presto.Context) (expression.Expression, *derp.Error) {

		result := expression.And()

		for _, param := range params {

			value := ctx.Param(param)

			if value == "" {
				return nil, derp.New(500, "scope.DefaultToken", "Parameter cannot be empty", param)
			}

			result = result.And(param, expression.OperatorEqual, value)
		}

		return result, nil
	}
}
