package presto

import (
	"github.com/benpate/data"
)

// RoleFunc is a function signature that validates a user's permission to access a particular object
type RoleFunc func(context Context, object data.Object) bool
