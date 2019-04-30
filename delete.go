package presto

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// Delete returns an HTTP handler that knows how to delete records from the collection
func (collection *Collection) Delete(roles ...RoleFunc) *Collection {

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

		// Check roles to make sure that we're allowed to touch this object
		for _, role := range roles {
			if role(context, object) == false {
				return context.NoContent(http.StatusUnauthorized)
			}
		}

		// Try to update the record in the database
		if err := service.DeleteObject(object, "DELETE COMMENT HERE"); err != nil {
			err = derp.Wrap(err, "presto.Delete", "Error deleting object", object, RequestInfo(context)).Report()
			return context.NoContent(err.Code)
		}

		// Try to remove the Etag from the cache
		if cache := collection.getCache(); cache != nil {
			if err := cache.Set(object.ID(), ""); err != nil {
				err = derp.Wrap(err, "presto.Delete", "Error flushing ETag cache", object)
				return context.NoContent(err.Code)
			}
		}

		// Return the newly updated record to the caller.
		return context.NoContent(http.StatusNoContent)
	}

	// Register the handler with the router
	collection.router.DELETE(collection.prefix+"/:id", handler)

	// Return the collection, so that we can chain function calls.
	return collection
}
