package lsearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	DefaultStorePath = "./lsearch_store"
)

// var (
// 	ERROR_CONFIG_OPEN_FILE       = errors.New("open config file error")
// 	ERROR_CONFIG_JSON_PARSE      = errors.New("json parse config file error")
// 	ERROR_CONFIG_STORE_PATH      = errors.New("store path not set")
// 	ERROR_CONGIF_MAKE_STORE_PATH = errors.New("mkdir store path  error")
// )

type config struct {
	storePath string `json:"store_path"`
}

// if filepath is nil  then use default file path `DefaultStorePath`
func newConfig(filepath string) *config {

	config := &config{}

	if filepath != "" {

		data, err := ioutil.ReadFile(filepath)

		if err != nil {
			logger.Fatal(err)
		}

		if err := json.Unmarshal(data, config); err != nil {
			logger.Fatal(err)
		}

		if config.storePath == "" {
			logger.Fatal(err)
		}

		return config
	}
	config.storePath = DefaultStorePath

	return config
}

func (c *config) String() string {
	return fmt.Sprintf("StorePath:[%s]", c.storePath)
}

func (c *config) initStorePath() {
	if err := os.MkdirAll(c.storePath, 0777); err != nil {
		logger.Fatal(err)
	}
}
