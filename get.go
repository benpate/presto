package presto

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/labstack/echo"
)

// Get returns an HTTP handler that knows how to retrieve a single record from the collection
func (collection *Collection) Get(roles ...RoleFunc) echo.HandlerFunc {

	return func(context echo.Context) error {

		service := collection.serviceFunc()

		objectID := context.Param("id")

		// TODO: Etags

		object, err := service.LoadObject(objectID)

		if err != nil {
			return derp.NewWithCode("presto.Get", "Error loading object", err, 500, objectID).Report()
		}

		// TODO: Update cache

		// Check roles to make sure that we're allowed to view this object
		for _, role := range roles {
			if role(context) == false {
				return context.String(http.StatusUnauthorized, "")
			}
		}

		return context.JSON(http.StatusOK, object)
	}
}
