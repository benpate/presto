package scope

import (
	"github.com/benpate/data/expression"
	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// DefaultScope maps all of the route parameters directly into a scope, matching the names used in the route itself.
// It is the default behavior for presto, and should serve most use cases.
func DefaultScope(ctx echo.Context) (expression.Expression, *derp.Error) {

	criteria := expression.And()

	for _, param := range ctx.ParamNames() {

		if value := ctx.Param(param); value != "" {
			criteria.And(param, "=", value)
		} else {
			return nil, derp.New(derp.CodeBadRequestError, "presto.RouteScope", "Parameter cannot be empty", param)
		}
	}

	// Otherwise, scan all items.
	return criteria, nil
}
