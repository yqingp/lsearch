package index

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"github.com/yqingp/lsearch/analyzer"
	"github.com/yqingp/lsearch/document"
	"github.com/yqingp/lsearch/mapping"
	"github.com/yqingp/lsearch/store"
	"github.com/yqingp/lsearch/util"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

const (
	IndexRWMutextCount     = 256
	IndexPostingMutexCount = 256
)

type Index struct {
	Id           int
	Name         string
	DocumentNum  int
	DB           *store.DB
	DocumentDB   *store.DB
	Meta         *Meta
	Analyzer     *analyzer.Analyzer
	idMutex      [IndexRWMutextCount]*sync.Mutex
	postingMutex [IndexPostingMutexCount]*sync.Mutex
}

var Logger *log.Logger = log.New(os.Stdout, "[DEGUG]", log.Lshortfile|log.Ldate|log.Ltime)

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
		Meta:       &Meta{},
		DocumentDB: documentDb,
	}

	for i := 0; i < IndexRWMutextCount; i++ {
		index.idMutex[i] = &sync.Mutex{}
	}

	for i := 0; i < IndexPostingMutexCount; i++ {
		index.postingMutex[i] = &sync.Mutex{}
	}

	index.Meta = newMeta(storePath, mapping)

	return index
}

func Recover(baseStorePath string) map[string]*Index {
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

		documentStorePath := filepath.Join(baseStorePath, name, "documents")

		documentDb, err := store.Open(documentStorePath, false)
		if err != nil {
			panic(err)
		}

		index := &Index{
			Name:       name,
			DB:         db,
			DocumentDB: documentDb,
		}

		for i := 0; i < IndexRWMutextCount; i++ {
			index.idMutex[i] = &sync.Mutex{}
		}

		for i := 0; i < IndexPostingMutexCount; i++ {
			index.postingMutex[i] = &sync.Mutex{}
		}

		storePath := filepath.Join(baseStorePath, name)
		index.Meta = recoverMeta(storePath)
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

	return "done", nil
}

func (i *Index) internalAddDocument(doc document.Document) {
	id := doc.Id

	i.lockId([]byte(id))
	defer i.unlockId([]byte(id))

	doc.InitTokens()
	doc.Analyze(i.Analyzer)

	data, err := doc.Encode()
	if err != nil {
		panic(err)
	}

	md5Val := md5.Sum(data)
	exist := false

	oldData, internalId := i.DocumentDB.GetAndReturnInternalId([]byte(id))
	if internalId > 0 {
		exist = true
	}

	oldDoc := &document.Document{}

	oldTokens := oldDoc.Tokens()

	if exist {
		oldMd5Val := md5.Sum(oldData)

		if oldMd5Val == md5Val {
			Logger.Println("==== md5 equal, same document")
			return
		}

		if err := json.Unmarshal(oldData, oldDoc); err != nil {
			panic(err)
		}

		oldTokens = oldDoc.Tokens()
	}

	internalId, err = i.DocumentDB.Set(-1, []byte(id), data)

	if err != nil {
		panic(err)
	}

	delTokens, addTokens := CheckTokensAndSplit(oldTokens, doc.Tokens())

	for k, _ := range delTokens {
		i.lockToken([]byte(k))
		data, _ := i.DB.Get([]byte(k)) // fix

		postings := make(util.Posting, 1000)
		json.Unmarshal(data, &postings)

		pos := sort.Search(len(postings), func(i int) bool {
			return postings[i] == internalId
		})

		postings = append(postings[:pos], postings[(pos+1):]...)
		data, _ = json.Marshal(postings)

		i.DB.Set(-1, []byte(k), data)

		i.unlockToken([]byte(k))
	}

	for k, _ := range addTokens {
		i.lockToken([]byte(k))

		data, ret := i.DB.Get([]byte(k)) // fix

		postings := make(util.Posting, 1000)

		if ret == 0 {
			json.Unmarshal(data, &postings)
		}

		pos := 0
		if postings[0] > 0 {
			pos = sort.Search(len(postings), func(i int) bool {
				return postings[i] >= internalId
			})
			copy(postings[pos+1:], postings[pos:])
		}

		postings[pos] = internalId

		data, _ = json.Marshal(postings)

		i.DB.Set(-1, []byte(k), data)

		i.unlockToken([]byte(k))
	}

}

func (i *Index) internalUpdateDocument() {

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

func (i *Index) View() *Meta {
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

func CheckTokensAndSplit(oldTokens map[string]string, newTokens map[string]string) (del map[string]string, add map[string]string) {
	if len(oldTokens) == 0 {
		add = newTokens
		return
	}

	if len(newTokens) == 0 {
		return
	}

	for k, _ := range oldTokens {
		if _, ok := newTokens[k]; !ok {
			del[k] = ""
		}
	}

	for k, _ := range newTokens {
		if _, ok := oldTokens[k]; !ok {
			add[k] = ""
		}
	}

	return
}

func (i *Index) lockToken(token []byte) {
	i.postingMutex[util.MurmurHash3(token)%IndexPostingMutexCount].Lock()
}

func (i *Index) unlockToken(token []byte) {
	i.postingMutex[util.MurmurHash3(token)%IndexPostingMutexCount].Unlock()
}

func (i *Index) lockId(id []byte) {
	i.idMutex[util.MurmurHash3(id)%IndexRWMutextCount].Lock()
}

func (i *Index) unlockId(id []byte) {
	i.idMutex[util.MurmurHash3(id)%IndexRWMutextCount].Unlock()
}
