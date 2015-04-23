package util

//-----------------------------------------------------------------------------
// MurmurHash3 was written by Austin Appleby, and is placed in the public
// domain. The author hereby disclaims copyright to this source code.

// Note - The x86 and x64 versions do _not_ produce the same results, as the
// algorithms are optimized for their respective platforms. You can still
// compile and run any of them on any platform, but your performance with the
// non-native version will be less than optimal.
//
// https://code.google.com/p/smhasher/source/browse/branches/chandlerc_dev/MurmurHash3.cpp

import (
    "unsafe"
)

const (
    c1  = 0xcc9e2d51
    c2  = 0x1b873593
)

// func rotl32(x, y uint32) uint32 {
//     return (x << y) | (x >> (32 - y))
// }

func MurmurHash3(key []byte) uint32 {
    len := uint32(len(key))
    blockNum := len / 4

    var seed uint32 = 0

    var i uint32 = 0

    for ; i < blockNum; i++ {
        t := key[i*4 : (i+1)*4]
        k1 := *(*uint32)(unsafe.Pointer(&t[0]))
        k1 *= c1
        k1 = (k1 << 15) | (k1 >> (32 - 15))
        k1 *= c2

        seed ^= k1
        seed = (seed << 13) | (seed >> (32 - 13))
        seed = seed*5 + 0xe6546b64
    }

    var k1 uint32 = 0

    switch len & 3 {
    case 3:
        k1 ^= uint32(key[len-1]) << 16
    case 2:
        k1 ^= uint32(key[len-1]) << 8
    case 1:
        k1 ^= uint32(key[len-1])
        k1 *= c1
        k1 = (k1 << 15) | (k1 >> (32 - 15))
        k1 *= c2
        seed ^= k1
    }

    seed ^= len
    seed ^= seed >> 16
    seed *= 0x85ebca6b
    seed ^= seed >> 13
    seed *= 0xc2b2ae35
    seed ^= seed >> 16

    return seed
}
