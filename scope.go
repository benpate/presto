package presto

import (
	"github.com/benpate/data/expression"
	"github.com/benpate/derp"
)

// ScopeFunc is the function signature for a function that can limit database
// queries to a particular "scope".  It inspects the provided context and
// returns criteria that will be passed to all database queries.
type ScopeFunc func(context Context) (expression.Expression, *derp.Error)

// ScopeFuncSlice defines behaviors for a slice of Scopes
type ScopeFuncSlice []ScopeFunc

// Evaluate resolves all scopes into an expression (or error) using the provided Context
func (scopes ScopeFuncSlice) Evaluate(ctx Context) (expression.AndExpression, *derp.Error) {

	result := expression.And()

	for _, scope := range scopes {

		next, err := scope(ctx)

		if err != nil {
			return result, derp.Wrap(err, "presto.getScope", "Error executing scope function")
		}

		result = expression.And(result, next)
	}

	return result, nil
}
