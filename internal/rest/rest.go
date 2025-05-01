package rest

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	mdwLib "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	mdw "users-api/internal/rest/middleware"
)

type IUserHandler interface {
	CreateUser(c echo.Context) error
	GetUserByID(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
}

type EchoServer struct {
	e       *echo.Echo
	UserHdl IUserHandler
}

func NewEchoServer(sh IUserHandler) *EchoServer {
	return &EchoServer{
		e:       echo.New(),
		UserHdl: sh,
	}
}

func (es *EchoServer) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Welcome to users-api project.")
}

func (es *EchoServer) RunServer(appName string, timeout int64) *echo.Echo {
	es.e.Use(mdw.Logger())
	es.e.Use(mdw.Trace(appName))
	es.e.Use(mdwLib.Recover())
	es.e.Use(mdwLib.TimeoutWithConfig(mdwLib.TimeoutConfig{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}))

	// monitoring
	es.e.GET("/healthcheck", es.HealthCheck)

	// allow swaggerUI to be served via /swagger/index.html
	es.e.GET("/swagger/*", echoSwagger.WrapHandler)
	es.e.File("/swagger/doc.json", "./docs/swagger.json")

	// symbol
	es.e.POST("/user", es.UserHdl.CreateUser)
	es.e.GET("/user/:id", es.UserHdl.GetUserByID)
	es.e.PATCH("/user/:id", es.UserHdl.UpdateUser)
	es.e.DELETE("/user/:id", es.UserHdl.DeleteUser)

	return es.e
}
