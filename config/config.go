package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	config_st := &Config{}
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	//json.Unmarshal(, v)
	err = json.Unmarshal(content, config_st)
	if err != nil {
		log.Fatal(err)
	}
	filepath.Abs(os.Args[0])
	log.Printf("%v", config_st.Port)

	return config_st
}
