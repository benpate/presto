package presto

import "github.com/labstack/echo/v4"

// Collection provides all of the HTTP hanlers for a specific domain object,
// or collection of records
type Collection struct {
	router      *echo.Echo
	serviceFunc ServiceFunc
	prefix      string
	scopes      []ScopeFunc
}

// NewCollection returns a fully populated Collection object
func NewCollection(router *echo.Echo, serviceFunc ServiceFunc, prefix string) *Collection {
	return &Collection{
		router:      router,
		serviceFunc: serviceFunc,
		prefix:      prefix,
		scopes:      []ScopeFunc{IDScope},
	}
}

// WithScopes replaces the default scope with a new list of ScopeFuncs
func (collection *Collection) WithScopes(scopes ...ScopeFunc) *Collection {
	collection.scopes = scopes

	return collection
}
