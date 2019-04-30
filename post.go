package presto

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/benpate/derp"
)

// Post returns an HTTP handler that knows how to create new objects in the collection
func (collection *Collection) Post(roles ...RoleFunc) *Collection {

	handler := func(context echo.Context) error {

		service := collection.factory.Service(collection.name)
		defer service.Close()

		// Create a new, empty object
		object := service.NewObject()

		// TODO: How to enforce SCOPE here?

		// Update the object with new information
		if er := context.Bind(object); er != nil {
			err := derp.New(derp.CodeBadRequestError, "presto.Post", "Error binding object", er.Error(), object, RequestInfo(context)).Report()
			return context.NoContent(err.Code)
		}

		// Check roles (after update) to make sure that we're making valid changes that still let us "own" this object.
		for _, role := range roles {
			if role(context, object) == false {
				return context.NoContent(http.StatusUnauthorized)
			}
		}

		// Try to update the record in the database
		if err := service.SaveObject(object, "SAVE COMMENT HERE"); err != nil {
			err = derp.Wrap(err, "presto.Post", "Error saving object", object, RequestInfo(context)).Report()
			return context.NoContent(err.Code)
		}

		// Try to reset the ETag cache
		if cache := collection.getCache(); cache != nil {
			if err := cache.Set(object.ID(), object.ETag()); err != nil {
				err = derp.Wrap(err, "presto.Post", "Error setting cache value", object)
				return context.NoContent(err.Code)
			}
		}

		// Return the newly updated record to the caller.
		return context.JSON(http.StatusOK, object)
	}

	collection.router.POST(collection.prefix, handler)

	return collection
}
