package index

import (
	"os"
	"path/filepath"
)

type IndexMeta struct {
	file     *os.File
	index    *Index
	fileName string
}

type IndexMetaFieldDump struct {
	id                 int
	name               string
	createdAt          int    `json:created_at`
	fieldType          int    `json:field_type`
	searchAnalyzerName string `json:search_analyzer_name`
	indexAnalyzerName  string `json:index_analyzer_name`
	isIndex            bool   `json:is_index`
}

type IndexMetaDump struct {
	id                  int
	name                string
	defaultAnalyzerName string `json:default_analyzer_name`
	createdAt           int    `json:created_at`
	updatedAt           int    `json:updated_at`
	fields              []IndexMetaFieldDump
}

const (
	MetaFileName = "i_meta"
)

func newIndexMeta(index *Index) *IndexMeta {
	return &IndexMeta{index: index}
}

func (self *IndexMeta) recoverMeta() (bool, error) {
	metaFileName := filepath.Join(self.index.indexPath, MetaFileName)
	if _, err := os.Stat(metaFileName); !os.IsExist(err) {
		return false, nil
	}

	self.fileName = metaFileName

	self.file, _ = os.OpenFile(self.fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)

	return true, nil
}

func (self *IndexMeta) dump() error {

	return nil
}
