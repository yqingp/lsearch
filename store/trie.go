package store

import (
    "errors"
    "os"
    "sync"
    "unsafe"
)

const (
    MaxLineCount  = 256
    BaseNodeCount = 1000000
    MaxNodeCount  = 100000000
)

type TrieState struct {
    id       int
    current  int
    total    int
    left     int
    list     [MaxLineCount]TrieList
    totalNum int
}

type TrieList struct {
    count int
    head  int
}

type Trie struct {
    state    *TrieState
    nodes    []TrieNode
    mmap     Mmap
    size     int64
    old      int64
    fileSize int64
    file     *os.File
    mutex    *sync.Mutex
    filePath string
}

var (
    SizeOfTrieState = int64(unsafe.Sizeof(TrieState{}))
    SizeOfTrieNode  = int64(unsafe.Sizeof(TrieNode{}))
)

func OpenTrie(filePath string) (*Trie, error) {
    if filePath == "" {
        return nil, errors.New("file name is blank")
    }

    m := &Trie{filePath: filePath}

    if err := m.init(); err != nil {
        return nil, err
    }

    return m, nil
}

func (t *Trie) Close() {
    if t.mmap != nil {
        t.mmap.Unmap()
    }

    if t.file != nil {
        t.file.Close()
    }
}

func (m *Trie) init() error {
    f, err := os.OpenFile(m.filePath, os.O_CREATE|os.O_RDWR, 0644)
    if err != nil {
        return err
    }
    m.file = f

    fstat, err := f.Stat()
    if err != nil {
        return err
    }

    fd := int(f.Fd())
    m.fileSize = fstat.Size()

    if m.mmap == nil {
        m.size = SizeOfTrieState + (MaxNodeCount * SizeOfTrieNode)
        mp, err := MmapFile(fd, int(m.size))
        if err != nil {
            return err
        }
        m.mmap = mp
        m.state = (*TrieState)(unsafe.Pointer(&m.mmap[0]))
        addr := &m.mmap[SizeOfTrieState]
        pointer := (*[MaxNodeCount]TrieNode)(unsafe.Pointer(addr))[:MaxNodeCount]
        m.nodes = pointer
    }

    if m.fileSize == 0 {
        m.fileSize = SizeOfTrieState + BaseNodeCount*SizeOfTrieNode
        if err := m.file.Truncate(m.fileSize); err != nil {
            return err
        }

        m.state.total = BaseNodeCount
        m.state.left = BaseNodeCount - MaxLineCount
        m.state.current = MaxLineCount
    }

    m.mutex = &sync.Mutex{}
    return nil
}

func (t *Trie) Set(key []byte) (int, error) {
    ret := -1

    if key == nil {
        return ret, errors.New("key is blank")
    }

    t.mutex.Lock()
    defer t.mutex.Unlock()

    i := int(key[0])

    m, n, z, j, k, min, max, pos, x := 1, 0, 0, 0, 0, 0, 0, 0, 0

    size := len(key)

    var err error

    for m < size {
        x = 0
        if t.nodes[i].childCount > 0 && t.nodes[i].childPos >= MaxLineCount {
            min = t.nodes[i].childPos
            max = min + int(t.nodes[i].childCount) - 1
            if key[m] == t.nodes[min].key {
                x = min
            } else if key[m] == t.nodes[max].key {
                x = max
            } else if key[m] < t.nodes[min].key {
                x = -1
            } else if key[m] > t.nodes[max].key {
                x = 1
            } else {
                for max > min {
                    z = (max + min) / 2
                    if z == min {
                        x = z
                        break
                    }
                    if t.nodes[z].key == key[m] {
                        x = z
                        break
                    } else if key[m] > t.nodes[z].key {
                        min = z
                    } else {
                        max = z
                    }
                }
            }
        }
        if x < MaxLineCount || t.nodes[x].key != key[m] {
            n = int(t.nodes[i].childCount) + 1
            z = t.nodes[i].childPos
            pos, err = t.pop(n)
            if err != nil {
                return -1, err
            }
            if pos < MaxLineCount || pos > t.state.current {
                return -1, nil
            }
            if x == 0 {
                t.nodes[pos].setKey(key[m])
                j = pos
            } else if x == -1 {
                t.nodes[pos].setKey(key[m])
                k = 1
                for k < n {
                    t.nodes[pos+k].nodeCopy(t.nodes[z])
                    z++
                    k++
                }
                j = pos
            } else if x == 1 {
                k = 0
                for k < (n - 1) {
                    t.nodes[pos+k].nodeCopy(t.nodes[z])
                    z++
                    k++
                }
                t.nodes[pos+k].setKey(key[m])
                j = pos + k
            } else {
                k = 0
                for (t.nodes[z].key) < key[m] {
                    t.nodes[pos+k].nodeCopy(t.nodes[z])
                    z++
                    k++
                }
                t.nodes[pos+k].setKey(key[m])
                x = k
                k++
                for k < n {
                    t.nodes[pos+k].nodeCopy(t.nodes[z])
                    z++
                    k++
                }
                j = pos + x
            }

            t.push(int(t.nodes[i].childCount), t.nodes[i].childPos)
            t.nodes[i].childCount++
            t.nodes[i].childPos = pos
            i = j
        } else {
            i = x
        }

        m++
    }

    ret = t.nodes[i].data

    if ret == 0 {
        t.state.id++
        t.nodes[i].data = t.state.id
        ret = t.nodes[i].data
        t.state.totalNum++
    }

    return ret, nil
}

