package store

import (
// "fmt"
// "strconv"
// "testing"
)

// func TestMmtrieGet(t *testing.T) {
//     trie, err := OpenTrie("a.txt")
//     if err != nil {
//         t.Error(err)
//     }

//     fmt.Println(trie.Get([]byte("932哈哈933")))
//     trie.Close()
// }

// func TestMmtrieDel(t *testing.T) {
//     trie, err := OpenTrie("a.txt")
//     if err != nil {
//         t.Error(err)
//     }
//     fmt.Println(trie.Get([]byte("981哈哈982")))
//     fmt.Println(trie.Del([]byte("981哈哈982")))
//     fmt.Println(trie.Get([]byte("981哈哈982")))
//     trie.Close()
// }

// func TestMmtrieSet(t *testing.T) {
//     trie, err := OpenTrie("a.txt")
//     if err != nil {
//         t.Error(err)
//     }

//     for i := 1; i < 1000000; i++ {
//         m := strconv.Itoa(i) + "哈哈" + strconv.Itoa(i+1)
//         // fmt.Println(m)
//         _, err := trie.Set([]byte(m))
//         if err != nil {
//             t.Fatal(err)
//         }
//     }

//     trie.Close()
// }
