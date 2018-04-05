package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	redirect "github.com/flo80/redirect/pkg/redirect"
	storage "github.com/flo80/redirect/pkg/storage"
	log "github.com/sirupsen/logrus"
)

//Build version (GIT SHA)
var Build = "development"

func main() {
	listenAddress := flag.String("listen", ":8080", "Sets listen address (ip:port) for redirector; empty ip for all interfaces")
	adminAddres := flag.String("admin", "", "Enable a REST API on a specific hostname (listen address has to cover this hostname))")
	redirectFile := flag.String("config", "redirects.json", "Save file for the redirector (loaded at beginning, saved at end)")
	redirectFileIgnoreErr := flag.Bool("ignoreError", false, "Ignore errors when opening redirector file and start with empty redirector")
	redirectNoSave := flag.Bool("noSave", false, "Do not save redirects into redirect file when closing server")
	version := flag.Bool("version", false, "Only print version and quit")
	flag.Parse()

	log.SetLevel(log.DebugLevel)

	if *version {
		fmt.Printf("Version %v \n", Build)
		return
	}

	var server *redirect.Server
	redirector := storage.MapRedirect{}

	if *redirectFile != "" {
		err := mapRedirectorFromFile(*redirectFile, &redirector)
		if err != nil {
			if !*redirectFileIgnoreErr {
				log.Fatalf("Could not create redirector: %v", err)
			}
			redirector = storage.MapRedirect{}
		}
		if !*redirectNoSave {
			defer func() {
				err := SaveMapRedirectorToFile(*redirectFile, &redirector)
				if err != nil {
					log.Fatalf("could not save configuration file: %v", err)
				}
				log.Printf("configuration file %v saved", *redirectFile)
			}()
		}
	}

	if *adminAddres == "" {
		server = redirect.NewServer(*listenAddress, redirect.WithRedirector(&redirector))

	} else {
		server = redirect.NewServer(*listenAddress, redirect.WithAdmin(*adminAddres), redirect.WithRedirector(&redirector))
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
	return
}

//LoadFromFile loads configuration from a file (as json)
func mapRedirectorFromFile(configFile string, redirector *storage.MapRedirect) error {

	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("configuration file %v not found: %v", configFile, err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(redirector)
	if err != nil {
		return fmt.Errorf("could not parse configuration file: %v", err)
	}
	log.Printf("decoded configuration file")

	if err = file.Close(); err != nil {
		return fmt.Errorf("could not close configuration file: %v ", err)
	}
	log.Printf("loaded configuration file %v", configFile)

	return nil
}

//SaveMapRedirectorToFile saves configuration to a file (as json)
func SaveMapRedirectorToFile(configFile string, redirector *storage.MapRedirect) error {
	log.Printf("Trying to save config to file: %v", configFile)

	b, err := json.MarshalIndent(redirector, "", " ")

	if err != nil {
		return fmt.Errorf("could not marshall config file: %v", err)
	}

	err = ioutil.WriteFile(configFile, b, 0644)
	if err != nil {
		return fmt.Errorf("could not write config file: %v", err)
	}
	log.Printf("Config file %v written", configFile)
	return nil
}
