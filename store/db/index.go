package db

import (
	. "github.com/yqingp/lsearch/store/mmap"
	"os"
	"path/filepath"
	"unsafe"
)

type Index struct {
	blockSize  int
	blockId    int
	dataLen    int
	index      int
	updateTime int64
}

func (d *DB) initIndex() {
	indexFilePath := filepath.Join(d.baseDir, IndexFileName)
	// Logger.Println(indexFilePath)
	f, err := os.OpenFile(indexFilePath, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		Logger.Fatal(err)
	}
	d.indexIO.file = f

	fstat, err := os.Stat(indexFilePath)
	if err != nil {
		Logger.Fatal(err)
	}

	d.indexIO.end = fstat.Size()
	d.indexIO.size = MaxIndexSize * SizeofIndex

	var errNo error

	fd := int(d.indexIO.file.Fd())

	if d.indexIO.mmap, errNo = MmapFile(fd, int(d.indexIO.size)); errNo != nil {
		Logger.Fatal(err)
	}

	d.indexes = (*[MaxIndexSize]Index)(unsafe.Pointer(&d.indexIO.mmap[0]))[:MaxIndexSize]
	if fstat.Size() == 0 {
		d.indexIO.end = BaseIndexSize * SizeofIndex

		if err := os.Truncate(indexFilePath, d.indexIO.end); err != nil {
			Logger.Fatal(err)
		}
	}

	// Logger.Println(d.indexes[0].index)
	d.indexIO.old = d.indexIO.end
}
