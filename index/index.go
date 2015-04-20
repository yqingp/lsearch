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

    index.Meta = newIndexMeta(storePath, mapping)

    return index
}

func RecoverIndexes(baseStorePath string) map[string]*Index {
    indexes := map[string]*Index{}

    dir, err := os.OpenFile(baseStorePath, os.O_RDONLY, 0755)
    if err != nil {
        panic(err)
    }
    dirs, err := dir.Readdir(0)
    if err != nil {
        panic(err)
    }
    for _, file := range dirs {
        if !file.IsDir() {
            continue
        }
        name := file.Name()
        dbStorePath := filepath.Join(baseStorePath, name, "db")
        db, err := store.Open(dbStorePath, true)
        if err != nil {
            panic(err)
        }

        index := &Index{
            Name: name,
            DB:   db,
        }
        storePath := filepath.Join(baseStorePath, name)
        index.Meta = recoverIndexMeta(storePath)
        indexes[name] = index
    }

    return indexes
}

func (i *Index) IndexDocuments(documents []document.Document) {

}

func (i *Index) internalIndexDocument(document document.Document) {

}

func (i *Index) Remove() {
    i.DB.Close()
    if i.Meta != nil {
        os.RemoveAll(i.Meta.StorePath)
    }
}

func (i *Index) View() *IndexMeta {
    return i.Meta
}

func initStorePath(baseStorePath string, name string) string {
    storePath := filepath.Join(baseStorePath, name)
    if err := os.MkdirAll(storePath, 0755); err != nil {
        panic(err)
    }
    return storePath
}
