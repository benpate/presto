package presto

import (
	"github.com/benpate/data"
	"github.com/benpate/data/expression"
	"github.com/benpate/derp"
)

// Collection provides all of the HTTP hanlers for a specific domain object,
// or collection of records
type Collection struct {
	serviceFunc ServiceFunc
	prefix      string
	scopes      ScopeFuncSlice
	cache       Cache
	token       string
}

// NewCollection returns a fully populated Collection object
func NewCollection(serviceFunc ServiceFunc, prefix string) *Collection {
	return &Collection{
		serviceFunc: serviceFunc,
		prefix:      prefix,
		scopes:      ScopeFuncSlice{},
		token:       "id",
	}
}

// UseScopes replaces the default scope with a new list of ScopeFuncs
func (collection *Collection) UseScopes(scopes ...ScopeFunc) *Collection {
	collection.scopes = ScopeFuncSlice(scopes)

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
func (collection *Collection) getScopes() ScopeFuncSlice {

	return append(globalScopes, collection.scopes...)
}

// getScopeWithToken returns all of the scopes that are valid for this collection, and adds a ScopeFunc for the collection.token
func (collection *Collection) getScopesWithToken() (ScopeFuncSlice, *derp.Error) {

	scopes := collection.getScopes()

	tokenScopeFunc := func(context Context) (expression.Expression, *derp.Error) {

		if value := context.Param(collection.token); value != "" {
			return expression.New(collection.token, expression.OperatorEqual, value), nil
		}

		return nil, derp.New(derp.CodeBadRequestError, "collection.getScopeWithToken", "Token cannot be empty", collection.token)
	}

	return append(scopes, tokenScopeFunc), nil
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
