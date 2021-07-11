package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func setupRoutes(r *mux.Router) http.Handler {
	r.Use(rateLimiter())

	r.HandleFunc("/echo", echo)
	r.HandleFunc("/admin", admin)

	return r
}

func echo(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "hello %v\n%v\n", req.RemoteAddr, req.Header)
}

func admin(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusUnauthorized)
}
