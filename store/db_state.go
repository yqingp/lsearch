package store

import (
    "github.com/yqingp/lsearch/mmap"
    "os"
    "path/filepath"
    "unsafe"
)

type DbState struct {
    status         int
    mode           int
    lastId         int
    lastOff        int
    dbIdMax        int
    dataLenMax     int
    blockIncreMode int
}

func (self *Db) initState() {
    stateFileName := filepath.Join(self.basedir, "db.state")

    f, err := os.OpenFile(stateFileName, os.O_CREATE|os.O_RDWR, 0664)
    if err != nil {
        self.logger.Fatal("create stateFile error")
    }
    self.stateIO.fd = int(f.Fd())
    self.stateIO.file = f

    fstat, err := os.Stat(stateFileName)
    if err != nil {
        self.logger.Fatal("fstat stateFile error")
    }

    self.stateIO.end = fstat.Size()
    if fstat.Size() == 0 {
        self.stateIO.end = SizeOfDbState
        self.stateIO.size = self.stateIO.end

        if err := os.Truncate(stateFileName, self.stateIO.end); err != nil {
            self.logger.Fatal("truncate state file error")
        }
    }

    var errNo error
    if self.stateIO.mmap, errNo = mmap.MmapFile(self.stateIO.fd, int(self.stateIO.end)); errNo != nil {
        self.logger.Fatal("mmap state file error")
    }

    self.state = (*DbState)(unsafe.Pointer(&self.stateIO.mmap[0]))
    self.state.mode = 0
    if self.isMmap {
        self.state.mode = 1
    }
}
