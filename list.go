package presto

// List returns an HTTP handler that knows how to list a series of records from the collection
func (collection *Collection) List(roles ...RoleFunc) *Collection {

	handler := func(context Context) error {
		return nil
	}

	// Register handler with the router
	globalRouter.GET(collection.prefix, handler)

	// Return collection, so that we can chain calls if needed.
	return collection
}
