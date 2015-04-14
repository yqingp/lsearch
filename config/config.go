package config

import (
    "encoding/json"
    "fmt"
    . "github.com/yqingp/lsearch/log"
    "io/ioutil"
    "os"
    "path/filepath"
)

const (
    DefaultStorePath = "./lsearch_store"
)

type Config struct {
    StorePath string `json:"store_path"`
}

// if basePath is nil  then use default file path `DefaultStorePath`
func newConfig(basePath string) *Config {

    config := &Config{}

    if basePath != "" {

        data, err := ioutil.ReadFile(basePath)

        if err != nil {
            LSLogger.Fatal(err)
        }

        if err := json.Unmarshal(data, config); err != nil {
            LSLogger.Fatal(err)
        }

        if config.StorePath == "" {
            LSLogger.Fatal(err)
        }
    } else {
        config.StorePath = DefaultStorePath
    }

    config.StorePath, _ = filepath.Abs(config.StorePath)
    fmt.Println(config.StorePath)
    return config
}

func (c *Config) String() string {
    return fmt.Sprintf("StorePath:[%s]", c.StorePath)
}

func (c *Config) initStorePath() {
    if err := os.MkdirAll(c.StorePath, 0777); err != nil {
        LSLogger.Fatal(err)
    }
}
