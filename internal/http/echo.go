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
	//appEchoCtx := &EchoContext{
	//	Context:   nil,
	//	Container: container,
	//}
	//e.Use(echoContextMiddleware(appEchoCtx))
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if _, ok := err.(*echo.HTTPError); !ok {
			container.GetLogger().WithError(err).Error("Error while handling request")
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
	//e.HTTPErrorHandler = HTTPErrorHandler
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

//func echoContextMiddleware(appEchoCtx *EchoContext) echo.MiddlewareFunc {
//	return func(h echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			appEchoCtx.Context = c
//			return h(appEchoCtx)
//		}
//	}
//}