// if not found, return 0
func (t *Trie) Get(key []byte) (int, error) {
    ret := 0

    if key == nil {
        return ret, errors.New("key is blank")
    }

    t.mutex.Lock()
    defer t.mutex.Unlock()

    i := int(key[0])

    m, z, min, max, x := 1, 0, 0, 0, 0

    size := len(key)

    if size == 1 && i >= 0 && i < t.state.total {
        return t.nodes[i].data, nil
    }

    for m < size {
        x = 0
        if t.nodes[i].childCount > 0 && t.nodes[i].childPos >= MaxLineCount {
            min = t.nodes[i].childPos
            max = min + int(t.nodes[i].childCount) - 1
            if key[m] == t.nodes[min].key {
                x = min
            } else if key[m] == t.nodes[max].key {
                x = max
            } else if key[m] < t.nodes[min].key {
                return ret, nil
            } else if key[m] > t.nodes[max].key {
                return ret, nil
            } else {
                for max > min {
                    z = (max + min) / 2
                    if z == min {
                        x = z
                        break
                    }
                    if t.nodes[z].key == key[m] {
                        x = z
                        break
                    } else if t.nodes[z].key < key[m] {
                        min = z
                    } else {
                        max = z
                    }
                }
                if t.nodes[x].key != key[m] {
                    return ret, nil
                }
            }
            i = x
        }

        if i >= 0 && i < t.state.total && (t.nodes[i].childCount == 0 || (m+1 == size)) {
            if t.nodes[i].key != key[m] {
                return ret, nil
            }
            if m+1 == size {
                return t.nodes[i].data, nil
            }
            break
        }
        m++
    }
    return ret, nil
}

//if not found return 0 ,  else return val and  remove it
func (t *Trie) Del(key []byte) (int, error) {
    ret := 0

    if key == nil {
        return ret, errors.New("key is blank")
    }

    t.mutex.Lock()
    defer t.mutex.Unlock()

    i := int(key[0])

    m, z, min, max, x := 1, 0, 0, 0, 0

    size := len(key)

    if size == 1 && i >= 0 && i < t.state.total && t.nodes[i].data != 0 {
        ret = t.nodes[i].data
        t.nodes[i].data = 0
        t.state.totalNum--
        return ret, nil
    }

    for m < size {
        x = 0
        if t.nodes[i].childCount > 0 && t.nodes[i].childPos >= MaxLineCount {
            min = t.nodes[i].childPos
            max = min + int(t.nodes[i].childCount) - 1
            if key[m] == t.nodes[min].key {
                x = min
            } else if key[m] == t.nodes[max].key {
                x = max
            } else if key[m] < t.nodes[min].key {
                return ret, nil
            } else if key[m] > t.nodes[max].key {
                return ret, nil
            } else {
                for max > min {
                    z = (max + min) / 2
                    if z == min {
                        x = z
                        break
                    }
                    if t.nodes[z].key == key[m] {
                        x = z
                        break
                    } else if t.nodes[z].key < key[m] {
                        min = z
                    } else {
                        max = z
                    }
                }
                if t.nodes[x].key != key[m] {
                    return ret, nil
                }
            }
            i = x
        }

        if i >= 0 && i < t.state.total && (t.nodes[i].childCount == 0 || (m+1 == size)) {
            if t.nodes[i].key != key[m] {
                return ret, nil
            }
            if m+1 == size {
                ret = t.nodes[i].data
                t.nodes[i].data = 0
                t.state.totalNum--
                return ret, nil
            }
            break
        }
        m++
    }
    return ret, nil
}

//pop one node
func (t *Trie) pop(num int) (int, error) {
    pos := -1

    if num > 0 && num <= MaxLineCount && t.state != nil && t.nodes != nil {
        if t.state.list[num-1].count > 0 {
            pos = t.state.list[num-1].head
            t.state.list[num-1].head = t.nodes[pos].childPos
            t.state.list[num-1].count--
        } else {
            if t.state.left < num {
                if err := t.increment(); err != nil {
                    return pos, err
                }
            }
            pos = t.state.current
            t.state.current += num
            t.state.left -= num
        }
    }

    return pos, nil
}

//push one node to free list
func (t *Trie) push(num int, pos int) {
    if pos >= MaxLineCount && num > 0 && num <= MaxLineCount {
        t.nodes[pos].childPos = t.state.list[num-1].head
        t.state.list[num-1].head = pos
        t.state.list[num-1].count++
    }
}

// check file need truncate
func (t *Trie) increment() error {

    if t.fileSize < int64(t.size) {
        t.old = t.fileSize
        t.fileSize += int64(BaseNodeCount) * int64(SizeOfTrieNode)
        if err := t.file.Truncate(t.fileSize); err != nil {
            return err
        }
        t.state.total += BaseNodeCount
        t.state.left += BaseNodeCount
    }

    return nil
}
