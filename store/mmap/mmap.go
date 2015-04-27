package mmap

import (
	"syscall"
)

type Mmap []byte

func MmapFile(fd int, mmapSize int) (Mmap, error) {
	data, err := syscall.Mmap(fd, 0, mmapSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m Mmap) Unmap() error {
	err := syscall.Munmap(m)

	return err
}
