package idb

import (
	"errors"
	. "github.com/yqingp/lsearch/store/trie"
	"os"
	"sync"
	"time"
	// "bytes"
	// "encoding/gob"
)

type DB struct {
	status          int
	mutex           *sync.Mutex
	blockQueueMutex *sync.Mutex
	indexMutex      *sync.Mutex
	mutexs          [MaxMutexCount]*sync.Mutex
	state           *State
	stateIO         IO
	blockQueueIO    IO
	indexIO         IO
	blockQueues     []BlockQueue
	indexes         []Index
	IOs             [MaxDbFileCount]IO
	baseDir         string
	keyMapTrie      *Trie
	loggerFile      *os.File
	isMmap          bool
}

func Open(baseDir string, isMmap bool) (*DB, error) {
	if baseDir == "" {
		return nil, errors.New("basedir is blank")
	}

	db := &DB{}

	db.blockQueueMutex = &sync.Mutex{}
	db.indexMutex = &sync.Mutex{}
	db.mutex = &sync.Mutex{}
	db.baseDir = baseDir
	db.isMmap = isMmap

	if err := db.initDir(); err != nil {
		return nil, err
	}

	if err := db.initLogger(); err != nil {
		return nil, err
	}

	if err := db.initKmap(); err != nil {

	}

	db.initState()
	db.initBlockQueue()
	db.initIndex()
	db.initIOs()

	return db, nil
}

func (d *DB) init() {

}

func (d DB) RecordNum() int {
	return d.keyMapTrie.TotalNum()
}

func (d *DB) Close() {
	if d.loggerFile != nil {
		d.loggerFile.Close()
	}
	if d.keyMapTrie != nil {
		d.keyMapTrie.Close()
	}

	d.stateIO.close()
	d.indexIO.close()
	d.blockQueueIO.close()
	for _, v := range d.IOs {
		v.close()
	}
}

func (d *DB) CheckAndSetId(key []byte) (int, bool, error) {

	id, err := d.keyMapTrie.Get(key)
	if err != nil {
		return id, false, err
	}

	if id == 0 {
		id, err := d.keyMapTrie.Set(key)
		if err != nil {
			return id, false, err
		}

		return id, false, nil
	}

	return id, true, nil
}

// if id is < 1  generate auto increment id
func (d *DB) Set(id int, key []byte, value []byte) (int, error) {
	if key == nil || value == nil || len(value) == 0 {
		return -1, errors.New("key or value is blank")
	}

	var err error

	if id < 1 {
		id, err = d.keyMapTrie.Set(key)
		if err != nil {
			return -1, err
		}
	}

	ret := d.internalSet(id, value)
	if ret == -1 {
		Logger.Fatal("set error")
	}

	return id, nil
}

func (d *DB) internalSet(id int, value []byte) int {
	ret := -1
	indexes := d.indexes
	if d.status != 0 || indexes == nil {
		return ret
	}

	valueLen := len(value)

	blocksCountNum, index := 0, 0

	d.indexMutex.Lock()
	d.checkIndexIOWithId(id)
	d.indexMutex.Unlock()

	d.lockId(id)
	defer d.unlockId(id)

	oldBlockQueue := &BlockQueue{}
	var newBlockQueue *BlockQueue

	if indexes[id].blockSize < valueLen {
		if indexes[id].blockSize > 0 {
			oldBlockQueue.index = indexes[id].index
			oldBlockQueue.blockId = indexes[id].blockId
			oldBlockQueue.count = blocksCount(indexes[id].blockSize)
			indexes[id].blockSize = 0
			indexes[id].blockId = 0
			indexes[id].dataLen = 0
		}

		blocksCountNum = blocksCount(valueLen)
		newBlockQueue = d.popBlockQueue(blocksCountNum)
		if newBlockQueue != nil {
			indexes[id].index = newBlockQueue.index
			indexes[id].blockId = newBlockQueue.blockId
			indexes[id].blockSize = blocksCountNum * BaseDbSize
			if valueLen > indexes[id].blockSize {
				Logger.Fatal("Invalid  block")
			}
		} else {
			Logger.Fatal("pop block error")
		}
	}

	if indexes[id].blockSize >= valueLen && indexes[id].index >= 0 &&
		d.IOs[index].file != nil {

		index = indexes[id].index
		if d.isMmap && indexes[id].blockId >= 0 && d.IOs[index].mmap != nil {
			for k, v := range value {
				d.IOs[index].mmap[indexes[id].blockId*BaseDbSize+k] = v
			}

			indexes[id].dataLen = valueLen
			ret = id
		} else {
			_, err := d.IOs[index].file.WriteAt(value, int64(indexes[id].blockId*BaseDbSize))
			if err != nil {
				indexes[id].dataLen = 0
				Logger.Fatal("write index error")
			}

			indexes[id].dataLen = valueLen
			ret = id
		}
	}

	if indexes[id].dataLen > d.state.dataLenMax {
		d.state.dataLenMax = indexes[id].dataLen
	}

	indexes[id].updateTime = time.Now().Unix()
	if oldBlockQueue.count > 0 {
		d.pushBlockQueue(oldBlockQueue.index, oldBlockQueue.blockId, oldBlockQueue.count*BaseDbSize)
	}

	return ret
}

