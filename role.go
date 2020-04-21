package presto

import (
	"github.com/benpate/data"
)

// RoleFunc is a function signature that validates a user's permission to access a particular object
type RoleFunc func(context Context, object data.Object) bool

// RoleFuncSlice defines behaviors for a slice of RoleFuncs
type RoleFuncSlice []RoleFunc

// Evaluate resolves all of the RoleFuncs using a Context and a data.Object
func (roles RoleFuncSlice) Evaluate(ctx Context, object data.Object) bool {

	for _, role := range roles {
		if role(ctx, object) == false {
			return false
		}
	}

	return true
}
