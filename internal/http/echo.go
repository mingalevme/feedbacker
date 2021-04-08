package http

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mingalevme/feedbacker/internal/app/di"
)

type EchoContext struct {
	echo.Context
	di.Container
}

type EchoHTTPServer struct {
	echo *echo.Echo
	address string
}

func NewEchoServer(address string, container di.Container) *EchoHTTPServer {
	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if _, ok := err.(*echo.HTTPError); !ok {
			container.GetLogger().WithRequest(c.Request()).WithError(err).Fatal("Echo: handling request")
		}
		e.DefaultHTTPErrorHandler(err, c)
	}
	e.Use(middleware.Recover())
	prometheus.NewPrometheus("echo", nil).Use(e)
	e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	h := NewEchoHandlerBag(container)
	e.GET("/ping", h.Ping)
	e.POST("/feedback", h.LeaveFeedback)
	e.GET("/feedback/:id", h.ViewFeedback)
	return &EchoHTTPServer{
		echo:    e,
		address: address,
	}
}

func (s *EchoHTTPServer) GetAddr() string {
	return s.address
}

func (s *EchoHTTPServer) ListenAndServe() error {
	return s.echo.Start(s.address)
}

func (s *EchoHTTPServer) WaitForShutdown() {
	waitForShutdown(s.echo.Shutdown)
}
