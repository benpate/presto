package presto

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/benpate/derp"
)

// Put returns an HTTP handler that knows how to update in the collection
func (collection Collection) Put(roles ...RoleFunc) echo.HandlerFunc {

	return func(context echo.Context) error {

		service := collection.serviceFunc()

		// Try to load the record from the database
		object, err := service.LoadObject(context.Param("id"))

		if err != nil {
			return derp.New("presto.Get", "Error loading object", err, RequestInfo(context)).Report()
		}

		// Check roles (before update) to make sure that we're allowed to touch this object
		for _, role := range roles {
			if role(context) == false {
				return context.String(http.StatusUnauthorized, "")
			}
		}

		// Create a new object to populate from this point forward
		object = service.NewObject()

		// Update the object with new information
		if err := context.Bind(object); err != nil {
			return derp.NewWithCode("presto.Put", "Error binding object", err, 500, object, RequestInfo(context)).Report()
		}

		// Check roles again (after update) to make sure that we're making valid changes that still let us "own" this object.
		for _, role := range roles {
			if role(context) == false {
				return context.String(http.StatusUnauthorized, "")
			}
		}

		// Try to update the record in the database
		if err := service.SaveObject(object, "SAVE COMMENT HERE"); err != nil {
			return derp.New("presto.Put", "Error saving object", err, object, RequestInfo(context)).Report()
		}

		// TODO: Flush Etags & cache

		// Return the newly updated record to the caller.
		return context.JSON(http.StatusOK, object)
	}
}
