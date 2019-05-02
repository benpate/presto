package presto

import (
	"github.com/benpate/data"
	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// Collection provides all of the HTTP hanlers for a specific domain object,
// or collection of records
type Collection struct {
	router  *echo.Echo
	factory ServiceFactory
	name    string
	prefix  string
	scopes  []ScopeFunc
	cache   Cache
	token   string
}

// NewCollection returns a fully populated Collection object
func NewCollection(router *echo.Echo, factory ServiceFactory, name string, prefix string, token string) *Collection {
	return &Collection{
		router:  router,
		factory: factory,
		name:    name,
		prefix:  prefix,
		scopes:  []ScopeFunc{RouteScope},
		token:   token,
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

// getCache locates the correct cache to use for this collection ~ either the global cache, or a local one.
func (collection *Collection) getCache() Cache {

	if collection.cache != nil {
		return collection.cache
	}

	return globalCache
}

// isEtagConflict returns TRUE if the provided ETag DOES NOT match the value in the cache.
// This is used for (very) optimistic locking.  If this returns a FALSE, then the value
// must STILL be double checked AFTER we load the object, because its might not be in the cache.
func (collection *Collection) isETagConflict(ctx echo.Context, object data.Object) bool {

	if object, ok := object.(ETagger); ok {

		etag := object.ETag()

		// If there is NO ETag for this object, then there's no conflict.  Return FALSE
		if etag == "" {
			return false
		}

		// Try to get the ETag from the request headers.
		if headerValue := ctx.Request().Header.Get("ETag"); headerValue != "" {

			if etag != headerValue {
				return true
			}
		}
	}

	return false
}

// getScope executes each scoper function for this context and returns a data expression
func (collection *Collection) getScope(ctx echo.Context) (data.Expression, *derp.Error) {

	result := data.Expression{}

	for _, scope := range collection.scopes {

		next, err := scope(ctx)

		if err != nil {
			return result, derp.Wrap(err, "presto.getScope", "Error executing scope function")
		}

		result = result.Join(next)
	}

	return result, nil
}
