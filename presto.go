/*
Package presto gives you a strong foundation for creating REST interfaces using
Go and the [Echo Router](http://echo.labstack.com)
*/
package presto

var globalCache Cache
var globalScopes []ScopeFunc

// WithCache sets the global cache for all presto endpoints.
func WithCache(cache Cache) {
	globalCache = cache
}

// WithScopes sets global settings for all collections that are managed by presto
func WithScopes(scopes ...ScopeFunc) {
	globalScopes = scopes
}
