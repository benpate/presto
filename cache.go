package presto

import "github.com/benpate/derp"

// Cache maintains fast access to key/value pairs that are used to check ETags of incoming requests.
// By default, Presto uses a Null cache, that simply reports cache misses for every request.  However,
// this can be extended by the user, with any external caching system that matches this interface.
type Cache interface {

	// Get returns the cache value (ETag) corresponding to the argument (objectID) provided.
	// If a value is not found, then Get returns empty string ("")
	Get(objectID string) string

	// Set updates the value in the cache, returning a derp.Error in case there was a problem.
	Set(objectID string, value string) *derp.Error
}
