/*
Package presto gives you a strong foundation for creating REST interfaces using
Go and the [Echo Router](http://echo.labstack.com)
*/
package presto

var globalCache Cache

// WithCache sets the global cache for all presto endpoints.
func WithCache(cache Cache) {
	globalCache = cache
}
