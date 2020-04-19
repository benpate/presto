package presto

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// Put returns an HTTP handler that knows how to update in the collection
func (collection *Collection) Put(roles ...RoleFunc) *Collection {

	handler := func(ctx echo.Context) error {

		service := collection.factory()
		defer service.Close()

		// We'll use this to determine if we're (T) creating a new object, or (F) updating an existing one.
		isNewObject := false

		// Use scoper functions to create query criteria for this object
		filter, err := collection.getScopeWithToken(ctx)

		if err != nil {
			err = derp.Wrap(err, "presto.Patch", "Error determining scope", ctx).Report()
			return ctx.NoContent(err.Code)
		}

		// Try to load the record from the database
		object, err := service.LoadObject(filter)

		if err != nil {

			// If the error is ANYTHING BUT a "Not Found" error, then it's a legitimate error,
			// and we need to shut this whole thing down right now.
			if err.NotFound() == false {
				err = derp.Wrap(err, "presto.Put", "Error loading object", RequestInfo(ctx)).Report()
				return ctx.NoContent(err.Code)
			}

			// Fall through to here means that this is just a "Not Found", which is
			// OK at this stage.
			isNewObject = true
		}

		// Only run these checks if we're updating an existing object:
		// 1) Make sure that ETags match (for optimistic locking)
		// 2) Make sure that the user has access to the existing object.
		if isNewObject == false {

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
		}

		// Create a new object to populate from this point forward
		object = service.NewObject()

		// Update the object with new information
		if er := ctx.Bind(object); er != nil {
			err := derp.New(derp.CodeBadRequestError, "presto.Put", "Error binding object", er, object, RequestInfo(ctx)).Report()
			return ctx.NoContent(err.Code)
		}

		// Check roles again (after update, before save) to make sure that we're making valid changes that still let us "own" this object.
		for _, role := range roles {
			if role(ctx, object) == false {
				return ctx.NoContent(http.StatusUnauthorized)
			}
		}

		// Try to update the record in the database
		if err := service.SaveObject(object, ctx.Request().Header.Get("X-Comment")); err != nil {
			err = derp.Wrap(err, "presto.Put", "Error saving object", object, RequestInfo(ctx)).Report()
			return ctx.NoContent(err.Code)
		}

		// Try to update the ETag cache
		if object, ok := object.(ETagger); ok {
			if cache := collection.getCache(); cache != nil {
				if err := cache.Set(ctx.Path(), object.ETag()); err != nil {
					err = derp.Wrap(err, "presto.Put", "Error updating cache", object).Report()
					return ctx.NoContent(err.Code)
				}
			}
		}

		// Return the newly updated record to the caller.
		return ctx.JSON(http.StatusOK, object)
	}

	// Register the handler with the router
	globalRouter.PUT(collection.prefix+"/:"+collection.token, handler)

	// Return the collection so that users can chain requests
	return collection
}
