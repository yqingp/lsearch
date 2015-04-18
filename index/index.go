package index

import (
    "github.com/yqingp/lsearch/document"
    "github.com/yqingp/lsearch/field"
    "github.com/yqingp/lsearch/mapping"
    "github.com/yqingp/lsearch/store"
    "os"
    "path/filepath"
    "time"
)

type Index struct {
    Id          int
    Name        string
    CreatedAt   int64
    UpdatedAt   int64
    FieldNum    int
    DocumentNum int
    Fields      []*field.Filed
    DB          *store.DB
    Meta        *IndexMeta
    MetaFile    *os.File
    MetaMmap    store.Mmap
}

func New(mapping *mapping.Mapping, baseStorePath string) (*Index, error) {
    storePath, err := initStorePath(baseStorePath, mapping.Name)
    if err != nil {
        return nil, err
    }

    dbStorePath := filepath.Join(storePath, "db")
    db, err := store.Open(dbStorePath, true)
    if err != nil {
        return nil, err
    }

    index := &Index{
        Name:      mapping.Name,
        CreatedAt: time.Now().Unix(),
        UpdatedAt: time.Now().Unix(),
        FieldNum:  len(mapping.Fields),
        Fields:    mapping.Fields,
        DB:        db,
        Meta:      &IndexMeta{},
    }

    if err := index.newIndexMeta(storePath); err != nil {
        return nil, err
    }
    // index.initMeta(storePath)
    return index, nil
}

func (i *Index) IndexDocuments(documents []document.Document) {

}

func (i *Index) internalIndexDocument(document document.Document) {

}

func (i *Index) Remove() {
    i.DB.Close()
}

func initStorePath(baseStorePath string, name string) (string, error) {
    storePath := filepath.Join(baseStorePath, name)
    if err := os.MkdirAll(storePath, 0755); err != nil {
        return storePath, err
    }
    return storePath, nil
}
