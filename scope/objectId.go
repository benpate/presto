package scope

import (
	"github.com/benpate/data/expression"
	"github.com/benpate/derp"
	"github.com/benpate/presto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectID generates a presto.ScoperFunc using the values provided.  It will attempt to convert every context parameter into a mongodb primitive.ObjectId
func ObjectID(params ...string) presto.ScopeFunc {

	return func(ctx presto.Context) (expression.Expression, *derp.Error) {

		result := expression.And()

		for _, param := range params {

			value := ctx.Param(param)

			objectID, err := primitive.ObjectIDFromHex(value)

			if err != nil {
				return nil, derp.New(500, "scope.ObjectID", "Invalid parameter", param, value, err)
			}

			result = result.And(param, expression.OperatorEqual, objectID)
		}

		return result, nil
	}
}
