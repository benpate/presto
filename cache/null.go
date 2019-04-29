package cache

import "github.com/benpate/derp"

// Null cache is a 100% perfect cache that never stores any data, and always reports
// cache misses.  It is an empty placeholder where an actual cache would go.
type Null struct{}

func (null Null) Get(key string) string {
	return ""
}

func (null Null) Set(key string, value string) *derp.Error {
	return nil
}
