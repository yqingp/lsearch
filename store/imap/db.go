package imap

import (
	"os"
	"sync"
)

const (
	SlotMax           = 100000
	SlotIncrStep      = 2000
	SlotNum           = 1024
	Slot2Num          = 512
	NodeValueIncrStep = 1000000
	NodeMax           = 100000000
)

type DB struct {
	state      *State
	nodes      []*Node
	nodeValues []*NodeValue
	slots      []*Slot
	roots      []uint32
	file       *os.File
	vfile      *os.File
	size       uint32
	mutex      *sync.Mutex
}

func Open(filePath string) error {

	return nil
}

func (d *DB) Set() {

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
