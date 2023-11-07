package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	log "github.com/sirupsen/logrus"
	"gopkg.in/tylerb/graceful.v1"
	"netflix-cache/resources/cache"
	"netflix-cache/server"
	"os"
	"time"
)

const (
	PortRangeMin = 1024
	PortRangeMax = 65535
)

func main() {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetLevel(log.InfoLevel)

	// Build command line parser, port only option, defaults to 8080 rather than fail if not specified
	parser := argparse.NewParser("netflix-cache-parser", "Command line parser for netflix-cache")
	port := parser.Int("p", "port", &argparse.Options{Required: false, Help: "Supply the port for the service to run on", Default: 8080})

	// if the user tried to give us input and we failed to parse don't start on default port, let them know
	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatalf("Failed to parse command line arguments: %s", err)
	}

	// Port range sanity check
	if !(*port >= PortRangeMin && *port <= PortRangeMax) {
		log.Fatalf("Port value supplied out of user application range %v-%v", PortRangeMin, PortRangeMax)
	}

	// init cache, pass pointer into server definition
	orgCache, err := cache.New()
	if err != nil {
		log.Warn("failed to initialize cache! fix me!", err)
	}

	// Starting server with graceful restarts, throw fatal log if startup fails
	log.Info(fmt.Sprintf("Starting netflix-cache service on port %d", *port))
	log.Fatal(graceful.RunWithErr(fmt.Sprintf(":%v", *port), time.Second, server.New(orgCache)))
}