func (d *DB) lockId(id int) {
	d.mutexs[id%MaxMutexCount].Lock()
}

func (d *DB) unlockId(id int) {
	d.mutexs[id%MaxMutexCount].Unlock()
}

func (d *DB) Get(key []byte) (value []byte, ret int) {
	if key == nil {
		return
	}

	id, err := d.keyMapTrie.Get(key)
	if err != nil {
		return nil, -1
	}

	return d.GetByInternalId(id)
}

func (d *DB) GetAndReturnInternalId(key []byte) ([]byte, int) {
	if key == nil {
		return nil, -1
	}

	id, err := d.keyMapTrie.Get(key)
	if err != nil {
		return nil, -1
	}
	// Logger.Println(id)

	val, _ := d.GetByInternalId(id)

	return val, id
}

// get by internal integer ID, if found return val ,otherwise -1
func (d *DB) GetByInternalId(id int) (value []byte, ret int) {
	if id <= 0 || id > d.state.maxId {
		return nil, -1
	}
	d.lockId(id)
	defer d.unlockId(id)

	indexes := d.indexes
	if indexes == nil {
		Logger.Fatal("db index error")
	}

	dataLen := indexes[id].dataLen

	blockId := indexes[id].blockId
	index := indexes[id].index

	if dataLen > 0 && indexes[id].blockSize > 0 && blockId >= 0 &&
		index >= 0 && d.IOs[index].file != nil {

		offsetSize := blockId * BaseDbSize
		if d.isMmap && d.IOs[index].mmap != nil {
			value = d.IOs[index].mmap[offsetSize:(offsetSize + dataLen)]
			return value, 0
		} else {
			value = make([]byte, dataLen)
			readSize, err := d.IOs[index].file.ReadAt(value[:dataLen], int64(offsetSize))
			if err != nil || readSize != dataLen {
				Logger.Fatal("read index error")
			}
		}
	}

	return value, 0
}

