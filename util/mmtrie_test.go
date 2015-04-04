package util

import (
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
    v, err := trie.Add([]byte("哈哈11"))
    if err != nil {
        t.Fatal(err)
    }
    t.Log(v)
}
