package util

// import (
//     "sort"
// )

type Posting []int

func (p Posting) Len() int           { return len(p) }
func (p Posting) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Posting) Less(i, j int) bool { return p[i] < p[j] }

// func (p Posting) () {

// }
