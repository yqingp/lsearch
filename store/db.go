package store

import (
    "errors"
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

func (d *DB) Close() {
    if d.loggerFile != nil {
        d.loggerFile.Close()
    }
    if d.keyMapTrie != nil {
        d.keyMapTrie.Close()
    }

    d.stateIO.close()
    d.indexIO.close()
    d.blockQueueIO.close()
    for _, v := range d.IOs {
        v.close()
    }
}

// if id is < 1  generate auto increment id
func (d *DB) Set(id int, key []byte, value []byte) (int, error) {
    if key == nil || value == nil || len(value) == 0 {
        return -1, errors.New("key or value is blank")
    }

    var err error

    if id < 1 {
        id, err = d.keyMapTrie.Set(key)
        if err != nil {
            return -1, err
        }
    }

    ret := d.internalSet(id, value)
    if ret == -1 {
        Logger.Fatal("set error")
    }

    return id, nil
}

func (d *DB) internalSet(id int, value []byte) int {
    ret := -1
    indexes := d.indexes
    if d.status != 0 || indexes == nil {
        return ret
    }

    valueLen := len(value)

    blocksCountNum, index := 0, 0

    d.indexMutex.Lock()
    d.checkIndexIOWithId(id)
    d.indexMutex.Unlock()

    d.lockId(id)
    defer d.unlockId(id)

    oldBlockQueue := &BlockQueue{}
    var newBlockQueue *BlockQueue

    if indexes[id].blockSize < valueLen {
        if indexes[id].blockSize > 0 {
            oldBlockQueue.index = indexes[id].index
            oldBlockQueue.blockId = indexes[id].blockId
            oldBlockQueue.count = blocksCount(indexes[id].blockSize)
            indexes[id].blockSize = 0
            indexes[id].blockId = 0
            indexes[id].dataLen = 0
        }

        blocksCountNum = blocksCount(valueLen)
        newBlockQueue = d.popBlockQueue(blocksCountNum)
        if newBlockQueue != nil {
            indexes[id].index = newBlockQueue.index
            indexes[id].blockId = newBlockQueue.blockId
            indexes[id].blockSize = blocksCountNum * BaseDbSize
            if valueLen > indexes[id].blockSize {
                Logger.Fatal("Invalid  block")
            }
        } else {
            Logger.Fatal("pop block error")
        }
    }

    if indexes[id].blockSize >= valueLen && indexes[id].index >= 0 &&
        d.IOs[index].file != nil {

        index = indexes[id].index
        if d.isMmap && indexes[id].blockId >= 0 && d.IOs[index].mmap != nil {
            for k, v := range value {
                d.IOs[index].mmap[indexes[id].blockId*BaseDbSize+k] = v
            }

            indexes[id].dataLen = valueLen
            ret = id
        } else {
            writeSize, err := d.indexIO.file.WriteAt(value, int64(indexes[id].blockId*BaseDbSize))
            if err != nil || writeSize != valueLen {
                indexes[id].dataLen = 0
                Logger.Fatal("write index error")
            }

            indexes[id].dataLen = valueLen
            ret = id
        }
    }

    if indexes[id].dataLen > d.state.dataLenMax {
        d.state.dataLenMax = indexes[id].dataLen
    }

    indexes[id].updateTime = time.Now().Unix()
    if oldBlockQueue.count > 0 {
        d.pushBlockQueue(oldBlockQueue.index, oldBlockQueue.blockId, oldBlockQueue.count*BaseDbSize)
    }

    return ret
}

func (d *DB) lockId(id int) {
    d.mutexs[id%MaxMutexCount].Lock()
}

func (d *DB) unlockId(id int) {
    d.mutexs[id%MaxMutexCount].Unlock()
}

func (d *DB) Get() {

}

func (d *DB) Del() {

}
