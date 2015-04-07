package store

import (
    "github.com/yqingp/lsearch/mmap"
    "os"
    "path/filepath"
)

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
