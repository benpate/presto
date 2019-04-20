package cache

import "github.com/benpate/derp"

// Mongodb represents a TODO:
type Mongodb struct {
	connectionString string
	collection       string
}

// NewMongodb returns a fully initialized mongodb cache
func NewMongodb(connectionString string, collection string) *Mongodb {

	return &Mongodb{
		connectionString: connectionString,
		collection:       collection,
	}
}

// Get retrieves a single value from the cache.  If the value does not exist
// in the cache, then "" is returned
func (mongodb *Mongodb) Get(key string) string {
	return ""
}

// Set adds/updates a value in the cache.
func (mongodb *Mongodb) Set(key string, value string) *derp.Error {
	return nil
}
