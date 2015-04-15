package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	DefaultStorePath = "lsearch_store"
	Ip               = "127.0.0.1"
	Port             = "9919"
)

// define default ip and port

type Config struct {
	StorePath string `json:"store_path"`
	Ip        string
	Port      string
}

// parse the root dir "config.json" and init base store path
// if basePath is nil  then use default file path `DefaultStorePath`
// using default config "store" or "ip" or "port" if config.json not exist
func NewConfig() *Config {
	str := `{"store_path":"` + DefaultStorePath + `",` + `"ip":"` + Ip + `",` + `"port":"` + Port + `"` + `}`
	config := &Config{}
	json.Unmarshal([]byte(str), config)
	//log.Printf("%v", str)
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(content, config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v", config.Port)
	return config
}
