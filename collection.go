package presto

import "github.com/labstack/echo/v4"

// Collection provides all of the HTTP hanlers for a specific domain object,
// or collection of records
type Collection struct {
	router  *echo.Echo
	factory ServiceFactory
	prefix  string
}

// NewCollection returns a fully populated Collection object
func NewCollection(router *echo.Echo, factory ServiceFactory, prefix string) *Collection {
	return &Collection{
		router:  router,
		factory: factory,
		prefix:  prefix,
	}
}
