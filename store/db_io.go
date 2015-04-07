package store

import (
    "github.com/yqingp/lsearch/mmap"
    "os"
    "path/filepath"
    "strconv"
    "sync"
)

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

func (self *Db) initDbsIO() {
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
        self.dbsIO[i].fd = int(file.Fd())

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

func (db *Db) checkIndexIOWithId(id int) {
    if id > db.state.dbIdMax {
        db.state.dbIdMax = id
    }

    if id < DB_DBX_MAX && int64(id)*SizeofDbIndex >= db.indexIO.end {
        db.indexIO.old = db.indexIO.end
        db.indexIO.end = int64(id)/int64(DB_DBX_BASE) + 1
        db.indexIO.end += SizeofDbIndex * int64(DB_DBX_BASE)
        if err := db.indexIO.file.Truncate(db.indexIO.end); err != nil {
            db.logger.Fatal("db index file truncate error; exit")
        }
    }
}
