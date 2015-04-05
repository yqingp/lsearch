package util

import (
    "fmt"
    . "github.com/yqingp/lsearch/util"
    "os"
    "runtime/pprof"
    // "strconv"
    "testing"
    // "unsafe"
)

func TestMm(t *testing.T) {
    f, err := os.Create("pp.prof")
    if err != nil {
        panic(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    trie, _ := NewMmtrie("a.txt")
    err = trie.Init()
    if err != nil {
        t.Error(err)
    }
    // for i := 0; i < 50000; i++ {
    //     m := strconv.Itoa(i) + "哈哈" + strconv.Itoa(i+1)
    //     v, err := trie.Set([]byte(m))
    //     if err != nil {
    //         t.Fatal(err)
    //     }
    //     fmt.Println(m + "===>" + strconv.Itoa(v))
    // }

    fmt.Println(trie.Get([]byte("49980哈哈49981")))
}
