package index

import (
    "github.com/yqingp/lsearch/analyzer"
    "github.com/yqingp/lsearch/document"
    "github.com/yqingp/lsearch/field"
    "github.com/yqingp/lsearch/store"
)

type Index struct {
    Id              int
    Name            string
    CreatedAt       int64
    UpdatedAt       int64
    DefaultAnalyzer *analyzer.Analyzer
    Fields          []*field.Filed
    FieldNum        int
    DocumentNum     int
    DB              *store.DB
}

func (i *Index) IndexDocuments(documents []document.Document) {

}

func (i *Index) internalIndexDocument(document document.Document) {

}
