package store

import (
    "github.com/yqingp/lsearch/mmap"
    "github.com/yqingp/lsearch/util"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "unsafe"
)

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
