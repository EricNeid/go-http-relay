// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: MIT

// Package server contains webserver with gracefull shutdown functions and logger.
package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const timeFormat string = "Mon Jan 2 15:04:05 2006"

// A ApplicationServer defines parameters and companions for running an HTTP server.
type ApplicationServer struct {
	Webserver *http.Server
}

// NewApplicationServer creates a new instance of ApplicationServer.
func NewApplicationServer(listenAddr string, basePath string) *ApplicationServer {
	log.Println("NewApplicationServer", "creating server", listenAddr, basePath)

	// create webserver
	router := http.NewServeMux()

	// configure routes
	base := normalizePath(basePath)
	router.HandleFunc(base+"/", logCall(welcome))

	return &ApplicationServer{
		Webserver: &http.Server{
			Addr:         listenAddr,
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}

func normalizePath(path string) string {
	if path == "" {
		return path
	}
	// ensure leading slash
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	// remove trailing slash
	path = strings.TrimSuffix(path, "/")
	return path
}

// GracefullShutdown request that the webserver shuts down.
// Server wait for the quit signal and cancels all current connections.
// It pushes true to the done channel after finished shuting down.
func (srv ApplicationServer) GracefullShutdown(quit <-chan os.Signal, done chan<- bool) {
	<-quit
	log.Println("GracefullShutdown", "server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.Webserver.SetKeepAlivesEnabled(false)
	if err := srv.Webserver.Shutdown(ctx); err != nil {
		log.Fatalln("GracefullShutdown", "could not gracefully shutdown the server", err)
	}

	close(done)
}

// ListenAndServe instructs the webserver to listens on
// the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
func (srv ApplicationServer) ListenAndServe() error {
	return srv.Webserver.ListenAndServe()
}

func logCall(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now()
		log.Printf("%s - %s\n", timestamp.Format(timeFormat), r.URL.Path)
		handler(w, r)
	}
}
