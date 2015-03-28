package lsearch

// "github.com/yqingp/lsearch/"

import (
	"fmt"
)

type LSearch struct {
	config   *config
	indexer  *indexer
	searcher *searcher
}

func NewLSearch(configFilePath string) *LSearch {
	return &LSearch{
		config:   newConfig(configFilePath),
		indexer:  nil,
		searcher: nil,
	}
}

func (lsearch *LSearch) Init() {
	fmt.Println("init")
	lsearch.config.initStorePath()
}
