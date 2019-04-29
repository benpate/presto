package presto

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// Get returns an HTTP handler that knows how to retrieve a single record from the collection
func (collection *Collection) Get(roles ...RoleFunc) *Collection {

	handler := func(context echo.Context) error {

		service := collection.serviceFunc()
		defer service.Close()

		objectID := context.Param("id")

		// If the object has an ETag, and it matches the value in the cache,
		// then we don't need to proceed any further.
		if etag := context.Request().Header.Get("ETag"); etag != "" {
			if ETagCache.Get(objectID) == etag {
				return context.NoContent(http.StatusNotModified)
			}
		}

		// Load the object from the database
		object, err := service.LoadObject(objectID)

		if err != nil {
			return derp.Wrap(err, "presto.Get", "Error loading object", objectID).Report()
		}

		// Try to update the ETag in the cache
		if err := ETagCache.Set(objectID, object.ETag()); err != nil {
			return derp.Wrap(err, "presto.Get", "Error setting cache value", object).Report()
		}

		// TODO: Update cache

		// Check roles to make sure that we're allowed to view this object
		for _, role := range roles {
			if role(context, object) == false {
				return context.NoContent(http.StatusUnauthorized)
			}
		}

		return context.JSON(http.StatusOK, object)
	}

	// Register the handler with the router
	collection.router.GET(collection.prefix+"/:id", handler)

	// Return the collection, so that we can chain function calls.
	return collection
}