func (d *DB) Add(id int, value []byte) {
	if id < 1 || value == nil {
		return
	}

	ret := -1

	_ = ret
	indexes := d.indexes
	if d.status != 0 || indexes == nil {
		return
	}

	valueLen := len(value)

	blocksCountNum, index := 0, 0

	d.indexMutex.Lock()
	d.checkIndexIOWithId(id)
	d.indexMutex.Unlock()

	d.lockId(id)
	defer d.unlockId(id)

	oldBlockQueue := &BlockQueue{}
	var newBlockQueue *BlockQueue

	oldValueLen := indexes[id].dataLen

	size := oldValueLen + valueLen
	newSize := 0
	oldIndex := 0

	if size > indexes[id].blockSize {

		oldBlockQueue.index = indexes[id].index
		oldBlockQueue.blockId = indexes[id].blockId
		oldBlockQueue.count = blocksCount(indexes[id].blockSize)

		newSize = indexes[id].blockSize
		if newSize > 0 {
			for size > newSize {
				newSize *= 2
			}
			size = newSize
		}

		blocksCountNum = blocksCount(size)
		newBlockQueue = d.popBlockQueue(blocksCountNum)

		if newBlockQueue == nil {
			Logger.Fatal("pop block error")
		}

		oldIndex = oldBlockQueue.index
		index = newBlockQueue.index
		if indexes[id].blockSize > 0 && indexes[id].dataLen >= 0 &&
			oldIndex >= 0 && index >= 0 && d.IOs[index].file != nil &&
			d.IOs[oldIndex].file != nil {

			if d.isMmap {
				if d.IOs[index].mmap != nil && d.IOs[oldIndex].mmap != nil {
					oldPos := oldBlockQueue.blockId * BaseDbSize
					newPos := newBlockQueue.blockId * BaseDbSize

					for i := 0; i < oldValueLen; i++ {
						d.IOs[index].mmap[newPos+i] = d.IOs[oldIndex].mmap[oldPos+i]
					}

				} else {
					Logger.Fatal("mmap error")
				}

			} else {
				buf := make([]byte, oldValueLen)
				offsetSize := int64(oldBlockQueue.blockId * BaseDbSize)
				_, err := d.indexIO.file.ReadAt(buf[:oldValueLen], offsetSize)

				if err != nil {
					Logger.Fatal("read index file error", err)
				}
				_, err = d.IOs[index].file.WriteAt(buf, int64(newBlockQueue.blockId*BaseDbSize))
				if err != nil {
					Logger.Fatal("write index error")
				}
			}
		}
		indexes[id].index = newBlockQueue.index
		indexes[id].blockId = newBlockQueue.blockId
		indexes[id].blockSize = newBlockQueue.count * BaseDbSize
	}

	index = indexes[id].index

	if index >= 0 && d.IOs[index].file != nil {
		newPos := indexes[id].blockId*BaseDbSize + indexes[id].dataLen
		if d.isMmap && d.IOs[index].mmap != nil {
			for i := 0; i < valueLen; i++ {
				d.IOs[index].mmap[newPos+i] = value[i]
			}
		} else {
			_, err := d.IOs[index].file.WriteAt(value, int64(newPos))
			if err != nil {
				Logger.Fatal("write index error")
			}
		}

		indexes[id].dataLen += valueLen
		ret = id
	}

	if indexes[id].dataLen > d.state.dataLenMax {
		d.state.dataLenMax = indexes[id].dataLen
	}

	indexes[id].updateTime = time.Now().Unix()
	if oldBlockQueue.count > 0 {
		d.pushBlockQueue(oldBlockQueue.index, oldBlockQueue.blockId, oldBlockQueue.count*BaseDbSize)
	}

	// return
}

func (d *DB) SetTag(id int, tag int) {
	indexes := d.indexes
	if indexes == nil {
		Logger.Fatal("db index error")
	}

	if id >= 0 && id < MaxIndexSize && d.status == 0 {
		d.indexMutex.Lock()
		defer d.indexMutex.Unlock()
		d.checkIndexIOWithId(id)
		indexes[id].tag = tag
	}
}

func (d *DB) GetTag(id int) int {
	indexes := d.indexes
	if indexes == nil {
		Logger.Fatal("db index error")
	}

	if id >= 0 && id < MaxIndexSize && d.status == 0 {
		d.indexMutex.Lock()
		defer d.indexMutex.Unlock()
		d.checkIndexIOWithId(id)

		return indexes[id].tag
	}
	return -1
}

func (d *DB) Del(key []byte) (int, error) {
	if key == nil {
		return -1, errors.New("key is blank")
	}
	id, err := d.keyMapTrie.Set(key)

	if err != nil {
		return -1, err
	}

	return id, nil

}

func (d *DB) DelById(id int) int {
	if id >= 0 && id <= d.state.maxId {
		indexes := d.indexes
		index := indexes[id]
		if index.blockSize > 0 {
			d.pushBlockQueue(index.index, index.blockId, index.blockSize)

			d.lockId(id)
			defer d.unlockId(id)

			index.blockSize = 0
			index.blockId = 0
			index.dataLen = 0
			index.tag = 0
		}
		index.updateTime = time.Now().UnixNano()
		return id
	}

	return -1
}
