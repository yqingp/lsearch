package lsearch

import (
// "errors"
)

const (
	DefaultShardNum          = -1
	DefaultReplicatNum       = 2
	IndexMappingMetaFileName = "lsearch.indexmapping.meta"
)

type IndexMapping struct {
	Name                   string
	Fields                 map[string]FieldMapping
	DefaultIndexerAnalyzer string
	DeafultSearchAanlyzer  string
	ReplicaNum             int
	ShardNum               int
}

// {:name => "", :fields => [{:name => aa, :}]}
func NewIndexMapping() IndexMapping {
	return IndexMapping{
		Name:                   "",
		ShardNum:               DefaultShardNum,
		ReplicaNum:             DefaultReplicatNum,
		DefaultIndexerAnalyzer: "none",
		DeafultSearchAanlyzer:  "none",
		Fields:                 make(map[string]FieldMapping),
	}
}

// func (self *IndexMapping) AddField(field map[string]string) error {
// }
