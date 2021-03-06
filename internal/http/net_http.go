package http

import (
	"context"
	netHTTP "net/http"
	"time"
)

type NetHTTPServer struct {
	netHTTP.Server
}

func NewNetHTTPServer(address string, router netHTTP.Handler) *NetHTTPServer {
	return &NetHTTPServer{
		netHTTP.Server{
			Handler:      router,
			Addr:         address,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *NetHTTPServer) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
