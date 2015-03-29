package index

import (
	"path/filepath"
)

type Index struct {
	indexMeta *IndexMeta
}

func recoverIndex(indexPath string) (*Index, error) {
	index := &Index{}
	index.indexMeta = newIndexMeta()
	index.indexMeta.recoverFromMetaFile(filepath.Join(indexPath, MetaFileName))
	return nil, nil
}
