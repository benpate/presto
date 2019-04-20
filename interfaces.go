package presto

import (
	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// ServiceFactory is an interface for objects that generate service sessions.
// Each session represents a single HTTP request, which can potentially span
// multiple database calls.  This gives the factory an opportunity to
// initialize a new database session for each HTTP request.
type ServiceFactory interface {
	Service() GenericService
}

// GenericService defines all of the functions that a service must provide to work with Presto.
// It is called a Generic service, because it relies on the generic Object interface to load and
// save objects of any type.
// GenericServices will likely include additional business logic that is triggered when a
// domain object is created, edited, or deleted, but this is hidden from presto.
type GenericService interface {

	// New creates a newly initialized object that is ready to use
	GenericNew() Object

	// Load retrieves a single object from the database
	GenericLoad(objectID string) (Object, *derp.Error)

	// Save inserts/updates a single object in the database
	GenericSave(object Object, comment string) *derp.Error

	// Delete removes a single object from the database
	GenericDelete(object Object, comment string) *derp.Error
}

// RoleFunc is a function signature that validates a user's permission to access a particular object
type RoleFunc func(echo.Context, Object) bool

// Object wraps all of the methods that a Domain Object must provide to Presto
type Object interface {

	// ID returns the primary key of the object
	ID() string

	// IsNew returns TRUE if the object has not yet been saved to the database
	IsNew() bool

	// SetCreated stamps the CreateDate and UpdateDate of the object, and makes a note
	SetCreated(comment string)

	// SetUpdated stamps the UpdateDate of the object, and makes a note
	SetUpdated(comment string)

	// SetDeleted marks the object virtually "deleted", and makes a note
	SetDeleted(comment string)

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
