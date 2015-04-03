package util

import (
    "fmt"
    "os"
    "sync"
    "syscall"
    "unsafe"
)

const (
    MMTRIE_PATH_MAX  = 256
    MMTRIE_LINE_MAX  = 256
    MMTRIE_BASE_NUM  = 1000000
    MMTRIE_NODES_MAX = 1000000000
    MMTRIE_WORD_MAX  = 4096
)

type MmtrieState struct {
    id      uint64
    current uint64
    total   uint64
    left    uint64
}

type MmtrieList struct {
    count uint64
    head  uint64
}

type MmtrieNode struct {
    key     byte
    nchilds uint8
    data    uint64
    childs  uint64
    list    [MMTRIE_LINE_MAX]MmtrieList
}

type Mmtrie struct {
    state    *MmtrieState
    nodes    []*MmtrieNode
    mmap     MMAP
    size     uint64
    old      uint64
    fileSize uint64
    fd       int
    bits     int
    mutex    *sync.Mutex
}
type MMAP []byte

func MmapFile(f *os.File) []byte {
    data, err := syscall.Mmap(int(f.Fd()), 0, MMTRIE_NODES_MAX, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
    if err != nil {
        fmt.Println(err)
    }

    return data
}

func Unmmap(f *os.File) {

}

func (m MMAP) GetState() *MmtrieState {
    return (*MmtrieState)(unsafe.Pointer(&m[0]))
}

func (m *MmtrieState) ToS() {
    fmt.Println(m)
}
