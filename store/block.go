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

func (d *DB) initBlockQueue() {
    blockQueueFilePath := filepath.Join(d.baseDir, BlockQueueFileName)

    f, err := os.OpenFile(blockQueueFilePath, os.O_CREATE|os.O_RDWR, 0664)
    if err != nil {
        Logger.Fatal(err)
    }
    d.blockQueueIO.file = f

    fstat, err := os.Stat(blockQueueFilePath)
    if err != nil {
        Logger.Fatal(err)
    }

    d.blockQueueIO.end = fstat.Size()

    if fstat.Size() == 0 {

        d.blockQueueIO.end = MaxBlockQueueCount * SizeOfBlockQueue
        d.blockQueueIO.size = d.blockQueueIO.end

        if err := os.Truncate(blockQueueFilePath, d.blockQueueIO.size); err != nil {
            Logger.Fatal(err)
        }
    }

    var errNo error
    fd := int(d.blockQueueIO.file.Fd())
    mmapSize := int(d.blockQueueIO.end)

    if d.blockQueueIO.mmap, errNo = MmapFile(fd, mmapSize); errNo != nil {
        Logger.Fatal(errNo)
    }

    d.blockQueues = (*[MaxBlockQueueCount]BlockQueue)(unsafe.Pointer(&d.blockQueueIO.mmap[0]))[:MaxBlockQueueCount]
}

func blocksCount(blen int) int {
    val := blen / BaseDbSize

    if (blen % BaseDbSize) > 0 {
        val += 1
    }

    return val
}

func (d *DB) popBlockQueue(blocksCountNum int) *BlockQueue {
    if blocksCountNum < 1 {
        return nil
    }

    d.blockQueueMutex.Lock()

    blockQueue, blockQueues := &BlockQueue{}, d.blockQueues

    var tmpBlockQueue *BlockQueue
    var decodeBlockQueue BlockQueue

    var buf []byte
    var rwBuffer bytes.Buffer

    pos, index, offset, dbId, blockId, blockSize := blocksCountNum, -1, 0, -1, -1, 0

    var offsetSize int64 = 0

    if blockQueues != nil {
        index = blockQueues[pos].index
    }

    if blockQueues != nil && index >= 0 && blockQueues[pos].count > 0 &&
        index < MaxDbFileCount && d.IOs[index].file != nil {
        blockQueue.count = blocksCountNum
        blockQueue.index = index
        blockQueue.blockId = blockQueues[pos].blockId

        blockQueues[pos].count--
        lcount := blockQueues[pos].count

        if lcount > 0 {
            if d.IOs[index].mmap != nil {
                addr := &d.IOs[index].mmap[(blockQueues[pos].blockId * BaseDbSize)]
                tmpBlockQueue = (*BlockQueue)(unsafe.Pointer(addr))
                blockQueues[pos].index = tmpBlockQueue.index
                blockQueues[pos].blockId = tmpBlockQueue.blockId
            } else {
                offsetSize = int64(blockQueues[pos].blockId * BaseDbSize)
                readSize, err := d.indexIO.file.ReadAt(buf[:SizeOfBlockQueue], offsetSize)

                if err != nil || readSize < 0 {
                    Logger.Fatal("read index file error", err)
                }
                rwBuffer.Write(buf)
                dec := gob.NewDecoder(&rwBuffer)
                dec.Decode(&decodeBlockQueue)

                blockQueues[pos].index = decodeBlockQueue.index
                blockQueues[pos].blockId = decodeBlockQueue.blockId
            }
        }
    } else {
        pos = d.state.lastId
        offset = int(d.IOs[pos].size) - d.state.lastOff
        if offset < BaseDbSize*blocksCountNum {
            dbId = pos
            blockId = d.state.lastOff / BaseDbSize
            blockSize = offset
            d.state.lastOff = BaseDbSize * blocksCountNum
            d.state.lastId++
            pos = d.state.lastId

            if pos >= MaxDbFileCount {
                Logger.Fatal("pop block dbs error")
            }

            currentDbPath := filepath.Join(d.baseDir, DbFileDirName, strconv.Itoa(pos/MaxDirFileCount))
            if err := os.MkdirAll(currentDbPath, 0755); err != nil {
                Logger.Fatal(err)
            }

            currentDbFileName := filepath.Join(currentDbPath, strconv.Itoa(pos)+DbFileSuffix)
            file, err := os.OpenFile(currentDbFileName, os.O_CREATE|os.O_RDWR, 0644)

            if err != nil {
                Logger.Fatal(err)
            }

            d.IOs[pos].file = file

            if err := file.Truncate(MaxDbFileSize); err != nil {
                Logger.Fatal(err)
            }

            d.IOs[pos].mutex = &sync.Mutex{}
            d.IOs[pos].end = MaxDbFileSize
            d.IOs[pos].size = MaxDbFileSize
            d.IOs[pos].initIOMmap()
            blockQueue.count = blocksCountNum
            blockQueue.index = pos
            blockQueue.blockId = 0
        } else {
            blockQueue.count = blocksCountNum
            blockQueue.index = pos
            blockQueue.blockId = (d.state.lastOff / BaseDbSize)
            d.state.lastOff += BaseDbSize * blocksCountNum
        }
    }
    d.blockQueueMutex.Unlock()
    if blockId >= 0 {
        d.pushBlockQueue(dbId, blockId, blockSize)
    }

    return blockQueue
}

func (d *DB) pushBlockQueue(index, blockId, blockSize int) {
    pos := 0
    pos = blocksCount(blockSize)
    var buf bytes.Buffer
    var blockQueue *BlockQueue

    d.blockQueueMutex.Lock()
    defer d.blockQueueMutex.Unlock()

    if d != nil && blockId >= 0 && pos > 0 && d.status == 0 &&
        pos < MaxBlockQueueCount && index >= 0 && index < MaxDbFileCount {

        if d.blockQueues != nil {
            if d.blockQueues[pos].count > 0 {
                if d.IOs[index].mmap != nil {
                    addr := &d.IOs[index].mmap[(blockId * BaseDbSize)]
                    blockQueue = (*BlockQueue)(unsafe.Pointer(addr))
                    blockQueue.index = d.blockQueues[pos].index
                    blockQueue.blockId = d.blockQueues[pos].blockId
                } else {
                    enc := gob.NewEncoder(&buf)
                    enc.Encode(&blockQueue)
                    blockQueue.index = d.blockQueues[pos].index
                    blockQueue.blockId = d.blockQueues[pos].blockId
                    writeSize, err := d.indexIO.file.WriteAt(buf.Bytes()[:SizeOfBlockQueue], int64(blockId*BaseDbSize))
                    if err != nil || writeSize < 0 {
                        Logger.Fatal("write index file error")
                    }
                }
            }
            d.blockQueues[pos].index = index
            d.blockQueues[pos].blockId = blockId
            d.blockQueues[pos].count++
        }
    }
}
