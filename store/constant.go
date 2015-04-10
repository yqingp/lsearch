package store

import (
    "unsafe"
)

const (
    MaxBlockQueueCount = 2 * 1024 * 1024
    MaxIndexSize       = 2000000000
    BaseIndexSize      = 1000000
    BaseDbSize         = 64
    MaxDirFileCount    = 64
    MaxMutexCount      = 65536
    MaxDbFileSize      = 256 * 1024 * 1024
    MaxDbFileCount     = 8129
)

const (
    SizeOfState      = int64(unsafe.Sizeof(State{}))
    SizeOfBlockQueue = int64(unsafe.Sizeof(BlockQueue{}))
    SizeofIndex      = int64(unsafe.Sizeof(Index{}))
)
