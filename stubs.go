package presto

import (
	"net/http"

	"github.com/benpate/data"
	"github.com/benpate/derp"
)

// Get does *most* of the work of handling an http.GET request.  It checks scopes and roles, then loads an object and
// returns it along with a valid HTTP status code.  You can use this as a shortcut in your own HTTP handler functions,
// and can wrap additional logic before and after this work.
func Get(ctx Context, service Service, cache Cache, scopes ScopeFuncSlice, roles RoleFuncSlice) (int, data.Object) {

	// If the object has an ETag, and it matches the value in the cache,
	// then we don't need to proceed any further.
	if cache != nil {

		// If the request includes an ETag
		if etag := ctx.Request().Header.Get("ETag"); etag != "" {

			// Try to find the correct ETag in the cache, using he context.Path() as the cache key
			if cache.Get(ctx.Path()) == etag {

				// If found, then we don't need to go any further.  Return NotModified
				return http.StatusNotModified, nil
			}
		}
	}

	// Use scoper functions to create query filter for this object
	filter, err := scopes.Evaluate(ctx)

	if err != nil {
		err = derp.Wrap(err, "presto.Get", "Error determining scope", ctx)
		derp.Report(err)
		return err.Code, nil
	}

	// Load the object from the database
	object, err := service.LoadObject(filter)

	if err != nil {
		err = derp.Wrap(err, "presto.Get", "Error loading object", filter)
		derp.Report(err)
		return err.Code, nil
	}

	// Try to update the ETag in the cache
	if cache != nil {
		if object, ok := object.(ETagger); ok {
			if err := cache.Set(ctx.Path(), object.ETag()); err != nil {
				err = derp.Wrap(err, "presto.Get", "Error setting cache value", object)
				derp.Report(err)
				return err.Code, nil
			}
		}
	}

	// Check roles to make sure that we're allowed to view this object
	if roles.Evaluate(ctx, object) == false {
		return http.StatusUnauthorized, nil
	}

	return http.StatusOK, object
}
