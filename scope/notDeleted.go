package scope

import (
	"github.com/benpate/data/expression"
	"github.com/benpate/derp"
	"github.com/benpate/presto"
)

// NotDeleted filters out all records that have not been "virtually deleted" from the database.
func NotDeleted(ctx presto.Context) (expression.Expression, *derp.Error) {
	return expression.New("journal.deleteDate", expression.OperatorEqual, 0), nil
}
