package util

import (
    "errors"
    "fmt"
    . "github.com/yqingp/lsearch/mmap"
    "os"
    "sync"
    "unsafe"
)

const (
    MMTRIE_PATH_MAX  = 256
    MMTRIE_LINE_MAX  = 256
    MMTRIE_BASE_NUM  = 1000
    MMTRIE_NODES_MAX = 100000
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

type MmtrieNode struct {
    key     uint8
    nchilds uint8
    data    int
    childs  int
}

type Mmtrie struct {
    state    *MmtrieState
    nodes    []MmtrieNode
    mmap     Mmap
    size     int
    old      int64
    filesize int64
    fd       int
    bits     int
    mutex    *sync.Mutex
    isInit   bool
    filename string
}

var (
    t                 MmtrieState
    t1                MmtrieNode
    MmtrieNodeSizeOf  = int(unsafe.Sizeof(t1))
    MmtrieStateSizeOf = int(unsafe.Sizeof(t))
)

func NewMmtrie(filename string) (*Mmtrie, error) {
    if filename == "" {
        return nil, errors.New("file name is blank")
    }
    return &Mmtrie{filename: filename}, nil
}

func (m *Mmtrie) Init() error {
    if m.isInit {
        return errors.New("is inited")
    }

    f, err := os.OpenFile(m.filename, os.O_CREATE|os.O_RDWR, 0644)
    if err != nil {
        return err
    }
    m.fd = int(f.Fd())

    fstat, err := os.Stat(m.filename)
    if err != nil {
        return err
    }

    if m.mmap == nil {
        m.size = MmtrieStateSizeOf + MMTRIE_NODES_MAX*MmtrieNodeSizeOf
        mp, err := MmapFile(m.fd, m.size)
        if err != nil {
            return err
        }
        m.mmap = mp
        m.state = (*MmtrieState)(unsafe.Pointer(&m.mmap[0]))
        addr := &m.mmap[MmtrieStateSizeOf]
        pointer := (*[MMTRIE_NODES_MAX]MmtrieNode)(unsafe.Pointer(addr))[:MMTRIE_NODES_MAX]
        m.nodes = pointer
    }

    if fstat.Size() == 0 {
        m.filesize = int64(MmtrieStateSizeOf) + MMTRIE_BASE_NUM*int64(MmtrieNodeSizeOf)
        if err := f.Truncate(m.filesize); err != nil {
            return err
        }
        m.state.total = MMTRIE_BASE_NUM
        m.state.left = MMTRIE_BASE_NUM - MMTRIE_LINE_MAX
        m.state.current = MMTRIE_LINE_MAX
    }

    m.mutex = &sync.Mutex{}
    return nil
}

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

func (self *Mmtrie) increment() error {
    if self.filesize < int64(self.size) {
        self.old = self.filesize
        self.filesize += int64(MMTRIE_BASE_NUM * MmtrieNodeSizeOf)
        if err := os.Truncate(self.filename, self.filesize); err != nil {
            return err
        }
        self.state.total += MMTRIE_BASE_NUM
        self.state.left += MMTRIE_BASE_NUM
    }

    return nil
}

func (self *MmtrieNode) setKey(k byte) {
    self.key = k
}

func (self *MmtrieNode) nodeCopy(old MmtrieNode) {
    self.childs = old.childs
    self.data = old.data
    self.key = old.key
    self.nchilds = old.nchilds
}

func (self *Mmtrie) Add(key []byte) int {
    ret := -1

    if key == nil {
        return ret
    }

    self.mutex.Lock()
    defer self.mutex.Unlock()

    i := key[0]
    m := 1

    n := 0
    z := 0
    j := 0
    k := 0

    pos := 0

    _ = n
    _ = z
    _ = pos
    _ = j
    _ = k

    childs := &MmtrieNode{}
    _ = childs

    size := len(key)
    for m < size {
        x := 0
        if self.nodes[i].nchilds > 0 && self.nodes[i].childs >= MMTRIE_LINE_MAX {

        }
        if x < MMTRIE_LINE_MAX || self.nodes[x].key != key[m] {
            n = int(self.nodes[i].nchilds) + 1
            z = self.nodes[i].childs
            pos, _ = self.pop(n)
            if pos < MMTRIE_LINE_MAX || pos > self.state.current {
                return ret
            }
            childs = &(self.nodes[pos])
            if x == 0 {
                childs.setKey(key[m])
                j = pos
            } else if x == -1 {
                childs.setKey(key[m])
                k = 1
                for k < n {
                    tc := &(self.nodes[pos+k])
                    tc.nodeCopy(self.nodes[z])
                    z++
                    k++
                }
                j = pos
            }
        }

        m++
    }

    if self.nodes[i].data == 0 {
        self.state.id++
        self.nodes[i].data = self.state.id
        ret = self.nodes[i].data
    }

    return ret
}

func (m *Mmtrie) ToS() {
    fmt.Println(len(m.nodes))
}
