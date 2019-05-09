/*
Package nullcache implements the presto.Cache interface with an empty
data structure that never stores any data, and always reports a cache
"miss".  It is an empty placeholder for testing purposes only, that
can be put anywhere an actual cache would go.
*/
package nullcache

import "github.com/benpate/derp"

// NullCache implements the presto.Cache interface with an empty
// data structure that never stores any data, and always reports a cache
// "miss".  It is an empty placeholder for testing purposes only, that
// can be put anywhere an actual cache would go.
type NullCache struct{}

// New returns a new cache that is ready to "use"
func New() *NullCache {
	return &NullCache{}
}

// Get returns the (null) value for any required cache key
func (null *NullCache) Get(key string) string {
	return ""
}

// Set is a NOOP for setting a value in the cache.
func (null *NullCache) Set(key string, value string) *derp.Error {
	return nil
}
