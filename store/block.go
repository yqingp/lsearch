package store

import (
    "bytes"
    "encoding/gob"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "unsafe"
)

type BlockQueue struct {
    index   int
    blockId int
    count   int
}

func (self *DB) initBlockQueue() {
    blockQueueFileName := filepath.Join(self.baseDir, BlockQueueFileName)

    f, err := os.OpenFile(blockQueueFileName, os.O_CREATE|os.O_RDWR, 0664)
    if err != nil {
        self.logger.Fatal(err)
        os.Exit(-1)
    }
    self.blockQueueIO.file = f

    fstat, err := os.Stat(blockQueueFileName)
    if err != nil {
        self.logger.Fatal(err)
    }

    self.blockQueueIO.end = fstat.Size()
    if fstat.Size() == 0 {
        self.blockQueueIO.end = MaxBlockQueueCount * SizeOfBlockQueue
        self.blockQueueIO.size = self.blockQueueIO.end

        if err := os.Truncate(blockQueueFileName, self.blockQueueIO.size); err != nil {
            self.logger.Fatal(err)
        }
    }

    var errNo error
    if self.blockQueueIO.mmap, errNo = MmapFile(int(self.blockQueueIO.file.Fd()), int(self.blockQueueIO.end)); errNo != nil {
        self.logger.Fatal(errNo)
    }

    self.blockQueues = (*[MaxBlockQueueCount]BlockQueue)(unsafe.Pointer(&self.blockQueueIO.mmap[0]))[:MaxBlockQueueCount]
}

func blocksCount(blen int) int {
    ret := blen / BaseDbSize
    if blen%BaseDbSize > 0 {
        ret += 1
    }

    return ret
}

func (self *BlockQueue) pop(db *DB, bcount int) (ret int) {
    ret = -1
    if db == nil || bcount < 1 {
        return
    }

    db.blockQueueMutex.Lock()
    links := db.blockQueues
    var plink *BlockQueue
    _ = plink

    var link BlockQueue

    var buf []byte
    var buf1 bytes.Buffer
    x, index, left, dbId, blockId, blockSize := bcount, -1, 0, -1, -1, 0
    _ = dbId
    _ = blockId
    _ = blockSize
    if links != nil {
        index = links[x].index
    }

    if links != nil && index >= 0 && links[x].count > 0 && index < MaxDbFileCount && db.IOs[index].file != nil {
        self.count = bcount
        self.index = index
        self.blockId = links[x].blockId

        ret = 0

        links[x].count--
        lcount := links[x].count

        if lcount > 0 {

            if db.IOs[index].mmap != nil {
                addr := &db.IOs[index].mmap[links[x].blockId*BaseDbSize]
                plink := (*BlockQueue)(unsafe.Pointer(addr))
                links[x].index = plink.index
                links[x].blockId = plink.blockId
            } else {
                readCount, err := db.indexIO.file.ReadAt(buf[:SizeOfBlockQueue], int64(links[x].blockId*BaseDbSize))
                if err != nil || readCount < 0 {
                    db.logger.Fatal("read index file error")
                }
                buf1.Write(buf)
                dec := gob.NewDecoder(&buf1)
                dec.Decode(&link)
                links[x].index = link.index
                links[x].blockId = link.blockId
            }

        }
    } else {
        x = db.state.lastId
        left = int(db.IOs[x].size) - db.state.lastOff
        if left < BaseDbSize*bcount {
            dbId = x
            blockId = db.state.lastOff / BaseDbSize
            blockSize = left
            db.state.lastOff = BaseDbSize * bcount
            db.state.lastId++
            x = db.state.lastId

            db.logger.Println(x)
            if x >= MaxDbFileCount {
                db.logger.Fatal("pop block dbs error")
            }

            currentDbPath := filepath.Join(db.baseDir, DbFileDirName, strconv.Itoa(x/MaxDirFileCount))
            if err := os.MkdirAll(currentDbPath, 0755); err != nil {
                db.logger.Fatal(err)
            }
            currentDbFileName := filepath.Join(currentDbPath, strconv.Itoa(x)+DbFileSuffix)
            file, err := os.OpenFile(currentDbFileName, os.O_CREATE|os.O_RDWR, 0644)
            if err != nil {
                db.logger.Fatal(err)
            }

            db.IOs[x].file = file

            if err := file.Truncate(MaxDbFileSize); err != nil {
                db.logger.Fatal(err)
            }

            db.IOs[x].mutex = &sync.Mutex{}
            db.IOs[x].end = MaxDbFileSize
            db.IOs[x].size = MaxDbFileSize
            db.IOs[x].initIOMmap(db)
            self.count = bcount
            self.index = x
            self.blockId = 0
            ret = 0
        } else {
            self.count = bcount
            self.index = x
            self.blockId = (db.state.lastOff / BaseDbSize)
            db.state.lastOff += BaseDbSize * bcount
            ret = 0
        }
    }
    db.blockQueueMutex.Unlock()
    if blockId >= 0 {
        self.push(db, dbId, blockId, blockSize)
    }

    return
}

func (self *BlockQueue) push(db *DB, index int, blockId int, blockSize int) int {
    ret, x := -1, 0
    x = blocksCount(blockSize)
    var buf bytes.Buffer
    var link *BlockQueue
    db.blockQueueMutex.Lock()
    defer db.blockQueueMutex.Unlock()
    if db != nil && blockId >= 0 && x > 0 && db.status == 0 && x < MaxBlockQueueCount && index >= 0 && index < MaxDbFileCount {

        if db.blockQueues != nil {
            if db.blockQueues[x].count > 0 {
                if db.IOs[index].mmap != nil {
                    addr := &db.IOs[index].mmap[blockId*BaseDbSize]
                    link = (*BlockQueue)(unsafe.Pointer(addr))
                    link.index = db.blockQueues[x].index
                    link.blockId = db.blockQueues[x].blockId
                } else {
                    enc := gob.NewEncoder(&buf)
                    enc.Encode(&link)
                    link.index = db.blockQueues[x].index
                    link.blockId = db.blockQueues[x].blockId
                    writeCount, err := db.indexIO.file.WriteAt(buf.Bytes()[:SizeOfBlockQueue], int64(blockId*BaseDbSize))
                    if err != nil || writeCount < 0 {
                        db.logger.Fatal("write index file error")
                    }
                }
            }
            db.blockQueues[x].index = index
            db.blockQueues[x].blockId = blockId
            db.blockQueues[x].count++
            ret = 0
        }
    }

    return ret
}
