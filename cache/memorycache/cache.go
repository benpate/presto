package memorycache

import "github.com/benpate/derp"

// Memory represents a very crude in-memory cache that should only be used for testing purposes.
// It DOES NOT implement a maximum size, and can easily overflow server memory if used in a
// production system.  DO NOT USE THIS IN PRODUCTION.
type Memory map[string]string

// New returns a fully initialized memory cache
func New() *Memory {
	return &Memory{}
}

// Get retrieves a single value from the cache.  If the value does not exist
// in the cache, then "" is returned
func (memory *Memory) Get(key string) string {
	return (*memory)[key]
}

// Set adds/updates a value in the cache.
func (memory *Memory) Set(key string, value string) *derp.Error {
	(*memory)[key] = value
	return nil
}
