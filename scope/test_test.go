package scope

import (
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func getTestContext() echo.Context {

	r := httptest.NewRequest("GET", "http://localhost/example", nil)
	w := httptest.NewRecorder()

	return echo.New().NewContext(r, w)
}
