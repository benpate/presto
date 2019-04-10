package presto

// Collection provides all of the HTTP hanlers for a specific domain object,
// or collection of records
type Collection struct {
	serviceFunc ServiceFunc
}

// NewCollection returns a fully populated Collection object
func NewCollection(serviceFunc ServiceFunc, scope ScopeFunc) *Collection {
	return &Collection{
		serviceFunc: serviceFunc,
	}
}
