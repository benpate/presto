package presto

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/benpate/derp"
)

// Delete returns an HTTP handler that knows how to delete records from the collection
func (collection *Collection) Delete(role RoleFunc) echo.HandlerFunc {
	return func(context echo.Context) error {

		service := collection.serviceFunc()

		// Try to load the record from the database
		object, err := service.LoadObject(context.Param("id"))

		if err != nil {
			return derp.New("presto.Get", "Error loading object", err, RequestInfo(context)).Report()
		}

		// Check role to make sure that we're allowed to touch this object
		if role(context) == false {
			return context.String(http.StatusUnauthorized, "")
		}

		// Try to update the record in the database
		if err := service.DeleteObject(object, "DELETE COMMENT HERE"); err != nil {
			return derp.New("presto.Delete", "Error deleting object", err, object, RequestInfo(context)).Report()
		}

		// TODO: Flush Etags & cache

		// Return the newly updated record to the caller.
		return context.String(http.StatusNoContent, "")
	}
}
