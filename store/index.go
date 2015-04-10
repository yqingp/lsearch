package store

import (
    "os"
    "path/filepath"
    "unsafe"
)

type Index struct {
    blockSize  int
    blockId    int
    dataLen    int
    index      int
    updateTime int64
}

func (self *DB) initIndex() {
    indexFileName := filepath.Join(self.basedir, "db.dbx")

    f, err := os.OpenFile(indexFileName, os.O_CREATE|os.O_RDWR, 0664)
    if err != nil {
        self.logger.Fatal(err)
    }

    self.indexIO.file = f

    fstat, err := os.Stat(indexFileName)
    if err != nil {
        self.logger.Fatal(err)
    }

    self.indexIO.end = fstat.Size()
    self.indexIO.size = MaxIndexSize * SizeofIndex

    var errNo error

    if self.indexIO.mmap, errNo = MmapFile(int(self.indexIO.file.Fd()), int(self.indexIO.size)); errNo != nil {
        self.logger.Fatal(err)
    }

    self.indexes = (*[MaxIndexSize]Index)(unsafe.Pointer(&self.indexIO.mmap[0]))[:MaxIndexSize]

    if fstat.Size() == 0 {
        self.indexIO.end = BaseIndexSize * SizeofIndex

        if err := os.Truncate(indexFileName, self.indexIO.end); err != nil {
            self.logger.Fatal(err)
        }
    }

    self.indexIO.old = self.indexIO.end
}
