package store

import (
    "errors"
    "github.com/yqingp/lsearch/util"
    "log"
    "os"
    "sync"
    "time"
)

type Db struct {
    status          int
    blockMax        int
    mmTotal         int64
    xxTotal         int64
    mutex           *sync.Mutex
    blockQueueMutex *sync.Mutex
    indexMutex      *sync.Mutex
    blockMutex      *sync.Mutex
    state           *DbState
    stateIO         DbIO
    blockQueueIO    DbIO
    blockQueues     []DbBlockQueue
    indexIO         DbIO
    indexes         []DbIndex
    dbsIO           [DB_MFILE_MAX]DbIO
    blocks          [DB_XBLOCKS_MAX]DbBlock
    basedir         string
    kmap            *util.Mmtrie
    loggerFile      *os.File
    logger          *log.Logger
    isMmap          bool
    mutexs          [DB_MUTEX_MAX]*sync.Mutex
}

func Open(basedir string, isMmap bool) (*Db, error) {
    if basedir == "" {
        return nil, errors.New("basedir is blank")
    }

    db := &Db{}

    db.blockQueueMutex = &sync.Mutex{}
    db.indexMutex = &sync.Mutex{}
    db.blockMutex = &sync.Mutex{}
    db.mutex = &sync.Mutex{}
    db.basedir = basedir
    db.isMmap = isMmap

    if err := db.initDir(); err != nil {
        return nil, err
    }

    if err := db.initLogger(); err != nil {
        return nil, err
    }

    if err := db.initKmap(); err != nil {

    }

    db.initState()
    db.initBlockQueue()
    db.initIndex()
    db.initDbsIO()

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
        if err != nil {
            self.logger.Println(err)
            return -1, err
        }
    }

    ret := self.internalSet(id, value)
    if ret == -1 {
        self.logger.Fatal("set error")
    }

    return id, nil
}

func (self *Db) internalSet(id int, value []byte) int {
    self.logger.Println(string(value))
    ret := -1
    dbIndexes := self.indexes
    if self.status != 0 || dbIndexes == nil {
        return ret
    }

    valueLen := len(value)

    blocksCountNum := 0

    index := 0
    _ = blocksCountNum
    self.indexMutex.Lock()
    self.checkIndexIOWithId(id)
    self.indexMutex.Unlock()

    self.logger.Println(id)
    self.lockId(id)
    defer self.unlockId(id)
    link, old := &DbBlockQueue{}, &DbBlockQueue{}

    if dbIndexes[id].blockSize < valueLen {
        self.logger.Println(dbIndexes[id].blockSize)
        if dbIndexes[id].blockSize > 0 {
            old.index = dbIndexes[id].index
            old.blockId = dbIndexes[id].blockId
            old.count = blocksCount(dbIndexes[id].blockSize)
            dbIndexes[id].blockSize = 0
            dbIndexes[id].blockId = 0
            dbIndexes[id].ndata = 0
        }

        blocksCountNum = blocksCount(valueLen)
        self.logger.Println("blocksCountNum", blocksCountNum)
        if link.pop(self, blocksCountNum) == 0 {
            self.logger.Println("pop val", 0)
            dbIndexes[id].index = link.index
            dbIndexes[id].blockId = link.blockId
            dbIndexes[id].blockSize = blocksCountNum * DB_BASE_SIZE
            if valueLen > dbIndexes[id].blockSize {
                self.logger.Fatal("Invalid  block")
            }
        } else {
            self.logger.Fatal("pop block error")
        }
    }

    self.logger.Println(dbIndexes[id].blockSize)
    self.logger.Println(index)
    self.logger.Println(self.dbsIO[index].fd)
    if dbIndexes[id].blockSize >= valueLen && dbIndexes[id].index >= 0 && self.dbsIO[index].file != nil {
        // self.logger.Println(1)
        index = dbIndexes[id].index
        if self.isMmap && dbIndexes[id].blockId >= 0 && self.dbsIO[index].mmap != nil {
            for k, v := range value {
                self.logger.Println(k, v)
                self.dbsIO[index].mmap[dbIndexes[id].blockId*DB_BASE_SIZE+k] = v
            }

            dbIndexes[id].ndata = valueLen
            ret = id
        } else {
            writeCount, err := self.indexIO.file.WriteAt(value, int64(dbIndexes[id].blockId*DB_BASE_SIZE))
            if err != nil || writeCount != valueLen {
                dbIndexes[id].ndata = 0
                self.logger.Fatal("write index error")
            }

            dbIndexes[id].ndata = valueLen
            ret = id
        }
    }

    if dbIndexes[id].ndata > self.state.dataLenMax {
        self.state.dataLenMax = dbIndexes[id].ndata
    }
    dbIndexes[id].modTime = time.Now().Unix()
    // self.unlockId(id)
    if old.count > 0 {
        link.push(self, old.index, old.blockId, old.count*DB_BASE_SIZE)
    }

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
