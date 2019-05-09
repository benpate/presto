package presto

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// Patch returns an HTTP handler that knows how to update in the collection
func (collection *Collection) Patch(roles ...RoleFunc) *Collection {

	handler := func(ctx echo.Context) error {

		service := collection.factory()
		defer service.Close()

		// Use scoper functions to create query criteria for this object
		filter, err := collection.getScope(ctx)

		if err != nil {
			err = derp.Wrap(err, "presto.Patch", "Error determining scope", ctx).Report()
			return ctx.NoContent(err.Code)
		}

		// Try to load the record from the database
		object, err := service.LoadObject(filter)

		if err != nil {
			err = derp.Wrap(err, "presto.Patch", "Error loading object", RequestInfo(ctx)).Report()
			return ctx.NoContent(err.Code)
		}

		// Check roles (before update) to make sure that we're allowed to touch this object
		for _, role := range roles {
			if role(ctx, object) == false {
				return ctx.NoContent(http.StatusUnauthorized)
			}
		}

		// Double check that the ETag matches the object ~ used for optimistic locking.
		if collection.isETagConflict(ctx, object) {
			return ctx.NoContent(http.StatusConflict)
		}

		// Update the object with new information
		if er := ctx.Bind(object); er != nil {
			err := derp.New(derp.CodeBadRequestError, "presto.Patch", "Error binding object", er, object, RequestInfo(ctx)).Report()
			return ctx.NoContent(err.Code)
		}

		// Check roles again (after update) to make sure that we're making valid changes that still let us "own" this object.
		for _, role := range roles {
			if role(ctx, object) == false {
				return ctx.NoContent(http.StatusUnauthorized)
			}
		}

		// Try to update the record in the database
		if err := service.SaveObject(object, ctx.Request().Header.Get("X-Comment")); err != nil {
			err = derp.Wrap(err, "presto.Patch", "Error saving object", object, RequestInfo(ctx)).Report()
			return ctx.NoContent(err.Code)
		}

		// Try to update the ETag cache
		if object, ok := object.(ETagger); ok {
			if cache := collection.getCache(); cache != nil {
				if err := cache.Set(ctx.Path(), object.ETag()); err != nil {
					err = derp.Wrap(err, "presto.Patch", "Error updating ETag cache", object).Report()
					return ctx.NoContent(err.Code)
				}
			}
		}

		// Return the newly updated record to the caller.
		return ctx.JSON(http.StatusOK, object)
	}

	// Register the handler with the router
	globalRouter.PATCH(collection.prefix+"/:"+collection.token, handler)

	// Return the collection so that we can chain requests.
	return collection
}
