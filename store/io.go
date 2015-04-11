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

func (i *IO) close() {
    if i.mmap != nil {
        i.mmap.Unmap()
    }
    if i.file != nil {
        i.file.Close()
    }
}

func (i *IO) initIOMmap() {
    if i.file == nil {
        return
    }
    if i.mmap == nil {
        m, err := MmapFile(int(i.file.Fd()), MaxDbFileSize)
        if err != nil {
            Logger.Fatal(err)
        }
        i.mmap = m
    }
}

func (d *DB) initIOs() {
    for i := 0; i <= d.state.lastId; i++ {
        dbNum := strconv.Itoa(i / MaxDirFileCount)
        currentDbPath := filepath.Join(d.baseDir, DbFileDirName, dbNum)
        if err := os.MkdirAll(currentDbPath, 0755); err != nil {
            Logger.Fatal(err)
        }

        dbFileName := strconv.Itoa(i) + DbFileSuffix
        currentDbFilePath := filepath.Join(currentDbPath, dbFileName)

        d.IOs[i].mutex = &sync.Mutex{}
        file, err := os.OpenFile(currentDbFilePath, os.O_CREATE|os.O_RDWR, 0644)
        if err != nil {
            Logger.Fatal(err)
        }

        fstat, err := file.Stat()
        if err != nil {
            Logger.Fatal(err)
        }

        d.IOs[i].file = file

        if fstat.Size() == 0 {
            d.IOs[i].size = MaxDbFileCount

            if err := file.Truncate(d.IOs[i].size); err != nil {
                Logger.Fatal(err)
            }
        } else {
            d.IOs[i].size = fstat.Size()
        }

        if d.isMmap {
            d.IOs[i].initIOMmap()
        }
    }

    for i := 0; i < MaxMutexCount; i++ {
        d.mutexs[i] = &sync.Mutex{}
    }
}

func (d *DB) checkIndexIOWithId(id int) {
    if id > d.state.maxId {
        d.state.maxId = id
    }

    if id < MaxIndexSize && int64(id)*SizeofIndex >= d.indexIO.end {
        d.indexIO.old = d.indexIO.end
        d.indexIO.end = int64(id)/int64(BaseIndexSize) + 1
        d.indexIO.end += SizeofIndex * int64(BaseIndexSize)
        if err := d.indexIO.file.Truncate(d.indexIO.end); err != nil {
            Logger.Fatal(err)
        }
    }
}
