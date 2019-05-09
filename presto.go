/*
Package presto gives you a strong foundation for creating REST interfaces using
Go and the [Echo Router](http://echo.labstack.com)
*/
package presto

import "github.com/labstack/echo/v4"

var globalRouter *echo.Echo
var globalCache Cache
var globalScopes []ScopeFunc

// UseRouter sets the echo router that presto will use to register HTTP handlers.
func UseRouter(router *echo.Echo) {
	globalRouter = router
}

// UseCache sets the global cache for all presto endpoints.
func UseCache(cache Cache) {
	globalCache = cache
}

// UseScopes sets global settings for all collections that are managed by presto
func UseScopes(scopes ...ScopeFunc) {
	globalScopes = scopes
}
