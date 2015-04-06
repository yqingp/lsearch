package util

import (
    "fmt"
    . "github.com/yqingp/lsearch/util"
    "testing"
)

func TestMmtrieGet(t *testing.T) {
    trie, err := Open("a.txt")
    if err != nil {
        t.Error(err)
    }

    fmt.Println(trie.Get([]byte("932哈哈933")))
    trie.Close()
}
