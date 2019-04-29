package presto

import "github.com/labstack/echo"
import "github.com/benpate/derp"

type ScopeFunc func(context echo.Context) (map[string]interface{}, *derp.Error)

func IDScope(context echo.Context) (map[string]interface{}, *derp.Error) {

	names := context.ParamNames()
	result := map[string]interface{}{}

	for index, param := range names {

		// If "ID" is one of the route parameters, then look for a non-empty value
		if param == "id" {
			values := context.ParamValues()
			if value := values[index]; value != "" {
				result["_id"] = value
				return result, nil
			}

			return result, derp.New(derp.CodeBadRequestError, "presto.scope.ID", "Invalid Object ID - Cannot be empty")
		}
	}

	// Otherwise, scan all items.
	return result, nil
}
