package lsearch

type IndexMapping struct {
	Name            string
	CreateTime      int64
	Fields          map[string]*FieldMapping
	IndexerAnalyzer string
	SearchAanlyzer  string
	ReplicaNum      int8
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

type IndexMappingTemplate map[string]interface{}

// {:name => "", :fields => [{:name => aa, :}]}

func NewDefaultIndexMapping() *IndexMapping {
	return &IndexMapping{
		Name:   "",
		Fields: make(map[string]*FieldMapping),
	}
}
