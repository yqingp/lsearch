package util

import (
    "fmt"
    . "github.com/yqingp/lsearch/util"
    "testing"
)

func TestMmtrieDel(t *testing.T) {
    trie, err := Open("a.txt")
    if err != nil {
        t.Error(err)
    }
    fmt.Println(trie.Get([]byte("981哈哈982")))
    fmt.Println(trie.Del([]byte("981哈哈982")))
    fmt.Println(trie.Get([]byte("981哈哈982")))
    trie.Close()
}
