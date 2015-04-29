package imap

import (
	"errors"
	"fmt"
	Mmap "github.com/yqingp/lsearch/store/mmap"
	"log"
	"os"
	"sync"
	"unsafe"
)

const (
	SlotMax           = 100000
	SlotIncrStep      = 2000
	SlotNum           = 1024
	Slot2Num          = 512
	NodeValueIncrStep = 1000000
	NodeMax           = 100000000
	SizeOfState       = int64(unsafe.Sizeof(State{}))
	SizeOfNode        = int64(unsafe.Sizeof(Node{}))
	SizeOfNodeValue   = int64(unsafe.Sizeof(NodeValue{}))
)

type DB struct {
	state      *State
	nodes      []Node
	nodeValues []NodeValue
	slots      []*Slot
	roots      []uint32
	file       *os.File
	vfile      *os.File
	size       int64
	msize      int64
	vmsize     int64
	mutex      *sync.Mutex
	mmap       Mmap.Mmap
	vsize      int64
}

var Logger *log.Logger = log.New(os.Stdout, "imap", log.Lshortfile|log.Ltime)

func Open(filePath string) *DB {
	if filePath == "" {
		Logger.Fatal("file path is blank")
	}

	db := &DB{}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		Logger.Fatal("open file error")
	}
	db.file = file

	size := SizeOfState + SizeOfNode*NodeMax
	db.msize = size

	fd := int(file.Fd())

	mmap, err := Mmap.MmapFile(fd, int(size))

	if err != nil {
		Logger.Fatal(err)
	}
	db.state = (*State)(unsafe.Pointer(&mmap[0]))

	addr := unsafe.Pointer(&mmap[SizeOfState])
	db.nodes = (*[NodeMax]Node)(addr)[:NodeMax]

	fstat, err := file.Stat()
	if err != nil {
		Logger.Fatal(err)
	}

	db.size = fstat.Size()

	if db.size < SizeOfState {
		file.Truncate(SizeOfState)
		db.size = SizeOfState

		for i := 0; i < SlotMax; i++ {
			db.state.slots[i].nodeId = -1
		}

	}

	vfilePath := fmt.Sprintf("%s.v", filePath)

	vfile, err := os.OpenFile(vfilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		Logger.Fatal(err)
	}

	db.vfile = vfile

	db.vmsize = SizeOfNodeValue * NodeMax
	size = db.vmsize

	mmap, err = Mmap.MmapFile(int(vfile.Fd()), int(size))
	if err != nil {
		Logger.Fatal(err)
	}

	addr = unsafe.Pointer(&mmap[0])
	db.nodeValues = (*[NodeMax]NodeValue)(addr)[:NodeMax]
	fstat, err = vfile.Stat()
	if err != nil {
		Logger.Fatal(err)
	}

	db.vsize = fstat.Size()
	db.mutex = &sync.Mutex{}

	return db
}

func (d *DB) Set(num uint32, key int32) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.vSet(num, key)

	if d.nodeValues[num].off < 0 {
		d.insert(num, key)
	} else {
		if key != d.nodeValues[num].val {
			d.remove(num)
			d.insert(num, key)
		}
	}

	d.nodeValues[num].val = key
}

func (d *DB) vSet(num uint32, key int32) {
	size := (int64(num)/NodeValueIncrStep + 1) * NodeValueIncrStep * SizeOfNodeValue
	n, i := int64(0), int64(0)

	if d.state != nil && num >= 0 && num < NodeMax {

		if size > d.vsize {
			err := d.vfile.Truncate(size)
			if err != nil {
				Logger.Fatal(err)
			}

			i = d.vsize / SizeOfNodeValue
			n = size / SizeOfNodeValue

			for i < n {
				d.nodeValues[i].off = -1
				d.nodeValues[i].val = 0
				i++
			}
		}
	}

}

func (d *DB) insert(num uint32, key int32) {

}

func (d *DB) remove(num uint32) {

}

func (d *DB) Get() {

}

func (d *DB) Del() {

}

func (d *DB) Range() {

}

func (d *DB) From() {

}

func (d *DB) To() {

}

//支持多个值
func (d *DB) In(keys []uint32) {

}
