package store

import (
    "os"
    "path/filepath"
    "strconv"
    "sync"
)

type IO struct {
    bits  int
    mmap  Mmap
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

func (self *IO) close() {
    if self.mmap != nil {
        self.mmap.Unmap()
    }
    if self.file != nil {
        self.file.Close()
    }
}

func (self *IO) checkIOMmap(db *DB) {
    if self.file == nil {
        return
    }
    if self.mmap == nil {
        m, err := MmapFile(int(self.file.Fd()), MaxDbFileSize)
        if err != nil {
            db.logger.Fatal(err)
        }
        self.mmap = m
    }
}

func (self *DB) initDbsIO() {
    for i := 0; i <= self.state.lastId; i++ {
        currentDbPath := filepath.Join(self.basedir, "base", strconv.Itoa(i/MaxDirFileCount))
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

        if fstat.Size() == 0 {
            self.dbsIO[i].size = MaxDbFileCount

            if err := file.Truncate(self.dbsIO[i].size); err != nil {
                self.logger.Fatal(err)
            }
        } else {
            self.dbsIO[i].size = fstat.Size()
        }

        if self.isMmap {
            self.dbsIO[i].checkIOMmap(self)
        }
    }

    for i := 0; i < MaxMutexCount; i++ {
        self.mutexs[i] = &sync.Mutex{}
    }
}

func (db *DB) checkIndexIOWithId(id int) {
    if id > db.state.dbIdMax {
        db.state.dbIdMax = id
    }

    if id < MaxIndexSize && int64(id)*SizeofIndex >= db.indexIO.end {
        db.indexIO.old = db.indexIO.end
        db.indexIO.end = int64(id)/int64(BaseIndexSize) + 1
        db.indexIO.end += SizeofIndex * int64(BaseIndexSize)
        if err := db.indexIO.file.Truncate(db.indexIO.end); err != nil {
            db.logger.Fatal(err)
        }
    }
}
