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
	Ip               = "127.0.0.1"
	Port             = "9919"
)

// define default ip and port

type Config struct {
	StorePath string `json:"store_path"`
	Ip        string `json:"ip"`
	Port      string `json:"port"`
	Path      string
}

// parse the root dir "config.json" and init base store path
// using default config "store" or "ip" or "port" if config.json not exist
func New() *Config {
	config := &Config{}
	content, err := ioutil.ReadFile("config.json")
	if err == nil {
		err = json.Unmarshal(content, config)
		if err != nil {
			log.Fatal(err)
		}
	}

	if config.StorePath == "" {
		config.StorePath = DefaultStorePath
	}

	if config.Ip == "" {
		config.Ip = Ip
	}

	if config.Port == "" {
		config.Ip = Ip
	}

	config.initStorePath()

	log.Println(config)
	return config
}

func (c *Config) BindAddr() string {
	return c.Ip + ":" + c.Port
}

func (c *Config) initStorePath() {
	absPath, err := filepath.Abs(c.StorePath)
	if err != nil {
		log.Fatal(err)
	}
	c.Path = absPath

	if err := os.MkdirAll(c.Path, 0755); err != nil {
		log.Fatal(err)
	}
}
