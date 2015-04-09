package store

import (
    "bytes"
    "encoding/gob"
    "github.com/yqingp/lsearch/mmap"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "unsafe"
)

type DbBlockQueue struct {
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
    blockQueueFileName := filepath.Join(self.basedir, "db.blkq")

    f, err := os.OpenFile(blockQueueFileName, os.O_CREATE|os.O_RDWR, 0664)
    if err != nil {
        self.logger.Fatal(err)
        os.Exit(-1)
    }
    self.blockQueueIO.fd = int(f.Fd())
    self.blockQueueIO.file = f

    fstat, err := os.Stat(blockQueueFileName)
    if err != nil {
        self.logger.Fatal(err)
    }

    self.blockQueueIO.end = fstat.Size()
    if fstat.Size() == 0 {
        self.blockQueueIO.end = DB_LNK_MAX * SizeOfDbBlockQueue
        self.blockQueueIO.size = self.blockQueueIO.end

        if err := os.Truncate(blockQueueFileName, self.blockQueueIO.size); err != nil {
            self.logger.Fatal(err)
        }
    }

    var errNo error
    if self.blockQueueIO.mmap, errNo = mmap.MmapFile(self.blockQueueIO.fd, int(self.blockQueueIO.end)); errNo != nil {
        self.logger.Fatal(errNo)
    }

    self.blockQueues = (*[DB_LNK_MAX]DbBlockQueue)(unsafe.Pointer(&self.blockQueueIO.mmap[0]))[:DB_LNK_MAX]
}

func blocksCount(blen int) int {
    ret := blen / DB_BASE_SIZE
    if blen%DB_BASE_SIZE > 0 {
        ret += 1
    }

    return ret
}

func (self *DbBlockQueue) pop(db *Db, bcount int) (ret int) {
    ret = -1
    if db == nil || bcount < 1 {
        return
    }

    db.freeBlockMutex.Lock()
    defer db.freeBlockMutex.Unlock()

    links := db.blockQueues
    var plink *DbBlockQueue
    _ = plink

    var link DbBlockQueue

    var buf []byte
    var buf1 bytes.Buffer
    x, index, left, db_id, block_id, block_size := bcount, 0, 0, -1, -1, 0
    _ = db_id
    _ = block_id
    _ = block_size
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

        if lcount > 0 {
            if db.dbsIO[index].mmap != nil {
                addr := &db.dbsIO[index].mmap[links[x].blockId*DB_BASE_SIZE]
                plink := (*DbBlockQueue)(unsafe.Pointer(addr))
                links[x].index = plink.index
                links[x].blockId = plink.blockId
            } else {

                readCount, err := db.indexIO.file.ReadAt(buf[:SizeOfDbBlockQueue], int64(links[x].blockId*DB_BASE_SIZE))
                if err != nil {
                    return
                }
                if readCount > 0 {
                    buf1.Write(buf)
                    dec := gob.NewDecoder(&buf1)
                    dec.Decode(&link)
                    links[x].index = link.index
                    links[x].blockId = link.blockId
                }
            }
        }
    } else {
        x = db.state.lastId
        left = int(db.dbsIO[x].size) - db.state.lastOff
        if left < DB_BASE_SIZE*bcount {
            db_id = x
            block_id = db.state.lastOff / DB_BASE_SIZE
            block_size = left
            db.state.lastOff = DB_BASE_SIZE * bcount
            db.state.lastId++
            x = db.state.lastId

            if x >= DB_MFILE_MAX {
                db.logger.Fatal("pop block dbs error")
            }

            currentDbPath := filepath.Join(db.basedir, "base", strconv.Itoa(x/DB_DIR_FILES))
            if err := os.MkdirAll(currentDbPath, 0755); err != nil {
                db.logger.Fatal(err)
            }
            currentDbFileName := filepath.Join(currentDbPath, strconv.Itoa(x)+".db")
            file, err := os.OpenFile(currentDbFileName, os.O_CREATE|os.O_RDWR, 0644)
            if err != nil {
                db.logger.Fatal(err)
            }

            db.dbsIO[x].fd = int(file.Fd())
            db.dbsIO[x].file = file

            if err := file.Truncate(DB_MFILE_SIZE); err != nil {
                db.logger.Fatal(err)
            }

            db.dbsIO[x].mutex = &sync.Mutex{}
            db.dbsIO[x].end = DB_MFILE_SIZE
            db.dbsIO[x].size = DB_MFILE_SIZE
            db.dbsIO[x].checkDbIOMmap(db)
            self.count = bcount
            self.index = x
            self.blockId = 0
            ret = 0
        } else {
            self.count = bcount
            self.index = x
            self.blockId = (db.state.lastOff / DB_BASE_SIZE)
            db.state.lastOff += DB_BASE_SIZE * bcount
            ret = 0
        }
    }

    if block_id >= 0 {
        self.push(db, db_id, block_id, block_size)
    }

    return
}

func (self *DbBlockQueue) push(db *Db, index int, blockid int, block_size int) {

}
