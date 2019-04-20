package cache

import "github.com/benpate/derp"

// Redis represents a TODO:
type Redis struct {
}

// NewRedis returns a fully initialized redis cache
func NewRedis() *Redis {
	return &Redis{}
}

// Get retrieves a single value from the cache.  If the value does not exist
// in the cache, then "" is returned
func (redis *Redis) Get(key string) string {
	return ""
}

// Set adds/updates a value in the cache.
func (redis *Redis) Set(key string, value string) *derp.Error {
	return nil
}
