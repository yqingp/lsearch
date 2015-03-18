package lsearch

import (
	"fmt"
	"testing"
)

func TestNewIndexMapping(t *testing.T) {
	x := NewIndexMapping()
	x.Name = "test11"
	x.ShardNum = 8
	x.ReplicaNum = 1
	x.DefaultIndexerAnalyzer = "none"
	x.DefaultIndexerAnalyzer = "none"
	x.Fields["f1"] = FieldMapping{Name: "f1", Type: "INTEGER"}
	fmt.Println(x)
}
