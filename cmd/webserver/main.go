// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: MIT
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/EricNeid/go-webserver/server"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logFile    string = "logs/webserver.log"
	listenAddr string = ":5000"
	basePath   string = ""
)

func main() {
	// read arguments
	if value, isSet := os.LookupEnv("LISTEN_ADDR"); isSet {
		listenAddr = value
	}
	if value, isSet := os.LookupEnv("BASE_PATH"); isSet {
		basePath = value
	}
	// cli arguments can override environment variables
	flag.StringVar(&listenAddr, "listen-addr", listenAddr, "server listen address")
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
				MaxAge:     28, //days
			},
		),
	)

	// prepare gracefull shutdown channel
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// create server
	log.Println("main", "creating server")
	server := server.NewApplicationServer(listenAddr, basePath)
	go server.GracefullShutdown(quit, done)

	// start listening
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln("main", "could not start listening", err)
	}

	<-done
	log.Println("main", "Server stopped")
}
