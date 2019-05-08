package presto

import (
	"github.com/benpate/data"
	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// ServiceFactory is an interface for objects that generate service sessions.
// Each session represents a single HTTP request, which can potentially span
// multiple database calls.  This gives the factory an opportunity to
// initialize a new database session for each HTTP request.
type ServiceFactory interface {
	Service(name string) Service
}

// Service defines all of the functions that a service must provide to work with Presto.
// It relies on the generic Object interface to load and save objects of any type.
// GenericServices will likely include additional business logic that is triggered when a
// domain object is created, edited, or deleted, but this is hidden from presto.
type Service interface {

	// NewObject creates a newly initialized object that is ready to use
	NewObject() data.Object

	// Load retrieves a single object from the database
	LoadObject(criteria data.Expression) (data.Object, *derp.Error)

	// Save inserts/updates a single object in the database
	SaveObject(object data.Object, comment string) *derp.Error

	// Delete removes a single object from the database
	DeleteObject(object data.Object, comment string) *derp.Error

	// Close cleans up any connections opened by the service.
	Close()
}

// RoleFunc is a function signature that validates a user's permission to access a particular object
type RoleFunc func(context echo.Context, object data.Object) bool

// ETagger interface wraps the ETag function, which tells presto whether or not an object
// supports ETags.  Presto uses ETags to automatically support optimistic locking of files,
// as well as saving time and bandwidth using 304: "Not Modified" responses when possible.
type ETagger interface {

	// ETag returns a version-unique string that helps determine if an object has changed or not.
	ETag() string
}

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
