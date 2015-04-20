package index

import (
    "encoding/json"
    "errors"
    "github.com/yqingp/lsearch/analyzer"
    "github.com/yqingp/lsearch/document"
    "github.com/yqingp/lsearch/mapping"
    "github.com/yqingp/lsearch/store"
    "github.com/yqingp/lsearch/util"
    "os"
    "path/filepath"
    "sort"
    // "time"
)

type Index struct {
    Id          int
    Name        string
    DocumentNum int
    DB          *store.DB
    DocumentDB  *store.DB
    Meta        *IndexMeta
    Analyzer    *analyzer.Analyzer
}

func New(mapping *mapping.Mapping, baseStorePath string) *Index {
    storePath := initStorePath(baseStorePath, mapping.Name)
    dbStorePath := filepath.Join(storePath, "db")
    db, err := store.Open(dbStorePath, true)
    if err != nil {
        panic(err)
    }

    documentStorePath := filepath.Join(storePath, "documents")
    documentDb, err := store.Open(documentStorePath, false)
    if err != nil {
        panic(err)
    }

    index := &Index{
        Name:       mapping.Name,
        DB:         db,
        Meta:       &IndexMeta{},
        DocumentDB: documentDb,
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

        documentStorePath := filepath.Join(baseStorePath, "documents")
        documentDb, err := store.Open(documentStorePath, false)
        if err != nil {
            panic(err)
        }

        index := &Index{
            Name:       name,
            DB:         db,
            DocumentDB: documentDb,
        }

        storePath := filepath.Join(baseStorePath, name)
        index.Meta = recoverIndexMeta(storePath)
        indexes[name] = index
    }

    return indexes
}

func (i *Index) AddDocuments(documents []document.Document) (interface{}, error) {
    if !i.validateDocuments(documents) {
        return nil, errors.New("documents structure error")
    }

    for _, doc := range documents {
        i.internalAddDocument(doc)
    }

    return nil, nil
}

func (i *Index) internalAddDocument(document document.Document) {
    document.Analyze(i.Analyzer)

    id := document.Id()
    data, err := document.Encode()
    if err != nil {
        panic(err)
    }

    internalId, err := i.DocumentDB.Set(-1, []byte(id), data)
    if err != nil {
        panic(err)
    }

    for k, _ := range document.Tokens() {
        data, _ := i.DB.Get([]byte(k)) // fix
        postings := make(util.Posting, 1000)
        json.Unmarshal(data, postings)
        postings = append(postings, internalId)
        sort.Sort(util.Posting(postings))

        data, err = json.Marshal(postings)

        i.DB.Set(-1, []byte(k), []byte(data))
    }
}

func (i *Index) DeleteDocuments(documents []document.Document) {

}

func (i *Index) UpdateDocuments(documents []document.Document) {

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

func (i *Index) validateDocuments(documents []document.Document) bool {
    if len(documents) < 1 {
        return false
    }

    if !document.Validate(documents, i.Meta.Fields) {
        return false
    }

    return true
}
