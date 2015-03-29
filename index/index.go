package index

import (
	"github.com/yqingp/lsearch/analyzer"
	"github.com/yqingp/lsearch/field"
)

const (
// LockFileName = ""
)

type Index struct {
	id              int
	name            string
	createdAt       int64
	updatedAt       int64
	indexMeta       *IndexMeta
	defaultAnalyzer *analyzer.Analyzer
	fields          []*field.Filed
	fieldNum        int
	documentNum     int
	indexPath       string
}

func recoverIndex(indexPath string) (*Index, error) {
	index := &Index{}
	index.indexMeta = newIndexMeta(index)
	index.indexPath = indexPath
	isExistMeta, err := index.indexMeta.recoverMeta()
	if !isExistMeta {
		return nil, nil
	}

	return nil, err
}
