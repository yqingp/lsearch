package lsearch

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

const (
	DefaultStorePath = "./lsearch_store"
)

var (
	ERROR_CONFIG_Open_File_Error  = errors.New("open config file error")
	ERROR_CONFIG_Json_Parse_Error = errors.New("json parse config file error")
	ERROR_CONFIG_Dbpath_Not_Set   = errors.New("dbpath not set")
)

type Config struct {
	StorePath string `json:"store_path"`
}

func InitConfig(filepath string) (*Config, error) {

	config := &Config{}

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

// func (c *Config) InitStorePath() {

// }
