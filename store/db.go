package store

import (
    "errors"
    "github.com/yqingp/lsearch/mmap"
    "github.com/yqingp/lsearch/util"
    "log"
    "os"
    "path/filepath"
    "strconv"
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
    file  *os.File
}

type DbFreeBlockQueue struct {
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

func (self *Db) initKmap() error {
    var err error
    kmapfileName := filepath.Join(self.basedir, "db.kmap")
    if self.kmap, err = util.Open(kmapfileName); err != nil {
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

    self.logger = log.New(f, "[LSearch][DB]:", log.Lshortfile|log.Ldate|log.Ltime)
    return nil
}

func (self *Db) initState() {
    stateFileName := filepath.Join(self.basedir, "db.state")

    f, err := os.OpenFile(stateFileName, os.O_CREATE|os.O_RDWR, 0664)
    if err != nil {
        self.logger.Fatal("create stateFile error")
    }
    self.stateIO.fd = int(f.Fd())
    self.stateIO.file = f

    fstat, err := os.Stat(stateFileName)
    if err != nil {
        self.logger.Fatal("fstat stateFile error")
    }

    self.stateIO.end = fstat.Size()
    if fstat.Size() == 0 {
        self.stateIO.end = SizeOfDbState
        self.stateIO.size = self.stateIO.end

        if err := os.Truncate(stateFileName, self.stateIO.end); err != nil {
            self.logger.Fatal("truncate state file error")
        }
    }

    var errNo error
    if self.stateIO.mmap, errNo = mmap.MmapFile(self.stateIO.fd, int(self.stateIO.end)); errNo != nil {
        self.logger.Fatal("mmap state file error")
    }

    self.state = (*DbState)(unsafe.Pointer(&self.stateIO.mmap[0]))
    self.state.mode = 0
    if self.isMmap {
        self.state.mode = 1
    }
}

func (self *Db) initFreeBlockQueue() {
    freeBlockQueueFileName := filepath.Join(self.basedir, "db.freeq")

    f, err := os.OpenFile(freeBlockQueueFileName, os.O_CREATE|os.O_RDWR, 0664)
    if err != nil {
        self.logger.Fatal("create free block queue error")
        os.Exit(-1)
    }
    self.freeBlockQueueIO.fd = int(f.Fd())
    self.freeBlockQueueIO.file = f

    fstat, err := os.Stat(freeBlockQueueFileName)
    if err != nil {
        self.logger.Fatal("stat free block queue error")
    }

    self.freeBlockQueueIO.end = fstat.Size()
    if fstat.Size() == 0 {
        self.freeBlockQueueIO.end = DB_LNK_MAX * SizeOfDbFreeBlockQueue
        self.freeBlockQueueIO.size = self.freeBlockQueueIO.end

        if err := os.Truncate(freeBlockQueueFileName, self.freeBlockQueueIO.size); err != nil {
            self.logger.Fatal("truncate stat free block queue file error")
        }
    }

    var errNo error
    if self.freeBlockQueueIO.mmap, errNo = mmap.MmapFile(self.freeBlockQueueIO.fd, int(self.freeBlockQueueIO.end)); errNo != nil {
        self.logger.Fatal("mmap stat free block queue file error")
    }
}

func (self *Db) initIndex() {
    indexFileName := filepath.Join(self.basedir, "db.dbx")

    f, err := os.OpenFile(indexFileName, os.O_CREATE|os.O_RDWR, 0664)
    if err != nil {
        self.logger.Fatal("create free block queue error")
    }
    self.indexIO.fd = int(f.Fd())
    self.indexIO.file = f

    fstat, err := os.Stat(indexFileName)
    if err != nil {
        self.logger.Fatal("stat free block queue error")
    }

    self.indexIO.end = fstat.Size()
    self.indexIO.size = DB_DBX_MAX * SizeofDbIndex

    var errNo error

    if self.indexIO.mmap, errNo = mmap.MmapFile(self.indexIO.fd, int(self.indexIO.size)); errNo != nil {
        self.logger.Fatal("mmap stat free block queue file error")
    }

    if fstat.Size() == 0 {
        self.indexIO.end = DB_DBX_BASE * SizeofDbIndex

        if err := os.Truncate(indexFileName, self.indexIO.end); err != nil {
            self.logger.Fatal("truncate stat free block queue file error")
        }
    }

    self.indexIO.old = self.indexIO.end
}

func (self *Db) initDbs() {
    for i := 0; i < self.state.lastId; i++ {
        currentDbPath := filepath.Join(self.basedir, "base", strconv.Itoa(i/DB_DIR_FILES))
        if err := os.MkdirAll(currentDbPath, 0755); err != nil {
            self.logger.Fatal("init dbs: mkdir error; check perm")
        }

        currentDbFileName := filepath.Join(currentDbPath, strconv.Itoa(i)+".db")
        self.dbsIO[i].mutex = &sync.Mutex{}
        file, err := os.OpenFile(currentDbFileName, os.O_CREATE|os.O_RDWR, 0644)
        if err != nil {
            self.logger.Fatal("init dbs:open db file error")
        }

        fstat, err := file.Stat()
        if err != nil {
            self.logger.Fatal("init dbs: db file stat error")
        }

        self.dbsIO[i].file = file
        self.stateIO.fd = int(file.Fd())

        if fstat.Size() == 0 {
            self.dbsIO[i].size = DB_MFILE_MAX

            if err := file.Truncate(self.dbsIO[i].size); err != nil {
                self.logger.Fatal("init dbs: truncate error")
            }
        } else {
            self.dbsIO[i].size = fstat.Size()
        }

        if self.isMmap {
            self.checkDbIOMmap(i)
        }
    }

    for i := 0; i < DB_MUTEX_MAX; i++ {
        self.mutexs[i] = &sync.Mutex{}
    }
}

func (self *Db) checkDbIOMmap(i int) {
    dbsio := self.dbsIO[i]

    if dbsio.fd < 1 || dbsio.file == nil {
        return
    }
    if dbsio.mmap == nil {
        m, err := mmap.MmapFile(dbsio.fd, int(dbsio.size))
        if err != nil {
            self.logger.Fatal("init dbs mmap error")
        }
        dbsio.mmap = m
    }

}
