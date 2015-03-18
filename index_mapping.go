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

// create :settings => { :number_of_shards => (N || -1(auto)), :number_of_replicas => 2 } , :mappings => {
//   :xx => {
//     :_all => {
//         :indexAnalyzer => "ik",
//         :searchAnalyzer =>  "ik"
//     },
//     :properties => {
//       :id       => { :type => 'string', :index => 'not_analyzed', :include_in_all => false },
//       :name    => { :type => 'string', :index => 'not_analyzed'  },
//       :created_date    => { :type => 'string', :index => 'not_analyzed'  },#created_at.strftime("%y%m%d")
//       :created_time    => { :type => 'integer'},
//			 :content => {:type => "string"}
//     }
//   }
// }

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
