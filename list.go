package presto

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// List returns an HTTP handler that knows how to list a series of records from the collection
func (collection *Collection) List(roles ...RoleFunc) *Collection {

	handler := func(ctx echo.Context) error {

		service := collection.serviceFunc(ctx.Request().Context())
		defer service.Close()

		// Use scoper functions to create query filter for this object
		filter, err := collection.getScope(ctx)

		if err != nil {
			err = derp.Wrap(err, "presto.List", "Error determining scope", ctx).Report()
			return ctx.NoContent(err.Code)
		}

		// TODO: Add pagination logic here.

		// Load the object from the database
		it, err := service.ListObjects(filter)

		if err != nil {
			err = derp.Wrap(err, "presto.List", "Error loading object", filter).Report()
			return ctx.NoContent(err.Code)
		}

		// TODO: add HTTP headers here...

		// Get a new object we can populate data into
		object := service.NewObject()
		first := true

		var buffer bytes.Buffer

		buffer.WriteByte('[')

		// Loop through the iterator to return a data structure.
		for it.Next(object) {

			// Check all roles to make sure that we're allowed to view this object
			for _, role := range roles {

				// If any role fails, then we won't include this object in the results returned to the user.
				if role(ctx, object) == false {
					continue
				}
			}

			// Try to marshal the object into JSON.
			record, err := json.Marshal(object)

			// If we're unable to marshal the object, then the whole result is b0rked.
			// So, flag the error and exit without returning any real data.
			if err != nil {
				// Need a real error message here.
				return ctx.String(http.StatusInternalServerError, "")
			}

			if first {
				first = false
			} else {
				buffer.WriteByte(',')
			}

			buffer.Write(record)
		}

		buffer.WriteByte(']')

		// Check roles to make sure that we're allowed to view this object
		for _, role := range roles {
			if role(ctx, object) == false {
				return ctx.NoContent(http.StatusUnauthorized)
			}
		}

		return ctx.JSON(http.StatusOK, object)
	}

	// Register handler with the router
	globalRouter.GET(collection.prefix, handler)

	// Return collection, so that we can chain calls if needed.
	return collection
}
