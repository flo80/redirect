package main

import (
	"os"
	"os/signal"

	redirect "github.com/flo80/redirect/pkg/redirect"
	storage "github.com/flo80/redirect/pkg/storage"
	log "github.com/sirupsen/logrus"
)

//Build version (GIT SHA)
var Build = "development"

// server settings
var config struct {
	listenAddress         string
	adminAddress          string
	redirectFile          string
	redirectFileIgnoreErr bool
	redirectNoSave        bool
	debug                 bool
}

func main() {
	Execute()
}

func runServer() {
	if config.debug {
		log.SetLevel(log.DebugLevel)
	}

	var server *redirect.Server
	redirector := storage.MapRedirect{}

	if config.redirectFile != "" {
		err := mapRedirectorFromFile(config.redirectFile, &redirector)
		if err != nil {
			if !config.redirectFileIgnoreErr {
				log.Fatalf("Could not create redirector: %v", err)
			}
			redirector = storage.MapRedirect{}
		}
		if !config.redirectNoSave {
			defer func() {
				err := SaveMapRedirectorToFile(config.redirectFile, &redirector)
				if err != nil {
					log.Fatalf("could not save redirector to file: %v", err)
				}
				log.Printf("redirector configuration file %v saved", config.redirectFile)
			}()
		}
	}

	if config.adminAddress == "" {
		server = redirect.NewServer(config.listenAddress, redirect.WithRedirector(&redirector))

	} else {
		server = redirect.NewServer(config.listenAddress, redirect.WithAdmin(config.adminAddress), redirect.WithRedirector(&redirector))
	}

	go func() {
		err := server.StartServer()
		if err != nil {
			log.Fatalf("could not start server: %v", err)
		}
	}()

	log.Printf("server started, waiting for interrupt")

	// wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	_ = <-c
	log.Printf("server stopped")
}
