package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", app.port),
		Handler: app.routes(),
		IdleTimeout: 60 * time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting server on port %d\n", app.port)
	return server.ListenAndServe()
}