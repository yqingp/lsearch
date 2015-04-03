package util

import (
    "fmt"
    "os"
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
    Id      uint64
    Current uint64
    Total   uint64
    Left    uint64
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
