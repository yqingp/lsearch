package lsearch

// "github.com/yqingp/lsearch/"

type LSearch struct {
	config   *config
	indexer  *indexer
	searcher *searcher
}

func NewLSearch(filepath string) *LSearch {
	return &LSearch{
		dbpath:   "",
		indexer:  nil,
		searcher: nil,
	}
}
