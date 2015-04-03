package util

import (
    "errors"
    "fmt"
    . "github.com/yqingp/lsearch/mmap"
    "os"
    "sync"
    // "syscall"
    "unsafe"
)

const (
    MMTRIE_PATH_MAX  = 256
    MMTRIE_LINE_MAX  = 256
    MMTRIE_BASE_NUM  = 500000
    MMTRIE_NODES_MAX = 1000000000
    MMTRIE_WORD_MAX  = 4096
)

type MmtrieState struct {
    id      uint64
    current uint64
    total   uint64
    left    uint64
    list    [MMTRIE_LINE_MAX]MmtrieList
}

type MmtrieList struct {
    count uint64
    head  uint64
}

type MmtrieNode struct {
    key     uint8
    nchilds uint8
    data    int
    childs  int
}

type Mmtrie struct {
    state       *MmtrieState
    nodes       []*MmtrieNode
    nodesOffset uint64
    mmap        Mmap
    size        int
    old         uint64
    fileSize    int64
    fd          int
    bits        int
    mutex       *sync.Mutex
    isInit      bool
    filename    string
}

var (
    t                 MmtrieState
    t1                MmtrieNode
    MmtrieNodeSizeOf  = int(unsafe.Sizeof(t1))
    MmtrieStateSizeOf = int(unsafe.Sizeof(t))
)

func NewMmtrie(filename string) (*Mmtrie, error) {
    if filename == "" {
        return nil, errors.New("file name is blank")
    }
    return &Mmtrie{filename: filename}, nil
}

func (m *Mmtrie) Init() error {
    if m.isInit {
        return errors.New("is inited")
    }

    f, err := os.OpenFile(m.filename, os.O_CREATE|os.O_RDWR, 0644)
    if err != nil {
        return err
    }
    m.fd = int(f.Fd())

    fstat, err := os.Stat(m.filename)
    if err != nil {
        return err
    }

    if m.mmap == nil {
        m.size = MmtrieStateSizeOf + MMTRIE_NODES_MAX*MmtrieNodeSizeOf
        mp, err := MmapFile(m.fd, m.size)
        if err != nil {
            return err
        }
        m.mmap = mp
        m.state = (*MmtrieState)(unsafe.Pointer(&m.mmap[0]))
        // addr := &m.mmap[MmtrieStateSizeOf]
        m.nodesOffset = uint64(MmtrieStateSizeOf)
    }

    if fstat.Size() == 0 {
        m.fileSize = int64(MmtrieStateSizeOf) + MMTRIE_BASE_NUM*int64(MmtrieNodeSizeOf)
        if err := f.Truncate(m.fileSize); err != nil {
            return err
        }
        m.state.total = MMTRIE_BASE_NUM
        m.state.left = MMTRIE_BASE_NUM - MMTRIE_LINE_MAX
        m.state.current = MMTRIE_LINE_MAX
    }

    return nil
}

func (m *Mmtrie) ToS() {
    fmt.Println(m.state)
}
