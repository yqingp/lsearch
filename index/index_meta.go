package index

import (
	"os"
)

type IndexMeta struct {
}

const (
	MetaFileName = "i_meta"
)

func newIndexMeta() *IndexMeta {
	return &IndexMeta{}
}

func (i *IndexMeta) recoverFromMetaFile(metaFilePath string) (bool, error) {
	if _, err := os.Stat(metaFilePath); !os.IsExist(err) {
		return false, nil
	}

	return true, nil
}
