package util

import (
    "hash/fnv"
)

func Djb2Hash(str []byte) uint32 {
    if str == nil {
        return 0
    }

    var hash uint32 = 5381

    for _, v := range str {
        hash = (hash<<5 + hash) + uint32(v)
    }

    return hash
}

func BaseHash(str []byte) uint32 {
    h := fnv.New32a()
    h.Write(str)
    return h.Sum32()
}
