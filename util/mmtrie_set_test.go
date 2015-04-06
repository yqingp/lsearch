package util

import (
    "fmt"
    . "github.com/yqingp/lsearch/util"
    "strconv"
    "testing"
)

func TestMmtrieSet(t *testing.T) {
    trie, _ := NewMmtrie("a.txt")
    err := trie.Init()
    if err != nil {
        t.Error(err)
    }
    for i := 1000000; i < 1000005; i++ {
        m := strconv.Itoa(i) + "哈哈" + strconv.Itoa(i+1)
        v, err := trie.Set([]byte(m))
        if err != nil {
            t.Fatal(err)
        }
        fmt.Println(m + "===>" + strconv.Itoa(v))
    }
    trie.Close()
}
