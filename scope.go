package presto

import "github.com/benpate/data"
import "github.com/benpate/derp"

// ScopeFunc is the function signature for a function that can limit database
// queries to a particular "scope".  It inspects the provided context and
// returns criteria that will be passed to all database queries.
type ScopeFunc func(context Context) (data.Expression, *derp.Error)

// DefaultScope maps all of the route parameters directly into a scope, matching the names used in the route itself.
// It is the default behavior for presto, and should serve most use cases.
func DefaultScope(ctx Context) (data.Expression, *derp.Error) {

	criteria := data.Expression{}

	for _, param := range ctx.ParamNames() {

		if value := ctx.Param(param); value != "" {
			criteria.Add(param, "=", value)
		} else {
			return nil, derp.New(derp.CodeBadRequestError, "presto.RouteScope", "Parameter cannot be empty", param)
		}
	}

	// Otherwise, scan all items.
	return criteria, nil
}
