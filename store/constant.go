package store

import (
    "unsafe"
)

const (
    DB_LNK_MAX            = 2097152
    DB_LNK_INCREMENT      = 65536
    DB_DBX_MAX            = 2000000000
    DB_DBX_BASE           = 1000000
    DB_BASE_SIZE          = 64
    DB_PATH_MAX           = 1024
    DB_DIR_FILES          = 64
    DB_BUF_SIZE           = 4096
    DB_XBLOCKS_MAX        = 14
    DB_MBLOCKS_MAX        = 1024
    DB_MBLOCK_BASE        = 4096
    DB_MBLOCK_MAX         = 33554432
    DB_MUTEX_MAX          = 65536
    DB_USE_MMAP           = 0x01
    DB_MFILE_SIZE         = 268435456
    DB_MFILE_MAX          = 8129
    DB_BLOCK_INCRE_LEN    = 0x0
    DB_BLOCK_INCRE_DOUBLE = 0x1
)

var (
    SizeOfDbState          = int64(unsafe.Sizeof(DbState{}))
    SizeOfDbFreeBlockQueue = int64(unsafe.Sizeof(DbFreeBlockQueue{}))
    SizeofDbIndex          = int64(unsafe.Sizeof(DbIndex{}))
)
