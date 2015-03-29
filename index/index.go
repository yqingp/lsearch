package index

import (
	"github.com/yqingp/lsearch/analyzer"
	"github.com/yqingp/lsearch/field"
	"path/filepath"
)

const (
// LockFileName = ""
)

type Index struct {
	indexMeta       *IndexMeta
	defaultAnalyzer *analyzer.Analyzer
	fields          []*field.Filed
	fieldNum        int
	documentNum     int
}

func recoverIndex(indexPath string) (*Index, error) {
	index := &Index{
		indexMeta: newIndexMeta(),
	}

	isExistMeta, err := index.indexMeta.recoverFromMetaFile(filepath.Join(indexPath, MetaFileName))
	if !isExistMeta {
		return nil, nil
	}

	return nil, err
}
