package store

import (
    "fmt"
    "os"
    "strconv"
    "testing"
)

func TestTrieAll(t *testing.T) {
    trie, err := OpenTrie("a.txt")
    if err != nil {
        t.Error(err)
    }

    total := 10000
    for i := 0; i < total; i++ {
        m := strconv.Itoa(i) + "哈哈" + strconv.Itoa(i+1)
        _, err := trie.Set([]byte(m))
        if err != nil {
            t.Fatal(err)
        }
    }

    fmt.Println(trie.Get([]byte("981哈哈982")))
    fmt.Println(trie.Del([]byte("981哈哈982")))
    fmt.Println(trie.Get([]byte("981哈哈982")))

    trie.Close()
    os.Remove("a.txt")
}
