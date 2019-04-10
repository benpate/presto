package presto

import (
	"github.com/labstack/echo"
)

// List returns an HTTP handler that knows how to list a series of records from the collection
func (collection Collection) List(role RoleFunc) echo.HandlerFunc {

	return func(context echo.Context) error {
		return nil
	}
}
