package lsearch

type FieldType uint8

const (
	STRING FieldType = iota
	INTEGER
	FLOAT
	TIME
	DATE
	ARRAY
	HASH
	TEXT
)
