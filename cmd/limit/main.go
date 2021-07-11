package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logLevel := envInt("LOG_LEVEL", int(log.InfoLevel))
	log.SetLevel(log.Level(logLevel))

	r := setupRoutes(mux.NewRouter())

	addr := ":" + port
	server := http.Server{
		Addr:    addr,
		Handler: r,

		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Infof("listening on %v", addr)
	log.Fatal(server.ListenAndServe())
}
