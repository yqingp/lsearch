package store

import (
    "github.com/yqingp/lsearch/mmap"
    "os"
    "path/filepath"
    "unsafe"
)

type DbIndex struct {
    blockSize int
    blockId   int
    ndata     int
    index     int
    modTime   int64
}

func (self *Db) initIndex() {
    indexFileName := filepath.Join(self.basedir, "db.dbx")

    f, err := os.OpenFile(indexFileName, os.O_CREATE|os.O_RDWR, 0664)
    if err != nil {
        self.logger.Fatal(err)
    }
    self.indexIO.fd = int(f.Fd())
    self.indexIO.file = f

    fstat, err := os.Stat(indexFileName)
    if err != nil {
        self.logger.Fatal(err)
    }

    self.indexIO.end = fstat.Size()
    self.indexIO.size = DB_DBX_MAX * SizeofDbIndex

    var errNo error

    if self.indexIO.mmap, errNo = mmap.MmapFile(self.indexIO.fd, int(self.indexIO.size)); errNo != nil {
        self.logger.Fatal(err)
    }

    self.indexes = (*[DB_DBX_MAX]DbIndex)(unsafe.Pointer(&self.indexIO.mmap[0]))[:DB_DBX_MAX]

    if fstat.Size() == 0 {
        self.indexIO.end = DB_DBX_BASE * SizeofDbIndex

        if err := os.Truncate(indexFileName, self.indexIO.end); err != nil {
            self.logger.Fatal(err)
        }
    }

    self.indexIO.old = self.indexIO.end
}
