package presto

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// Get returns an HTTP handler that knows how to retrieve a single record from the collection
func (collection *Collection) Get(roles ...RoleFunc) *Collection {

	handler := func(ctx echo.Context) error {

		service := collection.factory()
		defer service.Close()

		// Use scoper functions to create query filter for this object
		filter, err := collection.getScope(ctx)

		if err != nil {
			err = derp.Wrap(err, "presto.Get", "Error determining scope", ctx).Report()
			return ctx.NoContent(err.Code)
		}

		// If the object has an ETag, and it matches the value in the cache,
		// then we don't need to proceed any further.
		if cache := collection.getCache(); cache != nil {
			if etag := ctx.Request().Header.Get("ETag"); etag != "" {

				// Use the context.Path() as the cache key
				if cache.Get(ctx.Path()) == etag {
					return ctx.NoContent(http.StatusNotModified)
				}
			}
		}

		// Load the object from the database
		object, err := service.LoadObject(filter)

		if err != nil {
			err = derp.Wrap(err, "presto.Get", "Error loading object", filter).Report()
			return ctx.NoContent(err.Code)
		}

		// Try to update the ETag in the cache
		if object, ok := object.(ETagger); ok {
			if cache := collection.getCache(); cache != nil {
				if err := cache.Set(ctx.Path(), object.ETag()); err != nil {
					err = derp.Wrap(err, "presto.Get", "Error setting cache value", object).Report()
					return ctx.NoContent(err.Code)
				}
			}
		}

		// Check roles to make sure that we're allowed to view this object
		for _, role := range roles {
			if role(ctx, object) == false {
				return ctx.NoContent(http.StatusUnauthorized)
			}
		}

		return ctx.JSON(http.StatusOK, object)
	}

	// Register the handler with the router
	globalRouter.GET(collection.prefix+"/:"+collection.token, handler)

	// Return the collection, so that we can chain function calls.
	return collection
}
