package presto

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/benpate/derp"
)

// Patch returns an HTTP handler that knows how to update in the collection
func (collection *Collection) Patch(roles ...RoleFunc) *Collection {

	handler := func(context echo.Context) error {

		service := collection.factory.Service(collection.name)
		defer service.Close()

		// TODO: Use SCOPE here.

		// Try to load the record from the database
		object, err := service.LoadObject(context.Param("id"))

		if err != nil {
			err = derp.Wrap(err, "presto.Get", "Error loading object", RequestInfo(context)).Report()
			return context.NoContent(err.Code)
		}

		// Check roles (before update) to make sure that we're allowed to touch this object
		for _, role := range roles {
			if role(context, object) == false {
				return context.NoContent(http.StatusUnauthorized)
			}
		}

		// Update the object with new information
		if er := context.Bind(object); er != nil {
			err := derp.New(derp.CodeBadRequestError, "presto.Put", "Error binding object", er, object, RequestInfo(context)).Report()
			return context.NoContent(err.Code)
		}

		// Check roles again (after update) to make sure that we're making valid changes that still let us "own" this object.
		for _, role := range roles {
			if role(context, object) == false {
				return context.NoContent(http.StatusUnauthorized)
			}
		}

		// Try to update the record in the database
		if err := service.SaveObject(object, "SAVE COMMENT HERE"); err != nil {
			err = derp.Wrap(err, "presto.Put", "Error saving object", object, RequestInfo(context)).Report()
			return context.NoContent(err.Code)
		}

		// Try to update the ETag cache
		if cache := collection.getCache(); cache != nil {
			if err := cache.Set(object.ID(), object.ETag()); err != nil {
				err = derp.Wrap(err, "presto.Put", "Error updating ETag cache", object).Report()
				return context.NoContent(err.Code)
			}
		}

		// Return the newly updated record to the caller.
		return context.JSON(http.StatusOK, object)
	}

	// Register the handler with the router
	collection.router.PATCH(collection.prefix+"/:id", handler)

	// Return the collection so that we can chain requests.
	return collection
}
