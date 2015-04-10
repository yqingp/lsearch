package store

type TrieNode struct {
    key        uint8
    childCount uint8
    data       int
    childPos   int
}

func (t *TrieNode) setKey(k byte) {
    t.key = k
    t.childCount = 0
    t.childPos = 0
    t.data = 0
}

func (t *TrieNode) nodeCopy(old TrieNode) {
    t.childPos = old.childPos
    t.data = old.data
    t.key = old.key
    t.childCount = old.childCount
}
