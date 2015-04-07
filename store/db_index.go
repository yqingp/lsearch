package store

import (
    "github.com/yqingp/lsearch/mmap"
    "os"
    "path/filepath"
)

type DbIndex struct {
    blockSize int
    blockId   int
    ndata     int
    index     int
    modTime   int
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
