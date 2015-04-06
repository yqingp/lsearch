package util

import (
    "fmt"
    . "github.com/yqingp/lsearch/util"
    "testing"
)

func TestMmtrieDel(t *testing.T) {
    trie, _ := NewMmtrie("a.txt")
    err := trie.Init()
    if err != nil {
        t.Error(err)
    }
    fmt.Println(trie.Get([]byte("49980哈哈49981")))
    fmt.Println(trie.Del([]byte("49980哈哈49981")))
    fmt.Println(trie.Get([]byte("49980哈哈49981")))
    trie.Close()
}
