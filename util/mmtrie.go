package util

import (
    "errors"
    . "github.com/yqingp/lsearch/mmap"
    "os"
    "sync"
    "unsafe"
)

const (
    MMTRIE_PATH_MAX  = 256
    MMTRIE_LINE_MAX  = 256
    MMTRIE_BASE_NUM  = 1000000
    MMTRIE_NODES_MAX = 100000000
    MMTRIE_WORD_MAX  = 4096
)

type MmtrieState struct {
    id      int
    current int
    total   int
    left    int
    list    [MMTRIE_LINE_MAX]MmtrieList
}

type MmtrieList struct {
    count int
    head  int
}

type Mmtrie struct {
    state    *MmtrieState
    nodes    []MmtrieNode
    mmap     Mmap
    size     int64
    old      int64
    filesize int64
    fd       int
    file     *os.File
    bits     int
    mutex    *sync.Mutex
    filename string
}

var (
    MmtrieStateSizeOf = int64(unsafe.Sizeof(MmtrieState{}))
    MmtrieNodeSizeOf  = int64(unsafe.Sizeof(MmtrieNode{}))
)

func Open(filename string) (*Mmtrie, error) {
    if filename == "" {
        return nil, errors.New("file name is blank")
    }

    m := &Mmtrie{filename: filename}

    if err := m.init(); err != nil {
        return nil, err
    }

    return m, nil
}

func (self *Mmtrie) Close() {
    if self.mmap != nil {
        self.mmap.Unmap()
    }

    if self.file != nil {
        self.file.Close()
    }
}

func (m *Mmtrie) init() error {
    f, err := os.OpenFile(m.filename, os.O_CREATE|os.O_RDWR, 0644)
    if err != nil {
        return err
    }
    m.file = f
    m.fd = int(f.Fd())

    fstat, err := f.Stat()
    if err != nil {
        return err
    }

    m.filesize = fstat.Size()

    if m.mmap == nil {
        m.size = MmtrieStateSizeOf + MMTRIE_NODES_MAX*MmtrieNodeSizeOf
        mp, err := MmapFile(m.fd, int(m.size))
        if err != nil {
            return err
        }
        m.mmap = mp
        m.state = (*MmtrieState)(unsafe.Pointer(&m.mmap[0]))
        addr := &m.mmap[MmtrieStateSizeOf]
        pointer := (*[MMTRIE_NODES_MAX]MmtrieNode)(unsafe.Pointer(addr))[:MMTRIE_NODES_MAX]
        m.nodes = pointer
    }

    if m.filesize == 0 {
        m.filesize = MmtrieStateSizeOf + MMTRIE_BASE_NUM*MmtrieNodeSizeOf
        if err := m.file.Truncate(m.filesize); err != nil {
            return err
        }

        m.state.total = MMTRIE_BASE_NUM
        m.state.left = MMTRIE_BASE_NUM - MMTRIE_LINE_MAX
        m.state.current = MMTRIE_LINE_MAX
    }

    m.mutex = &sync.Mutex{}
    return nil
}

func (self *Mmtrie) Set(key []byte) (int, error) {
    ret := -1

    if key == nil {
        return ret, errors.New("key is blank")
    }

    self.mutex.Lock()
    defer self.mutex.Unlock()

    i := int(key[0])

    m, n, z, j, k, min, max, pos, x := 1, 0, 0, 0, 0, 0, 0, 0, 0

    size := len(key)

    var err error

    for m < size {
        x = 0
        if self.nodes[i].nchilds > 0 && self.nodes[i].childs >= MMTRIE_LINE_MAX {
            min = self.nodes[i].childs
            max = min + int(self.nodes[i].nchilds) - 1
            if key[m] == self.nodes[min].key {
                x = min
            } else if key[m] == self.nodes[max].key {
                x = max
            } else if key[m] < self.nodes[min].key {
                x = -1
            } else if key[m] > self.nodes[max].key {
                x = 1
            } else {
                for max > min {
                    z = (max + min) / 2
                    if z == min {
                        x = z
                        break
                    }
                    if self.nodes[z].key == key[m] {
                        x = z
                        break
                    } else if key[m] > self.nodes[z].key {
                        min = z
                    } else {
                        max = z
                    }
                }
            }
        }
        if x < MMTRIE_LINE_MAX || self.nodes[x].key != key[m] {
            n = int(self.nodes[i].nchilds) + 1
            z = self.nodes[i].childs
            pos, err = self.pop(n)
            if err != nil {
                return -1, err
            }
            if pos < MMTRIE_LINE_MAX || pos > self.state.current {
                return -1, nil
            }
            if x == 0 {
                self.nodes[pos].setKey(key[m])
                j = pos
            } else if x == -1 {
                self.nodes[pos].setKey(key[m])
                k = 1
                for k < n {
                    self.nodes[pos+k].nodeCopy(self.nodes[z])
                    z++
                    k++
                }
                j = pos
            } else if x == 1 {
                k = 0
                for k < (n - 1) {
                    self.nodes[pos+k].nodeCopy(self.nodes[z])
                    z++
                    k++
                }
                self.nodes[pos+k].setKey(key[m])
                j = pos + k
            } else {
                k = 0
                for (self.nodes[z].key) < key[m] {
                    self.nodes[pos+k].nodeCopy(self.nodes[z])
                    z++
                    k++
                }
                self.nodes[pos+k].setKey(key[m])
                x = k
                k++
                for k < n {
                    self.nodes[pos+k].nodeCopy(self.nodes[z])
                    z++
                    k++
                }
                j = pos + x
            }

            self.push(int(self.nodes[i].nchilds), self.nodes[i].childs)
            self.nodes[i].nchilds++
            self.nodes[i].childs = pos
            i = j
        } else {
            i = x
        }

        m++
    }

    ret = self.nodes[i].data

    if ret == 0 {
        self.state.id++
        self.nodes[i].data = self.state.id
        ret = self.nodes[i].data
    }

    return ret, nil
}

