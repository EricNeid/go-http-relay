// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: MIT
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/EricNeid/go-http-relay/server"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logFile         = "logs/webserver.log"
	listenAddr      = ":5000"
	destinationAddr = ""
	basePath        = ""
)

func main() {
	// read arguments
	if value, isSet := os.LookupEnv("LISTEN_ADDR"); isSet {
		listenAddr = value
	}
	if value, isSet := os.LookupEnv("DESTINATION_ADDR"); isSet {
		destinationAddr = value
	}
	if value, isSet := os.LookupEnv("BASE_PATH"); isSet {
		basePath = value
	}
	// cli arguments can override environment variables
	flag.StringVar(&listenAddr, "listen-addr", listenAddr, "server listen address")
	flag.StringVar(&destinationAddr, "destination-addr", destinationAddr, "destination address")
	flag.StringVar(&basePath, "base-path", basePath, "base path to serve endpoints")
	flag.Parse()

	// prepare logging
	log.SetPrefix("[APP] ")
	log.SetOutput(
		LazyMultiWriter(
			os.Stdout,
			&lumberjack.Logger{
				Filename:   logFile,
				MaxSize:    500, // megabytes
				MaxBackups: 3,
				MaxAge:     28, // days
			},
		),
	)

	if destinationAddr == "" {
		log.Println("main", "destination not set", "you have to set destination-addr")
		os.Exit(1)
	}

	// prepare graceful shutdown channel
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// create server
	log.Println("main", "creating server")
	srv := server.NewApplicationServer(listenAddr, basePath, destinationAddr)
	go srv.GracefullShutdown(quit, done)

	// start listening
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln("main", "could not start listening", err)
	}

	<-done
	log.Println("main", "Server stopped")
}
