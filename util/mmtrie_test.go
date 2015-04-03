package util

import (
    // "fmt"
    . "github.com/yqingp/lsearch/util"
    // "os"
    "testing"
    // "unsafe"
)

func TestMm(t *testing.T) {
    trie, _ := NewMmtrie("a.txt")
    err := trie.Init()
    if err != nil {
        t.Error(err)
    }
    a := trie.Add([]byte("哈哈11"))
    t.Log(a)
}
