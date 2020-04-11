package presto

import "github.com/benpate/data"

// ServiceFunc is a function that can generate new services/sessions.
// Each session represents a single HTTP request, which can potentially span
// multiple database calls.  This gives the factory an opportunity to
// initialize a new database session for each HTTP request.
type ServiceFunc func() Service

// Service defines all of the functions that a service must provide to work with Presto.
// It relies on the generic Object interface to load and save objects of any type.
// GenericServices will likely include additional business logic that is triggered when a
// domain object is created, edited, or deleted, but this is hidden from presto.
type Service interface {

	// NewObject creates a newly initialized object that is ready to use
	NewObject() data.Object

	// ListObjects returns an iterator the returns all objects
	ListObjects(criteria data.Expression, options ...data.Option) (data.Iterator, *derp.Error)

	// Load retrieves a single object from the database
	LoadObject(criteria data.Expression, options ...data.Option) (data.Object, *derp.Error)

	// Save inserts/updates a single object in the database
	SaveObject(object data.Object, comment string) *derp.Error

	// Delete removes a single object from the database
	DeleteObject(object data.Object, comment string) *derp.Error

	// Close cleans up any connections opened by the service.
	Close()
}


