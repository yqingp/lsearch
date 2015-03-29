package lsearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

// if basePath is nil  then use default file path `DefaultStorePath`
func newConfig(basePath string) *config {

	config := &config{}

	if basePath != "" {

		data, err := ioutil.ReadFile(basePath)

		if err != nil {
			Logger.Fatal(err)
		}

		if err := json.Unmarshal(data, config); err != nil {
			Logger.Fatal(err)
		}

		if config.storePath == "" {
			Logger.Fatal(err)
		}
	} else {
		config.storePath = DefaultStorePath
	}

	config.storePath, _ = filepath.Abs(config.storePath)
	fmt.Println(config.storePath)
	return config
}

func (c *config) String() string {
	return fmt.Sprintf("StorePath:[%s]", c.storePath)
}

func (c *config) initStorePath() {
	if err := os.MkdirAll(c.storePath, 0777); err != nil {
		Logger.Fatal(err)
	}
}
