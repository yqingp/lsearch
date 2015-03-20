package lsearch

import (
// "github.com/yqingp/lsearch/"
)

type LSearch struct {
	dbpath   string
	indexer  *Indexer
	searcher *Searcher
}

func NewLSearch() *LSearch {
	return &LSearch{
		dbpath:   "",
		indexer:  nil,
		searcher: nil,
	}
}
