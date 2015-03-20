package lsearch

import (
	"fmt"
	"testing"
)

func TestConfigInit(t *testing.T) {
	config, err := InitConfig("")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(config)
	}
}