// if not found, return 0
func (self *Mmtrie) Get(key []byte) (int, error) {
    ret := 0

    if key == nil {
        return ret, errors.New("key is blank")
    }

    self.mutex.Lock()
    defer self.mutex.Unlock()

    i := int(key[0])

    m, z, min, max, x := 1, 0, 0, 0, 0

    size := len(key)

    if size == 1 && i >= 0 && i < self.state.total {
        return self.nodes[i].data, nil
    }

    for m < size {
        x = 0
        if self.nodes[i].nchilds > 0 && self.nodes[i].childs >= MMTRIE_LINE_MAX {
            min = self.nodes[i].childs
            max = min + int(self.nodes[i].nchilds) - 1
            if key[m] == self.nodes[min].key {
                x = min
            } else if key[m] == self.nodes[max].key {
                x = max
            } else if key[m] < self.nodes[min].key {
                return ret, nil
            } else if key[m] > self.nodes[max].key {
                return ret, nil
            } else {
                for max > min {
                    z = (max + min) / 2
                    if z == min {
                        x = z
                        break
                    }
                    if self.nodes[z].key == key[m] {
                        x = z
                        break
                    } else if self.nodes[z].key < key[m] {
                        min = z
                    } else {
                        max = z
                    }
                }
                if self.nodes[x].key != key[m] {
                    return ret, nil
                }
            }
            i = x
        }

        if i >= 0 && i < self.state.total && (self.nodes[i].nchilds == 0 || (m+1 == size)) {
            if self.nodes[i].key != key[m] {
                return ret, nil
            }
            if m+1 == size {
                return self.nodes[i].data, nil
            }
            break
        }
        m++
    }
    return ret, nil
}

//if not found return 0 ,  else return val and  remove it
func (self *Mmtrie) Del(key []byte) (int, error) {
    ret := 0

    if key == nil {
        return ret, errors.New("key is blank")
    }

    self.mutex.Lock()
    defer self.mutex.Unlock()

    i := int(key[0])

    m, z, min, max, x := 1, 0, 0, 0, 0

    size := len(key)

    if size == 1 && i >= 0 && i < self.state.total && self.nodes[i].data != 0 {
        ret = self.nodes[i].data
        self.nodes[i].data = 0
        return ret, nil
    }

    for m < size {
        x = 0
        if self.nodes[i].nchilds > 0 && self.nodes[i].childs >= MMTRIE_LINE_MAX {
            min = self.nodes[i].childs
            max = min + int(self.nodes[i].nchilds) - 1
            if key[m] == self.nodes[min].key {
                x = min
            } else if key[m] == self.nodes[max].key {
                x = max
            } else if key[m] < self.nodes[min].key {
                return ret, nil
            } else if key[m] > self.nodes[max].key {
                return ret, nil
            } else {
                for max > min {
                    z = (max + min) / 2
                    if z == min {
                        x = z
                        break
                    }
                    if self.nodes[z].key == key[m] {
                        x = z
                        break
                    } else if self.nodes[z].key < key[m] {
                        min = z
                    } else {
                        max = z
                    }
                }
                if self.nodes[x].key != key[m] {
                    return ret, nil
                }
            }
            i = x
        }

        if i >= 0 && i < self.state.total && (self.nodes[i].nchilds == 0 || (m+1 == size)) {
            if self.nodes[i].key != key[m] {
                return ret, nil
            }
            if m+1 == size {
                ret = self.nodes[i].data
                self.nodes[i].data = 0
                return ret, nil
            }
            break
        }
        m++
    }
    return ret, nil
}

//pop one node
func (self *Mmtrie) pop(num int) (int, error) {
    pos := -1

    if num > 0 && num <= MMTRIE_LINE_MAX && self.state != nil && self.nodes != nil {
        if self.state.list[num-1].count > 0 {
            pos = self.state.list[num-1].head
            self.state.list[num-1].head = self.nodes[pos].childs
            self.state.list[num-1].count--
        } else {
            if self.state.left < num {
                if err := self.increment(); err != nil {
                    return pos, err
                }
            }
            pos = self.state.current
            self.state.current += num
            self.state.left -= num
        }
    }

    return pos, nil
}

//push one node to free list
func (self *Mmtrie) push(num int, pos int) {
    if pos >= MMTRIE_LINE_MAX && num > 0 && num <= MMTRIE_LINE_MAX {
        self.nodes[pos].childs = self.state.list[num-1].head
        self.state.list[num-1].head = pos
        self.state.list[num-1].count++
    }
}

// check file need truncate
func (self *Mmtrie) increment() error {

    if self.filesize < int64(self.size) {
        self.old = self.filesize
        self.filesize += int64(MMTRIE_BASE_NUM) * int64(MmtrieNodeSizeOf)
        if err := self.file.Truncate(self.filesize); err != nil {
            return err
        }
        self.state.total += MMTRIE_BASE_NUM
        self.state.left += MMTRIE_BASE_NUM
    }

    return nil
}
