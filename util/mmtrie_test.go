package util

import (
    "fmt"
    . "github.com/yqingp/lsearch/util"
    // "os"
    "strconv"
    "testing"
    // "unsafe"
)

func TestMmtrieSet(t *testing.T) {
    trie, _ := NewMmtrie("a.txt")
    err := trie.Init()
    if err != nil {
        t.Error(err)
    }
    for i := 0; i < 100; i++ {
        m := strconv.Itoa(i) + "哈哈" + strconv.Itoa(i+1)
        v, err := trie.Set([]byte(m))
        if err != nil {
            t.Fatal(err)
        }
        fmt.Println(m + "===>" + strconv.Itoa(v))
    }
}

func TestMmtrieGet(t *testing.T) {
    trie, _ := NewMmtrie("a.txt")
    err := trie.Init()
    if err != nil {
        t.Error(err)
    }

    fmt.Println(trie.Get([]byte("49980哈哈49981")))
}
