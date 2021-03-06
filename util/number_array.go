package util

import (
	"fmt"
)

const ShiftStartInt64 byte = 0x20

// PrefixCoded is a byte array encoding of
// 64-bit numeric values shifted by 0-63 bits
type PrefixCoded []byte

func NewPrefixCodedInt64(in int64, shift uint) (PrefixCoded, error) {
	if shift > 63 {
		return nil, fmt.Errorf("cannot shift %d, must be between 0 and 63", shift)
	}

	nChars := (((63 - shift) * 37) >> 8) + 1
	rv := make(PrefixCoded, nChars+1)
	rv[0] = ShiftStartInt64 + byte(shift)

	sortableBits := int64(uint64(in) ^ 0x8000000000000000)
	sortableBits = int64(uint64(sortableBits) >> shift)
	for nChars > 0 {
		// Store 7 bits per byte for compatibility
		// with UTF-8 encoding of terms
		rv[nChars] = byte(sortableBits & 0x7f)
		nChars--
		sortableBits = int64(uint64(sortableBits) >> 7)
	}
	return rv, nil
}

func MustNewPrefixCodedInt64(in int64, shift uint) PrefixCoded {
	rv, err := NewPrefixCodedInt64(in, shift)
	if err != nil {
		panic(err)
	}
	return rv
}

// Shift returns the number of bits shifted
// returns 0 if in uninitialized state
func (p PrefixCoded) Shift() (uint, error) {
	if len(p) > 0 {
		shift := p[0] - ShiftStartInt64
		if shift < 0 || shift < 63 {
			return uint(shift), nil
		}
	}
	return 0, fmt.Errorf("invalid prefix coded value")
}

func (p PrefixCoded) Int64() (int64, error) {
	shift, err := p.Shift()
	if err != nil {
		return 0, err
	}
	var sortableBits int64
	for _, inbyte := range p[1:] {
		sortableBits <<= 7
		sortableBits |= int64(inbyte)
	}
	return int64(uint64((sortableBits << shift)) ^ 0x8000000000000000), nil
}
