package store

import (
    "github.com/yqingp/lsearch/mmap"
    "os"
    "path/filepath"
    "syscall"
    "unsafe"
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

    self.freeBlockQueues = (*[DB_LNK_MAX]DbFreeBlockQueue)(unsafe.Pointer(&self.freeBlockQueueIO.mmap[0]))[:DB_LNK_MAX]
}

func blocksCount(blen int) int {
    ret := blen / DB_BASE_SIZE
    if blen%DB_BASE_SIZE > 0 {
        ret += 1
    }

    return ret
}

func (self *DbFreeBlockQueue) pop(db *Db, bcount int) (ret int) {
    ret = -1
    if db == nil || bcount < 1 {
        return
    }

    db.freeBlockMutex.Lock()
    defer db.freeBlockMutex.Unlock()

    links := db.freeBlockQueues
    var plink *DbFreeBlockQueue
    _ = plink
    x, index := bcount, 0
    if links != nil {
        index = links[x].index
    }

    if links != nil && index >= 0 && links[x].count > 0 && index < DB_MFILE_MAX && db.dbsIO[index].file != nil {
        self.count = bcount
        self.index = index
        self.blockId = links[x].blockId

        ret = 0

        links[x].count--
        lcount := links[x].count

        for lcount > 0 {
            if db.dbsIO[index].mmap != nil {
                addr := &db.dbsIO[index].mmap[links[x].blockId*DB_BASE_SIZE]
                plink := (*DbFreeBlockQueue)(unsafe.Pointer(addr))
                links[x].index = plink.index
                links[x].blockId = plink.blockId
            } else {
                // syscall.Pread(db.indexIO.fd, SizeOfDbFreeBlockQueue, links[x].blockId*DB_BASE_SIZE)
            }
        }
    } else {

    }

    return
}
