package main

import (
	"log"
	"os"
	"os/signal"

	redirect "github.com/flo80/redirect/redirectserver"
)

var _debug = false

func main() {
	// debugPtr := flag.Bool("debug", false, "Turn debugging log on")
	// configFilePtr := flag.String("config", "config.json", "filename of config file")
	// noSaveFilePtr := flag.Bool("nosave", false, "Turns on saving of config file")
	// flag.Parse()
	// _debug = *debugPtr

	server := redirect.NewServer(":8080", "localhost:8080")

	server.Redirects.AddRedirect("sonarr.poetscher.org", "/", "http://leuk.poetscher.org:8989/sonarr/")

	json, err := server.Redirects.GetJSON()
	if err != nil {
		log.Fatalf("cannot get JSON: %v", err)
	}
	log.Printf("JSON %s", json)

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
