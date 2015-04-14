package config

import (
	"encoding/json"
	"os"
)

const (
	DefaultStorePath = "lsearch_store"
)

type Config struct {
	StorePath string `json:"store_path"`
	IP        string
	Port      string
}

// parse the root dir "config.json" and init base store path
// if basePath is nil  then use default file path `DefaultStorePath`
func NewConfig() *Config {
	configFile, err := os.Open("../config.json")
	if err != nil {
		//Logger.Fatal(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&Config); err != nil {
		//printError("parsing config file", err.Error())
	}
	return nil
}
