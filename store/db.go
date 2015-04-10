package store

import (
    "errors"
    "log"
    "os"
    "sync"
    "time"
)

type DB struct {
    status          int
    mutex           *sync.Mutex
    blockQueueMutex *sync.Mutex
    indexMutex      *sync.Mutex
    mutexs          [MaxMutexCount]*sync.Mutex
    state           *State
    stateIO         IO
    blockQueueIO    IO
    indexIO         IO
    blockQueues     []BlockQueue
    indexes         []Index
    IOs             [MaxDbFileCount]IO
    baseDir         string
    keyMapTrie      *Mmtrie
    loggerFile      *os.File
    logger          *log.Logger
    isMmap          bool
}

func Open(baseDir string, isMmap bool) (*DB, error) {
    if baseDir == "" {
        return nil, errors.New("basedir is blank")
    }

    db := &DB{}

    db.blockQueueMutex = &sync.Mutex{}
    db.indexMutex = &sync.Mutex{}
    db.mutex = &sync.Mutex{}
    db.baseDir = baseDir
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
    db.initIOs()

    return db, nil
}

func (d *DB) init() {

}

func (self *DB) Close() {
    if self.loggerFile != nil {
        self.loggerFile.Close()
    }
    if self.keyMapTrie != nil {
        self.keyMapTrie.Close()
    }

    self.stateIO.close()
    self.indexIO.close()
    self.blockQueueIO.close()
    for _, v := range self.IOs {
        v.close()
    }
}

// if id is < 1  generate auto increment id
func (self *DB) Set(id int, key []byte, value []byte) (int, error) {
    if key == nil || value == nil || len(value) == 0 {
        return -1, errors.New("key or value is blank")
    }

    var err error

    if id < 1 {
        id, err = self.keyMapTrie.Set(key)
        if err != nil {
            return -1, err
        }
    }

    ret := self.internalSet(id, value)
    if ret == -1 {
        self.logger.Fatal("set error")
    }

    return id, nil
}

func (self *DB) internalSet(id int, value []byte) int {
    ret := -1
    indexes := self.indexes
    if self.status != 0 || indexes == nil {
        return ret
    }

    valueLen := len(value)

    blocksCountNum := 0

    index := 0
    _ = blocksCountNum
    self.indexMutex.Lock()
    self.checkIndexIOWithId(id)
    self.indexMutex.Unlock()

    self.lockId(id)
    defer self.unlockId(id)
    link, old := &BlockQueue{}, &BlockQueue{}

    if indexes[id].blockSize < valueLen {
        if indexes[id].blockSize > 0 {
            old.index = indexes[id].index
            old.blockId = indexes[id].blockId
            old.count = blocksCount(indexes[id].blockSize)
            indexes[id].blockSize = 0
            indexes[id].blockId = 0
            indexes[id].dataLen = 0
        }

        blocksCountNum = blocksCount(valueLen)
        if link.pop(self, blocksCountNum) == 0 {
            indexes[id].index = link.index
            indexes[id].blockId = link.blockId
            indexes[id].blockSize = blocksCountNum * BaseDbSize
            if valueLen > indexes[id].blockSize {
                self.logger.Fatal("Invalid  block")
            }
        } else {
            self.logger.Fatal("pop block error")
        }
    }

    if indexes[id].blockSize >= valueLen && indexes[id].index >= 0 && self.IOs[index].file != nil {
        index = indexes[id].index
        if self.isMmap && indexes[id].blockId >= 0 && self.IOs[index].mmap != nil {
            for k, v := range value {
                self.IOs[index].mmap[indexes[id].blockId*BaseDbSize+k] = v
            }

            indexes[id].dataLen = valueLen
            ret = id
        } else {
            writeCount, err := self.indexIO.file.WriteAt(value, int64(indexes[id].blockId*BaseDbSize))
            if err != nil || writeCount != valueLen {
                indexes[id].dataLen = 0
                self.logger.Fatal("write index error")
            }

            indexes[id].dataLen = valueLen
            ret = id
        }
    }

    if indexes[id].dataLen > self.state.dataLenMax {
        self.state.dataLenMax = indexes[id].dataLen
    }

    indexes[id].updateTime = time.Now().Unix()
    if old.count > 0 {
        link.push(self, old.index, old.blockId, old.count*BaseDbSize)
    }

    return ret
}

func (self *DB) lockId(id int) {
    self.mutexs[id%MaxMutexCount].Lock()
}

func (self *DB) unlockId(id int) {
    self.mutexs[id%MaxMutexCount].Unlock()
}

func (self *DB) Get() {

}

func (self *DB) Del() {

}
