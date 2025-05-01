package middleware

import (
	"github.com/labstack/echo/v4"
	etrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

func Trace(appName string) echo.MiddlewareFunc {
	return etrace.Middleware(etrace.WithServiceName(appName))
}
