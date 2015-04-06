package store

import (
    "errors"
    "github.com/yqingp/lsearch/mmap"
    "github.com/yqingp/lsearch/util"
    "log"
    "os"
    "path/filepath"
    "sync"
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

type DbIndex struct {
    blockSize int
    blockId   int
    ndata     int
    index     int
    modTime   int
}

type DbIO struct {
    fd    int
    bits  int
    mmap  mmap.Mmap
    mutex *sync.Mutex
    old   int64
    end   int64
    size  int64
}

type DbFreeBlock struct {
    index   int
    blockId int
    count   int
}

type DbBlockMap struct {
    blockSize int
    blocksMax int
}

type DbBlock struct {
    mblocks  [DB_MBLOCKS_MAX]string
    nmblocks int
    total    int
}

type DbState struct {
    status         int
    mode           int
    lastId         int
    lastOff        int
    dbIdMax        int
    dataLenMax     int
    blockIncreMode int
}

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
    freeBlockIO    DbIO
    indexIO        DbIO
    dbsIO          [DB_MFILE_MAX]DbIO
    blocks         [DB_XBLOCKS_MAX]DbBlock
    basedir        string
    kmap           *util.Mmtrie
    loggerFile     *os.File
    logger         *log.Logger
}

func NewDb(basedir string, isMmap bool) (*Db, error) {
    if basedir == "" {
        return nil, errors.New("basedir is blank")
    }

    db := &Db{}

    db.freeBlockMutex = &sync.Mutex{}
    db.indexMutex = &sync.Mutex{}
    db.blockMutex = &sync.Mutex{}
    db.mutex = &sync.Mutex{}
    db.basedir = basedir

    if err := db.initKmap(); err != nil {
        return nil, err
    }

    if err := db.initLogger(); err != nil {
        return nil, err
    }

    return db, nil
}

func (self *Db) initKmap() error {
    var err error
    kmapfileName := filepath.Join(self.basedir, "db.kmap")
    if self.kmap, err = util.NewMmtrie(kmapfileName); err != nil {
        return err
    }
    if err = self.kmap.Init(); err != nil {
        return err
    }

    return nil
}

func (self *Db) initLogger() error {
    loggerFileName := filepath.Join(self.basedir, "db.log")

    f, err := os.OpenFile(loggerFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)

    if err != nil {
        return err
    }

    self.logger = log.New(f, "LSearch:DB:", log.Lshortfile|log.Ldate|log.Ltime)
    return nil
}
