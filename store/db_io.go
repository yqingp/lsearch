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

// func (self *Db) checkDbIOMmap(i int) {
//     dbsio := self.dbsIO[i]

//     if dbsio.fd < 1 || dbsio.file == nil {
//         return
//     }
//     if dbsio.mmap == nil {
//         m, err := mmap.MmapFile(dbsio.fd, int(dbsio.size))
//         if err != nil {
//             self.logger.Fatal(err)
//         }
//         dbsio.mmap = m
//     }

// }

func (self *DbIO) close() {
    if self.mmap != nil {
        self.mmap.Unmap()
    }
    if self.file != nil {
        self.file.Close()
    }
}

func (self *DbIO) checkDbIOMmap(db *Db) {
    if self.fd < 1 || self.file == nil {
        return
    }
    if self.mmap == nil {
        m, err := mmap.MmapFile(self.fd, DB_MFILE_SIZE)
        if err != nil {
            db.logger.Fatal(err)
        }
        self.mmap = m
    }
}

func (self *Db) initDbsIO() {
    for i := 0; i <= self.state.lastId; i++ {
        currentDbPath := filepath.Join(self.basedir, "base", strconv.Itoa(i/DB_DIR_FILES))
        if err := os.MkdirAll(currentDbPath, 0755); err != nil {
            self.logger.Fatal(err)
        }

        currentDbFileName := filepath.Join(currentDbPath, strconv.Itoa(i)+".db")
        self.dbsIO[i].mutex = &sync.Mutex{}
        file, err := os.OpenFile(currentDbFileName, os.O_CREATE|os.O_RDWR, 0644)
        if err != nil {
            self.logger.Fatal(err)
        }

        fstat, err := file.Stat()
        if err != nil {
            self.logger.Fatal(err)
        }

        self.dbsIO[i].file = file
        self.dbsIO[i].fd = int(file.Fd())

        if fstat.Size() == 0 {
            self.dbsIO[i].size = DB_MFILE_MAX

            if err := file.Truncate(self.dbsIO[i].size); err != nil {
                self.logger.Fatal(err)
            }
        } else {
            self.dbsIO[i].size = fstat.Size()
        }

        if self.isMmap {
            self.dbsIO[i].checkDbIOMmap(self)
            // self.checkDbIOMmap(i)
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
            db.logger.Fatal(err)
        }
    }
}
