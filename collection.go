package presto

import "github.com/labstack/echo/v4"

// Collection provides all of the HTTP hanlers for a specific domain object,
// or collection of records
type Collection struct {
	router  *echo.Echo
	factory ServiceFactory
	name    string
	prefix  string
	scopes  []ScopeFunc
	cache   Cache
}

// NewCollection returns a fully populated Collection object
func NewCollection(router *echo.Echo, factory ServiceFactory, name string, prefix string) *Collection {
	return &Collection{
		router:  router,
		factory: factory,
		name:    name,
		prefix:  prefix,
		scopes:  []ScopeFunc{IDScope},
	}
}

// WithScopes replaces the default scope with a new list of ScopeFuncs
func (collection *Collection) WithScopes(scopes ...ScopeFunc) *Collection {
	collection.scopes = scopes

	return collection
}

// WithCache adds a local ETag cache for this collection only
func (collection *Collection) WithCache(cache Cache) *Collection {
	collection.cache = cache

	return collection
}

func (collection *Collection) getCache() Cache {

	if collection.cache != nil {
		return collection.cache
	}

	return globalCache
}
