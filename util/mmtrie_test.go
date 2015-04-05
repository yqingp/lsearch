package util

import (
    "fmt"
    . "github.com/yqingp/lsearch/util"
    "os"
    "runtime/pprof"
    "strconv"
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
    // v, _ := fmt.Println([]byte("5哈哈6"))
    // v, _ = fmt.Println([]byte("52哈哈53"))
    v, _ := trie.Add([]byte("1哈哈2"))
    fmt.Println("1哈哈2" + "===>" + strconv.Itoa(v))
    v, _ = trie.Add([]byte("10哈哈11"))
    fmt.Println("10哈哈11" + "===>" + strconv.Itoa(v))
    // for i := 0; i < 11; i++ {
    //     m := strconv.Itoa(i) + "哈哈" + strconv.Itoa(i+1)
    //     v, err := trie.Add([]byte(m))
    //     if err != nil {
    //         t.Fatal(err)
    //     }
    //     fmt.Println(m + "===>" + strconv.Itoa(v))
    // }

}
