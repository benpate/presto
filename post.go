package presto

import (
	"net/http"

	"github.com/benpate/derp"
)

// Post returns an HTTP handler that knows how to create new objects in the collection
func (collection *Collection) Post(roles ...RoleFunc) *Collection {

	handler := func(ctx Context) error {

		service := collection.factory()
		defer service.Close()

		// Create a new, empty object
		object := service.NewObject()

		// Update the object with new information
		if er := ctx.Bind(object); er != nil {
			err := derp.New(derp.CodeBadRequestError, "presto.Post", "Error binding object", er.Error(), object, RequestInfo(ctx)).Report()
			return ctx.NoContent(err.Code)
		}

		// Check roles (after update) to make sure that we're making valid changes that still let us "own" this object.
		for _, role := range roles {
			if role(ctx, object) == false {
				return ctx.NoContent(http.StatusUnauthorized)
			}
		}

		// Try to update the record in the database
		if err := service.SaveObject(object, ctx.Request().Header.Get("X-Comment")); err != nil {
			err = derp.Wrap(err, "presto.Post", "Error saving object", object, RequestInfo(ctx)).Report()
			return ctx.NoContent(err.Code)
		}

		// Try to reset the ETag cache
		if object, ok := object.(ETagger); ok {
			if cache := collection.getCache(); cache != nil {
				if err := cache.Set(ctx.Path(), object.ETag()); err != nil {
					err = derp.Wrap(err, "presto.Post", "Error setting cache value", object)
					return ctx.NoContent(err.Code)
				}
			}
		}

		// Return the newly updated record to the caller.
		return ctx.JSON(http.StatusOK, object)
	}

	globalRouter.POST(collection.prefix, handler)

	return collection
}
