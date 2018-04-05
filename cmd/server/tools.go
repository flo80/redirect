package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/flo80/redirect/pkg/storage"
)

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
