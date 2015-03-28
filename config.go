package lsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	DefaultStorePath = "./lsearch_store"
)

var (
	ERROR_CONFIG_Open_File_Error  = errors.New("open config file error")
	ERROR_CONFIG_Json_Parse_Error = errors.New("json parse config file error")
	ERROR_CONFIG_Dbpath_Not_Set   = errors.New("dbpath not set")
)

type config struct {
	storePath string `json:"store_path"`
}

// if filepath is nil  then use default file path `DefaultStorePath`
func InitConfig(filepath string) (*config, error) {

	config := &config{}

	if filepath != "" {

		data, err := ioutil.ReadFile(filepath)

		if err != nil {
			return nil, ERROR_CONFIG_Open_File_Error
		}

		if err := json.Unmarshal(data, config); err != nil {
			return nil, ERROR_CONFIG_Json_Parse_Error
		}

		if config.StorePath == "" {
			return nil, ERROR_CONFIG_Dbpath_Not_Set
		}

		return config, nil
	}
	config.StorePath = DefaultStorePath

	return config, nil
}

func (c *config) String() string {
	return fmt.Sprintf("StorePath:[%s]", c.StorePath)
}

func (c *config) InitStorePath() {
	if err := os.MkdirAll(c.StorePath, 0777); err != nil {
		fmt.Println(err)
		return
	}
}
