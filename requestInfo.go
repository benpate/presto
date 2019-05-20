package presto

// RequestInfo inspects a request and returns any information that might be useful
// for debugging problems.  It is primarily used by internal methods whenever there's
// a problem with a request.
func RequestInfo(context Context) map[string]string {
	return map[string]string{}
}
