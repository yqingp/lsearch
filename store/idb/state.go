package idb

import (
    . "github.com/yqingp/lsearch/store/mmap"
    "os"
    "path/filepath"
    "unsafe"
)

type State struct {
    status         int
    mode           int
    lastId         int
    lastOff        int
    maxId          int
    dataLenMax     int
    blockIncreMode int
}

func (d *DB) initState() {
    stateFilePath := filepath.Join(d.baseDir, StateFileName)

    f, err := os.OpenFile(stateFilePath, os.O_CREATE|os.O_RDWR, 0664)
    if err != nil {
        Logger.Fatal(err)
    }
    d.stateIO.file = f

    fstat, err := os.Stat(stateFilePath)
    if err != nil {
        Logger.Fatal(err)
    }

    d.stateIO.end = fstat.Size()
    if fstat.Size() == 0 {
        d.stateIO.end = SizeOfState
        d.stateIO.size = d.stateIO.end

        if err := os.Truncate(stateFilePath, d.stateIO.end); err != nil {
            Logger.Fatal(err)
        }
    }

    var errNo error

    fd := int(d.stateIO.file.Fd())

    if d.stateIO.mmap, errNo = MmapFile(fd, int(d.stateIO.end)); errNo != nil {
        Logger.Fatal(err)
    }

    d.state = (*State)(unsafe.Pointer(&d.stateIO.mmap[0]))
    d.state.mode = 0
    if d.isMmap {
        d.state.mode = 1
    }
}
