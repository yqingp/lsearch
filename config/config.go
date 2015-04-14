package config

import ()

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

    return nil
}
