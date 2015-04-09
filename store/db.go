package store

import (
    "errors"
    "github.com/yqingp/lsearch/util"
    "log"
    "os"
    "sync"
)

type Db struct {
    status         int
    blockMax       int
    mmTotal        int64
    xxTotal        int64
    mutex          *sync.Mutex
    freeBlockMutex *sync.Mutex
    indexMutex     *sync.Mutex
    blockMutex     *sync.Mutex
    state          *DbState
    stateIO        DbIO
    blockQueueIO   DbIO
    blockQueues    []DbBlockQueue
    indexIO        DbIO
    indexes        []DbIndex
    dbsIO          [DB_MFILE_MAX]DbIO
    blocks         [DB_XBLOCKS_MAX]DbBlock
    basedir        string
    kmap           *util.Mmtrie
    loggerFile     *os.File
    logger         *log.Logger
    isMmap         bool
    mutexs         [DB_MUTEX_MAX]*sync.Mutex
}

func Open(basedir string, isMmap bool) (*Db, error) {
    if basedir == "" {
        return nil, errors.New("basedir is blank")
    }

    db := &Db{}

    db.freeBlockMutex = &sync.Mutex{}
    db.indexMutex = &sync.Mutex{}
    db.blockMutex = &sync.Mutex{}
    db.mutex = &sync.Mutex{}
    db.basedir = basedir
    db.isMmap = isMmap

    if err := db.initKmap(); err != nil {
        return nil, err
    }

    if err := db.initLogger(); err != nil {
        return nil, err
    }

    return db, nil
}

// if id is < 1  generate auto increment id
func (self *Db) Set(id int, key []byte, value []byte) (int, error) {
    if key == nil || value == nil || len(value) == 0 {
        return -1, errors.New("key or value is blank")
    }

    var err error

    if id < 1 {
        id, err = self.kmap.Set(key)
        if err == nil {
            return -1, err
        }
    }

    self.set(id, value)

    return id, nil
}

func (self *Db) set(id int, value []byte) int {

    ret := -1
    dbIndexes := self.indexes
    if self.status != 0 || dbIndexes == nil {
        return ret
    }

    valueLen := len(value)

    blocksCountNum := 0

    _ = blocksCountNum
    self.indexMutex.Lock()
    self.checkIndexIOWithId(id)
    self.indexMutex.Unlock()

    self.lockId(id)
    freeQueue, oldFreeQueue := DbBlockQueue{}, DbBlockQueue{}
    _ = freeQueue
    if dbIndexes[id].blockSize < valueLen {
        if dbIndexes[id].blockSize > 0 {
            oldFreeQueue.index = dbIndexes[id].index
            oldFreeQueue.blockId = dbIndexes[id].blockId
            oldFreeQueue.count = blocksCount(dbIndexes[id].blockSize)
            dbIndexes[id].blockSize = 0
            dbIndexes[id].blockId = 0
            dbIndexes[id].ndata = 0
        }

        blocksCountNum = blocksCount(valueLen)
    }

    self.unlockId(id)

    return ret
}

func (self *Db) lockId(id int) {
    self.mutexs[id%DB_MUTEX_MAX].Lock()
}

func (self *Db) unlockId(id int) {
    self.mutexs[id%DB_MUTEX_MAX].Unlock()
}

func (self *Db) Get() {

}

func (self *Db) Del() {

}
