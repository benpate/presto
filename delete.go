package presto

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// Delete returns an HTTP handler that knows how to delete records from the collection
func (collection *Collection) Delete(roles ...RoleFunc) *Collection {

	handler := func(ctx echo.Context) error {

		service := collection.factory(ctx.Request().Context())
		defer service.Close()

		// Use scoper functions to create query criteria for this object
		criteria, err := collection.getScopeWithToken(ctx)

		if err != nil {
			err = derp.Wrap(err, "presto.Delete", "Error determining scope", ctx).Report()
			return ctx.NoContent(err.Code)
		}

		// Try to load the record from the database
		object, err := service.LoadObject(criteria)

		if err != nil {
			err = derp.Wrap(err, "presto.Delete", "Error loading object", RequestInfo(ctx)).Report()
			return ctx.NoContent(err.Code)
		}

		// Check roles to make sure that we're allowed to touch this object
		for _, role := range roles {
			if role(ctx, object) == false {
				return ctx.NoContent(http.StatusUnauthorized)
			}
		}

		// Double check that the ETag matches the object ~ used for optimistic locking.
		if collection.isETagConflict(ctx, object) {
			return ctx.NoContent(http.StatusConflict)
		}

		// Try to update the record in the database
		if err := service.DeleteObject(object, ctx.Request().Header.Get("X-Comment")); err != nil {
			err = derp.Wrap(err, "presto.Delete", "Error deleting object", object, RequestInfo(ctx)).Report()
			return ctx.NoContent(err.Code)
		}

		// Try to remove the Etag from the cache
		if cache := collection.getCache(); cache != nil {
			if err := cache.Set(ctx.Path(), ""); err != nil {
				err = derp.Wrap(err, "presto.Delete", "Error flushing ETag cache", object)
				return ctx.NoContent(err.Code)
			}
		}

		// Return the newly updated record to the caller.
		return ctx.NoContent(http.StatusNoContent)
	}

	// Register the handler with the router
	globalRouter.DELETE(collection.prefix+"/:"+collection.token, handler)

	// Return the collection, so that we can chain function calls.
	return collection
}
