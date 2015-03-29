package index

import (
	"fmt"
)

type IndexMeta struct {
}

const (
	MetaFileName = "i_meta"
)

func newIndexMeta() *IndexMeta {
	return &IndexMeta{}
}

func (i *IndexMeta) recoverFromMetaFile(metaFilePath string) {
	fmt.Println(metaFilePath)
}
