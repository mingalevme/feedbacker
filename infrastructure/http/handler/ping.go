package handler

import netHTTP "net/http"

func Ping(w netHTTP.ResponseWriter, r *netHTTP.Request) {
	_, _ = w.Write([]byte("Pong"))
}
