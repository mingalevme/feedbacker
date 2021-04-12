package http

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mingalevme/feedbacker/internal/app"
)

type EchoHTTPServer struct {
	echo    *echo.Echo
	address string
}

func NewEchoServer(address string, env app.Env) *EchoHTTPServer {
	server := echo.New()
	server.HTTPErrorHandler = func(err error, c echo.Context) {
		if _, ok := err.(*echo.HTTPError); !ok {
			env.Logger().WithRequest(c.Request()).WithError(err).Fatal("Echo: handling request")
		}
		server.DefaultHTTPErrorHandler(err, c)
	}
	server.Use(middleware.Recover())
	prometheus.NewPrometheus("echo", nil).Use(server)
	if env.LogRequests() {
		server.Use(middleware.Logger())
	}
	server.Use(middleware.Secure())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	h := NewEchoHandlerBag(env)
	server.GET("/ping", h.Ping)
	server.GET("/health", h.Health)
	server.POST("/feedback", h.LeaveFeedback)
	server.GET("/feedback/:id", h.ViewFeedback)
	return &EchoHTTPServer{
		echo:    server,
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
