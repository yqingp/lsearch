package store

import (
    "os"
    "path/filepath"
    "unsafe"
)

type State struct {
    status         int
    mode           int
    lastId         int
    lastOff        int
    dbIdMax        int
    dataLenMax     int
    blockIncreMode int
}

func (self *DB) initState() {
    stateFileName := filepath.Join(self.baseDir, StateFileName)

    f, err := os.OpenFile(stateFileName, os.O_CREATE|os.O_RDWR, 0664)
    if err != nil {
        self.logger.Fatal(err)
    }
    self.stateIO.file = f

    fstat, err := os.Stat(stateFileName)
    if err != nil {
        self.logger.Fatal(err)
    }

    self.stateIO.end = fstat.Size()
    if fstat.Size() == 0 {
        self.stateIO.end = SizeOfState
        self.stateIO.size = self.stateIO.end

        if err := os.Truncate(stateFileName, self.stateIO.end); err != nil {
            self.logger.Fatal(err)
        }
    }

    var errNo error
    if self.stateIO.mmap, errNo = MmapFile(int(self.stateIO.file.Fd()), int(self.stateIO.end)); errNo != nil {
        self.logger.Fatal(err)
    }

    self.state = (*State)(unsafe.Pointer(&self.stateIO.mmap[0]))
    self.state.mode = 0
    if self.isMmap {
        self.state.mode = 1
    }
}
