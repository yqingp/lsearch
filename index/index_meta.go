package index

import (
	// "github.com/yqingp/lsearch/analyzer"
	// "github.com/yqingp/lsearch/field"
	"encoding/json"
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
	createdAt          int64  `json:created_at`
	fieldType          int    `json:field_type`
	searchAnalyzerName string `json:search_analyzer_name`
	indexAnalyzerName  string `json:index_analyzer_name`
	isIndex            bool   `json:is_index`
}

type IndexMetaDump struct {
	id                  int
	name                string
	defaultAnalyzerName string `json:default_analyzer_name`
	createdAt           int64  `json:created_at`
	updatedAt           int64  `json:updated_at`
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
	if self.index == nil {
		return nil
	}
	indexMetaDump := IndexMetaDump{}
	indexMetaDump.id = self.index.id
	indexMetaDump.createdAt = self.index.createdAt
	indexMetaDump.updatedAt = self.index.updatedAt
	indexMetaDump.name = self.index.name
	indexMetaDump.defaultAnalyzerName = self.index.defaultAnalyzer.Name
	indexMetaDump.fields = []IndexMetaFieldDump{}

	for _, v := range self.index.fields {
		idf := IndexMetaFieldDump{}
		idf.name = v.Name
		idf.createdAt = v.CreatedAt
		idf.fieldType = int(v.FieldType)
		idf.indexAnalyzerName = v.IndexAnalyzer.Name
		idf.searchAnalyzerName = v.SearchAnalyzer.Name
		idf.isIndex = v.IsIndex
		indexMetaDump.fields = append(indexMetaDump.fields, idf)
	}

	t, _ := json.Marshal(indexMetaDump)
	self.file.Write(t)
	return nil
}
