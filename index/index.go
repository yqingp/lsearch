package index

import (
    "github.com/yqingp/lsearch/document"
    "github.com/yqingp/lsearch/mapping"
    "github.com/yqingp/lsearch/store"
    "os"
    "path/filepath"
    // "time"
)

type Index struct {
    Id          int
    Name        string
    DocumentNum int
    DB          *store.DB
    Meta        *IndexMeta
}

func New(mapping *mapping.Mapping, baseStorePath string) *Index {
    storePath := initStorePath(baseStorePath, mapping.Name)
    dbStorePath := filepath.Join(storePath, "db")
    db, err := store.Open(dbStorePath, true)
    if err != nil {
        panic(err)
    }

    index := &Index{
        Name: mapping.Name,
        DB:   db,
        Meta: &IndexMeta{},
    }

    index.Meta = newMeta(storePath, mapping)

    return index
}

func (i *Index) IndexDocuments(documents []document.Document) {

}

func (i *Index) internalIndexDocument(document document.Document) {

}

func (i *Index) Remove() {
    i.DB.Close()
}

func initStorePath(baseStorePath string, name string) string {
    storePath := filepath.Join(baseStorePath, name)
    if err := os.MkdirAll(storePath, 0755); err != nil {
        panic(err)
    }
    return storePath
}
