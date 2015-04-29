package util

import (
	"math"
)

func Float64ToInt64(f float64) int64 {
	fasint := int64(math.Float64bits(f))
	if fasint < 0 {
		fasint = fasint ^ 0x7fffffffffffffff
	}
	return fasint
}

func Int64ToFloat64(i int64) float64 {
	if i < 0 {
		i ^= 0x7fffffffffffffff
	}
	return math.Float64frombits(uint64(i))
}
