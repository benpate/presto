package presto

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/benpate/derp"
)

// Post returns an HTTP handler that knows how to create new objects in the collection
func (collection Collection) Post(roles ...RoleFunc) echo.HandlerFunc {

	return func(context echo.Context) error {

		service := collection.serviceFunc()

		// Create a new, empty object
		object := service.NewObject()

		// Update the object with new information
		if err := context.Bind(object); err != nil {
			return derp.NewWithCode("presto.Put", "Error binding object", err, 500, object, RequestInfo(context)).Report()
		}

		// Check roles (after update) to make sure that we're making valid changes that still let us "own" this object.
		for _, role := range roles {
			if role(context) == false {
				return context.String(http.StatusUnauthorized, "")
			}
		}

		// Try to update the record in the database
		if err := service.SaveObject(object, "SAVE COMMENT HERE"); err != nil {
			return derp.New("presto.Put", "Error saving object", err, object, RequestInfo(context)).Report()
		}

		// TODO: Flush Etags

		// Return the newly updated record to the caller.
		return context.JSON(http.StatusOK, object)
	}
}
