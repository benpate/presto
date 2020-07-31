package scope

import (
	"github.com/benpate/data/expression"
	"github.com/benpate/derp"
	"github.com/benpate/presto"
)

// And generates a presto.ScopeFunc that concatenates a set of other ScopeFunc's into a single expression
func And(params ...presto.ScopeFunc) presto.ScopeFunc {

	return func(ctx presto.Context) (expression.Expression, *derp.Error) {

		result := expression.And()

		for _, param := range params {

			exp, err := param(ctx)

			if err != nil {
				return nil, derp.Wrap(err, "scope.And", "Error executing presto.ScopeFunc")
			}

			result = result.Add(exp)
		}

		return result, nil
	}
}
