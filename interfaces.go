package presto

import (
	"github.com/benpate/derp"
	"github.com/labstack/echo"
)

// ServiceFunc is a function that generates and initializes the required service.
type ServiceFunc func() Service

// Service wraps all of the methods that a service managing Domain Objects must provide to Presto
type Service interface {

	// New creates a newly initialized object that is ready to use
	NewObject() Object

	// Load retrieves a single object from the database
	LoadObject(string) (Object, *derp.Error)

	// Save inserts/updates a single object in the database
	SaveObject(Object, string) *derp.Error

	// Delete removes a single object from the database
	DeleteObject(Object, string) *derp.Error
}

// RoleFunc is a function signature that validates a user's permission to access a particular object
type RoleFunc func(echo.Context) bool

// ScopeFunc defines the baseline query parameters that are present in ALL queries.  This
// sets the maximum range (or scope) of the records that the requester can access
type ScopeFunc func(echo.Context) map[string]interface{}

// Object wraps all of the methods that a Domain Object must provide to Presto
type Object interface {

	// ID returns the primary key of the object
	ID() string

	// IsNew returns TRUE if the object has not yet been saved to the database
	IsNew() bool

	// SetCreated stamps the CreateDate and UpdateDate of the object, and makes a note
	SetCreated(string)

	// SetUpdated stamps the UpdateDate of the object, and makes a note
	SetUpdated(string)

	// SetDeleted marks the object virtually "deleted", and makes a note
	SetDeleted(string)
}
