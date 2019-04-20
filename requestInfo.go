package presto

import "github.com/labstack/echo/v4"

// RequestInfo inspects a request and returns any information that might be useful
// for debugging problems.  It is primarily used by internal methods whenever there's
// a problem with a request.
func RequestInfo(context echo.Context) map[string]string {
	return map[string]string{}
}
