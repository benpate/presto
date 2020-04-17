package presto

import (
	"github.com/benpate/data"
	"github.com/benpate/data/expression"
	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// Collection provides all of the HTTP hanlers for a specific domain object,
// or collection of records
type Collection struct {
	factory ServiceFunc
	prefix  string
	scopes  []ScopeFunc
	cache   Cache
	token   string
}

// NewCollection returns a fully populated Collection object
func NewCollection(factory ServiceFunc, prefix string) *Collection {
	return &Collection{
		factory: factory,
		prefix:  prefix,
		scopes:  []ScopeFunc{DefaultScope},
		token:   "id",
	}
}

// UseScopes replaces the default scope with a new list of ScopeFuncs
func (collection *Collection) UseScopes(scopes ...ScopeFunc) *Collection {
	collection.scopes = scopes

	return collection
}

// UseCache adds a local ETag cache for this collection only
func (collection *Collection) UseCache(cache Cache) *Collection {
	collection.cache = cache

	return collection
}

// UseToken overrides the default "token" variable that is appended to all GET, PUT, PATCH, and DELETE
// routes, and is used as the unique identifier of the record being created, read, updated, or deleted.
func (collection *Collection) UseToken(token string) *Collection {
	collection.token = token

	return collection
}

// getCache locates the correct cache to use for this collection ~ either the global cache, or a local one.
func (collection *Collection) getCache() Cache {

	if collection.cache != nil {
		return collection.cache
	}

	return globalCache
}

// getScope executes each scoper function for this context and returns a data expression
func (collection *Collection) getScope(ctx echo.Context) (expression.Expression, *derp.Error) {

	result := expression.And()

	if globalScopes != nil {

		for _, scope := range globalScopes {

			next, err := scope(ctx)

			if err != nil {
				return result, derp.Wrap(err, "presto.getScope", "Error executing global scope function")
			}

			result = expression.And(result, next)
		}
	}

	for _, scope := range collection.scopes {

		next, err := scope(ctx)

		if err != nil {
			return result, derp.Wrap(err, "presto.getScope", "Error executing scope function")
		}

		result = expression.And(result, next)
	}

	return result, nil
}

// isEtagConflict returns TRUE if the provided ETag DOES NOT match the value in the cache.
// This is used for (very) optimistic locking.  If this returns a FALSE, then the value
// must STILL be double checked AFTER we load the object, because its might not be in the cache.
func (collection *Collection) isETagConflict(ctx Context, object data.Object) bool {

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
