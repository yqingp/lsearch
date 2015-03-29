package index

import (
	"path/filepath"
)

const (
// LockFileName = ""
)

type Index struct {
	indexMeta *IndexMeta
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
