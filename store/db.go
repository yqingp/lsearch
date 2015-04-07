package store

import (
    "errors"
    "github.com/yqingp/lsearch/util"
    "log"
    "os"
    "sync"
    "unsafe"
)

const (
    DB_LNK_MAX            = 2097152
    DB_LNK_INCREMENT      = 65536
    DB_DBX_MAX            = 2000000000
    DB_DBX_BASE           = 1000000
    DB_BASE_SIZE          = 64
    DB_PATH_MAX           = 1024
    DB_DIR_FILES          = 64
    DB_BUF_SIZE           = 4096
    DB_XBLOCKS_MAX        = 14
    DB_MBLOCKS_MAX        = 1024
    DB_MBLOCK_BASE        = 4096
    DB_MBLOCK_MAX         = 33554432
    DB_MUTEX_MAX          = 65536
    DB_USE_MMAP           = 0x01
    DB_MFILE_SIZE         = 268435456
    DB_MFILE_MAX          = 8129
    DB_BLOCK_INCRE_LEN    = 0x0
    DB_BLOCK_INCRE_DOUBLE = 0x1
)

type Db struct {
    status           int
    blockMax         int
    mmTotal          int64
    xxTotal          int64
    mutex            *sync.Mutex
    freeBlockMutex   *sync.Mutex
    indexMutex       *sync.Mutex
    blockMutex       *sync.Mutex
    state            *DbState
    stateIO          DbIO
    freeBlockQueueIO DbIO
    indexIO          DbIO
    dbsIO            [DB_MFILE_MAX]DbIO
    blocks           [DB_XBLOCKS_MAX]DbBlock
    basedir          string
    kmap             *util.Mmtrie
    loggerFile       *os.File
    logger           *log.Logger
    isMmap           bool
    mutexs           [DB_MUTEX_MAX]*sync.Mutex
}

var (
    SizeOfDbState          = int64(unsafe.Sizeof(DbState{}))
    SizeOfDbFreeBlockQueue = int64(unsafe.Sizeof(DbFreeBlockQueue{}))
    SizeofDbIndex          = int64(unsafe.Sizeof(DbIndex{}))
)

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
    if key == nil || value == nil {
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

func (self *Db) set(id int, value []byte) {

}

func (self *Db) Get() {

}

func (self *Db) Del() {

}
