package redirectserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//LoadFromFile loads configuration from a file (as json)
func LoadFromFile(configFile string, config *interface{}) error {
	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("configuration file %v not found: %v", configFile, err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return fmt.Errorf("could not parse configuration file: %v", err)
	}
	log.Printf("decoded configuration file into %t", config)

	if err = file.Close(); err != nil {
		return fmt.Errorf("could not close configuration file: %v ", err)
	}
	log.Printf("loaded configuration file %v", configFile)

	return nil
}

//SaveToFile saves configuration to a file (as json)
func SaveToFile(configFile string, config *interface{}) error {
	log.Printf("Trying to save config: %t %v", config, config)

	b, err := json.MarshalIndent(config, "", " ")
	// if _debug {
	// 	var c Configuration
	// 	json.Unmarshal(b, &c)
	// 	log.Debugf("unmarshalled data: %v", c)
	// }

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
