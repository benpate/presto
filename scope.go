package presto

import (
	"github.com/benpate/data/expression"
	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
)

// ScopeFunc is the function signature for a function that can limit database
// queries to a particular "scope".  It inspects the provided context and
// returns criteria that will be passed to all database queries.
type ScopeFunc func(context echo.Context) (expression.Expression, *derp.Error)
