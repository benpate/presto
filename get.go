package presto

import (
	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// Get returns an HTTP handler that knows how to retrieve a single record from the collection
func (collection *Collection) Get(roles ...RoleFunc) *Collection {

	handler := func(ctx echo.Context) error {

		service := collection.serviceFunc(ctx.Request().Context())
		defer service.Close()

		scopes, err := collection.getScopesWithToken()

		if err != nil {
			err = derp.Wrap(err, "collection.Get", "Error defining scopes")
			derp.Report(err)
			return ctx.String(err.Code, "")
		}

		code, object := Get(ctx, service, collection.getCache(), scopes, roles)

		if object == nil {
			return ctx.String(code, "")
		}

		return ctx.JSON(code, object)
	}

	// Register the handler with the router
	globalRouter.GET(collection.prefix+"/:"+collection.token, handler)

	// Return the collection, so that we can chain function calls.
	return collection
}
