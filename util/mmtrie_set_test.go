package util

import (
    // "fmt"
    . "github.com/yqingp/lsearch/util"
    "strconv"
    "testing"
)

func TestMmtrieSet(t *testing.T) {
    trie, err := Open("a.txt")
    // t.Log(trie)
    if err != nil {
        t.Error(err)
    }

    for i := 1; i < 1000000; i++ {
        m := strconv.Itoa(i) + "哈哈" + strconv.Itoa(i+1)
        // fmt.Println(m)
        _, err := trie.Set([]byte(m))
        if err != nil {
            t.Fatal(err)
        }
        // fmt.Println(m + "===>" + strconv.Itoa(v))
    }

    trie.Close()
}
