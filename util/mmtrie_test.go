package util

import (
    // "fmt"
    . "github.com/yqingp/lsearch/util"
    "os"
    "testing"
    "unsafe"
)

func TestMm(t *testing.T) {
    f, _ := os.OpenFile("a.txt", os.O_RDWR|os.O_CREATE, 0666)
    // os.Truncate("a.txt", MMTRIE_BASE_NUM)
    var data MMAP = MmapFile(f)
    a := (*MmtrieState)(unsafe.Pointer(&data[0]))
    t.Log(a)
    f.Close()
}
