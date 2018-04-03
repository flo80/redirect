package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	redirect "github.com/flo80/redirect/redirectserver"
)

func main() {
	listenAddress := flag.String("listen", ":8080", "Sets listen address (ip:port) for redirector; empty ip for all interfaces")
	adminAddres := flag.String("admin", "", "Enable a REST API on a specific hostname (listen address has to cover this hostname))")

	// configFilePtr := flag.String("config", "config.json", "filename of config file")
	// noSaveFilePtr := flag.Bool("nosave", false, "Turns on saving of config file")
	flag.Parse()

	var server *redirect.Server
	if *adminAddres == "" {
		server = redirect.NewServer(*listenAddress)

	} else {
		server = redirect.NewServer(*listenAddress, redirect.WithAdmin(*adminAddres))
	}
	server.Redirects.AddRedirect(redirect.Redirect{"sonarr.poetscher.org", "/", "http://leuk.poetscher.org:8989/sonarr/"})

	// if saveConfig {
	// 	defer func() {
	// 		err := config.SaveToFile(configurationFile)
	// 		if err != nil {
	// 			log.Fatalf("could not save configuration file: %v", err)
	// 		}
	// 		log.Printf("configuration file %v saved", configurationFile)
	// 	}()
	// }

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
	return
}
